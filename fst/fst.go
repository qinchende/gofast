// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

import (
	"context"
	"fmt"
	"gofast/skill"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// GoFast is the framework's instance.
// Create an instance of GoFast, by using CreateServer().
type GoFast struct {
	srv *http.Server
	*AppConfig
	*HomeRouter
	appEvents

	pool      sync.Pool
	readyOnce sync.Once
}

// 站点根目录是一个特殊的路由分组，所有其他分组都是他的子孙节点
type HomeRouter struct {
	RouterGroup
	// 有两个特殊 RouterItem： 1. noRoute  2. noMethod
	// 这两个节点不参与构建路由树
	routerItem404 *RouterItem
	miniNode404   *radixMiniNode

	routerItem405 *RouterItem
	miniNode405   *radixMiniNode

	// 虽然支持 RESTFUL 路由规范，但 GET 和 POST 是一等公民。
	// 绝大部分应用Get和Post路由居多，我们能尽快匹配就不需要无用的Method比较选择的过程
	treeGet    *methodTree
	treePost   *methodTree
	treeOthers methodTrees
	treeAll    methodTrees

	// 主要以数组结构的形式，存储了 Routes & Handlers
	fstMem *fstMemSpace
}

// 第一步：初始化一个 WebServer , 配置各种参数
func CreateServer(cfg *AppConfig) *GoFast {
	// 初始化当前环境变量
	gft := new(GoFast)
	if cfg == nil {
		gft.AppConfig = &AppConfig{}
	} else {
		gft.AppConfig = cfg
	}
	gft.initServerEnv()

	// 初始化 HomeRouter
	// 启动的时候，根分组"/"默认就有了，而且我们把他当做是一种特殊的最后一级节点
	// 方便将来加入 NoRoute、NoMethod 的处理Item
	hm := &HomeRouter{}
	hm.hdsGroupIdx = -1
	hm.prefix = "/"
	hm.gftApp = gft
	gft.HomeRouter = hm

	gft.treeGet = &methodTree{method: http.MethodGet}
	gft.treePost = &methodTree{method: http.MethodPost}
	gft.treeOthers = make(methodTrees, 0, 9)

	gft.pool.New = func() interface{} {
		return &Context{gftApp: gft}
	}

	gft.fstMem = new(fstMemSpace)
	return gft
}

// 一个快速创建Server的函数，使用默认配置参数，方便调用。
// 记住：使用之前一定要先调用 ReadyToListen方法。
func Default() *GoFast {
	skill.DebugPrintWARNINGDefault()
	app := CreateServer(&AppConfig{
		RunMode: ProductMode,
	})
	return app
}

// Ready to listen the ip address
// 在不执行真正Listen的场景中，调用此函数能初始化服务器
func (gft *GoFast) ReadyToListen() {
	// 服务Listen之前，只执行一次初始化
	gft.readyOnce.Do(func() {
		gft.checkDefaultHandler()
		// 设置 treeAll
		lenTreeOthers := len(gft.treeOthers)
		ifPanic(lenTreeOthers > 7, "Too many kind of methods")
		gft.treeAll = gft.treeOthers[:lenTreeOthers:9]
		gft.treeAll = append(gft.treeAll, gft.treeGet)
		gft.treeAll = append(gft.treeAll, gft.treePost)

		if gft.PrintRouteTrees {
			gft.printRouteTrees()
		}
		// 这里开始完整的重建整个路由树的数据结构
		gft.buildMiniRoutes()
		// 依次执行 onReady 事件处理函数
		gft.execHandlers(gft.eReadyHds)
	})
}

//func (gft *GoFast) Listen(addr ...string) (err error) {
//	gft.ReadyToListen()
//
//	defer func() { skill.DebugPrintError(err) }()
//	err = http.ListenAndServe(skill.ResolveAddress(addr), gft)
//	return
//}

// 第二步：启动端口监听
// 说明：第一步和第二步之间，需要做所有的工作，主要就是初始化参数，设置所有的路由和处理函数
func (gft *GoFast) Listen(addr ...string) (err error) {
	gft.ReadyToListen()

	defer func() { skill.DebugPrintError(err) }()
	gft.srv = &http.Server{Addr: skill.ResolveAddress(addr), Handler: gft}

	// 设置关闭前等待时间
	go func() {
		err = gft.srv.ListenAndServe()
	}()
	gft.GracefulShutdown()
	return
}

// 优雅关闭
func (gft *GoFast) GracefulShutdown() {
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// 执行 onClose 事件订阅函数
	gft.execHandlers(gft.eCloseHds)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gft.SecondsBeforeShutdown)*time.Second)
	defer cancel()
	if err := gft.srv.Shutdown(ctx); err != nil {
		fmt.Sprintln("Server Shutdown Error: ", err)
	}
	<-ctx.Done()
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
		// 开始在路由树中匹配 url path
		miniRoot.matchRoute(gft.fstMem, rPath, &c.matchRst)

		// 如果能匹配到路径
		if c.matchRst.ptrNode != nil {
			c.Params = c.matchRst.params

			// 第一种方案（默认）：两种不用的事件队列结构，看执行那一个
			c.execHandlers(c.matchRst.ptrNode)
			// 第二种方案
			//c.execHandlersMini(nodeVal.ptrNode)

			c.resW.WriteHeaderNow()
			return
		}
		// 匹配不到 先考虑 重定向
		if httpMethod != "CONNECT" && rPath != "/" {
			if c.matchRst.tsr && gft.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			//if gft.RedirectFixedPath && redirectFixedPath(c, miniRoot, gft.RedirectFixedPath) {
			//	return
			//}
		}
	}

	// 如果需要查找非本Method中的路由匹配，就尝试去找。
	// 找到了：就给出Method错误提示
	// 找不到：就走后面路由没匹配的逻辑
	if gft.HandleMethodNotAllowed {
		for _, tree := range gft.treeAll {
			if tree.method == httpMethod || tree.miniRoot == nil {
				continue
			}
			// 在别的 Method 路由树中匹配到了当前路径，返回提示 当前请求的 Method 错了。
			if tree.miniRoot.matchRoute(gft.fstMem, rPath, &c.matchRst); c.matchRst.ptrNode != nil {
				c.execHandlers(gft.miniNode405)
				return
			}
		}
	}

	// 如果没有匹配到任何路由，需要执行: 全局中间件 + noRoute handler
	c.execHandlers(gft.miniNode404)
}
