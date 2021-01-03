// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

import (
	"gofast/skill"
	"net/http"
	"sync"
)

// GoFast is the framework's instance.
// Create an instance of GoFast, by using CreateServer().
type GoFast struct {
	*AppConfig
	appEvents

	home       *HomeSite
	treeGet    *methodTree
	treePost   *methodTree
	treeOthers methodTrees
	pool       sync.Pool
}

// 站点根目录是一个特殊的分组
type HomeSite struct {
	RouterGroup
	// 有两个特殊 RouterItem： 1. noRoute  2. noMethod
	// 这两个节点不参与构建路由树
	routerItem404 *RouterItem
	miniNode404   *radixMiniNode
	routerItem405 *RouterItem
	miniNode405   *radixMiniNode
}

// 第一步：初始化一个 WebServer , 配置各种参数
func CreateServer(cfg *AppConfig) (*GoFast, *HomeSite) {
	// 初始化当前环境变量
	gft := new(GoFast)
	if cfg == nil {
		gft.AppConfig = &AppConfig{}
	} else {
		gft.AppConfig = cfg
	}
	gft.initServerEnv()

	// 初始化 HomeSite
	// 启动的时候，根分组"/"默认就有了，而且我们把他当做是一种特殊的最后一级节点
	// 方便将来加入 NoRoute、NoMethod 的处理Item
	hm := &HomeSite{}
	hm.hdsGroupIdx = -1
	hm.prefix = "/"
	hm.gftApp = gft
	gft.home = hm

	// 虽然支持 RESTFUL 路由规范，但 GET 和 POST 是一等公民.
	gft.treeGet = &methodTree{method: http.MethodGet}
	gft.treePost = &methodTree{method: http.MethodPost}
	gft.treeOthers = make(methodTrees, 0, 7)

	gft.pool.New = func() interface{} {
		return &Context{gftApp: gft}
	}
	return gft, gft.home
}

// Ready to listen the ip address
// 在不执行真正Listen的场景中，调用此函数能初始化服务器
func (gft *GoFast) ReadyToListen() {
	gft.checkDefaultHandler()

	if gft.PrintRouteTrees {
		gft.printRouteTrees()
	}
	// 这里开始完整的重建整个路由树的数据结构
	gft.rebuildRoutes()
	// 依次执行 onReady 事件处理函数
	gft.execHandlers(gft.eReadyHds)
}

// 第二步：启动端口监听
// 说明：第一步和第二步之间，需要做所有的工作，主要就是初始化参数，设置所有的路由和处理函数
func (gft *GoFast) Listen(addr ...string) (err error) {
	gft.ReadyToListen()

	defer func() { skill.DebugPrintError(err) }()
	err = http.ListenAndServe(skill.ResolveAddress(addr), gft)
	return
}

// http服务器，所有请求的入口，底层是用 goroutine 发起的一个协程任务
// 也就是说主线程，获取到任何请求事件（数据）之后，通过goroutine调用这个接口方法来并行处理
// 这里的代码就是在一个协程中运行的
func (gft *GoFast) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c := gft.pool.Get().(*Context)
	c.resW.Reset(res)
	c.Request = req
	c.reset()
	gft.handleHTTPRequest(c)
	gft.pool.Put(c)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 还有大量的逻辑需要在这里支持
// 处理所有请求，匹配路由并执行指定的路由处理函数
func (gft *GoFast) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	//unescape := false
	//if gft.UseRawPath && len(c.Request.URL.RawPath) > 0 {
	//	rPath = c.Request.URL.RawPath
	//	unescape = gft.UnescapePathValues
	//}

	if gft.RemoveExtraSlash {
		rPath = skill.CleanPath(rPath)
	}

	// 看能不能找到 http method 对应的路由树
	miniRoot := gft.getMethodMiniRoot(httpMethod)
	if miniRoot != nil {
		nodeVal := miniRoot.matchRoute(rPath, c.Params)
		// 如果能匹配到路径
		if nodeVal.ptrNode != nil {
			c.Params = nodeVal.params

			// 第一种方案（默认）：两种不用的事件队列结构，看执行那一个
			c.execHandlers(nodeVal.ptrNode)
			// 第二种方案
			//c.execHandlersMini(nodeVal.ptrNode)

			c.resW.WriteHeaderNow()
			return
		}
		// 匹配不到 先考虑 重定向
		if httpMethod != "CONNECT" && rPath != "/" {
			if nodeVal.tsr && gft.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			//if gft.RedirectFixedPath && redirectFixedPath(c, miniRoot, gft.RedirectFixedPath) {
			//	return
			//}
		}
	}

	if gft.HandleMethodNotAllowed {
		for _, tree := range gft.treeOthers {
			if tree.method == httpMethod {
				continue
			}
			// 在别的 Method 路由树中匹配到了当前路径，返回提示 当前请求的 Method 错了。
			if nodeVal := tree.miniRoot.matchRoute(rPath, c.Params); nodeVal.ptrNode != nil {
				c.execHandlers(gft.home.miniNode405)
				return
			}
		}
	}

	// 如果没有匹配到任何路由，需要执行: 全局中间件 + noRoute handler
	c.execHandlers(gft.home.miniNode404)
}
