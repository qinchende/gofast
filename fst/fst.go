// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"context"
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/httpx"
	"github.com/qinchende/gofast/skill/stat"
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
	srv         *http.Server // WebServer
	*AppConfig               // 引用配置
	*HomeRouter              // 根路由组（Root Group）
	appEvents                // 应用级事件
	fitHandlers IncHandlers  // 全局中间件处理函数，incoming request handlers
	resPool     sync.Pool    // 第一级：GFResponse context pools
	ctxPool     sync.Pool    // 第二级：Handler context pools
	readyOnce   sync.Once    // WebServer初始化只能执行一次
}

// 站点根目录是一个特殊的路由分组，所有其他分组都是他的子孙节点
type HomeRouter struct {
	// HomeRouter 本身就是一个路由分组
	RouterGroup

	// 记录当前Server所有的路由信息，方便后期重构路由树
	allRouters RouterItems
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

// 一个快速创建Server的函数，使用默认配置参数，方便调用。
// 记住：使用之前一定要先调用 ReadyToListen方法。
func Default() *GoFast {
	logx.DebugPrintWarningDefault()
	app := CreateServer(&AppConfig{
		RunMode: ProductMode,
	})
	return app
}

// 第一步：初始化一个 WebServer , 配置各种参数
func CreateServer(cfg *AppConfig) *GoFast {
	// 初始化当前环境变量
	app := new(GoFast)
	if cfg == nil {
		app.AppConfig = &AppConfig{}
	} else {
		app.AppConfig = cfg
	}
	app.initServerEnv()
	app.initResourcePool()
	app.initHomeRouter()
	return app
}

// 初始化资源池
func (gft *GoFast) initResourcePool() {
	gft.ctxPool.New = func() interface{} {
		c := &Context{}
		c.Pms = make(map[string]string)
		//c.GFResponse = &GFResponse{gftApp: gft}
		return c
	}

	gft.resPool.New = func() interface{} {
		cRes := &GFResponse{gftApp: gft, fitIdx: -1, ResWrap: &ResWrapriteWrap{}}
		return cRes
	}
}

// 初始化根路由树变量
func (gft *GoFast) initHomeRouter() {
	// 初始化 HomeRouter
	// 启动的时候，根分组"/"默认就有了，而且我们把他当做是一种特殊的Root节点
	// 方便将来加入 NoRoute、NoMethod 的处理Item
	gft.HomeRouter = &HomeRouter{}

	gft.hdsIdx = -1
	gft.prefix = "/"
	gft.gftApp = gft

	gft.allRouters = make(RouterItems, 0)
	gft.fstMem = new(fstMemSpace)
}

// 重构路由树，
// 在不执行真正Listen的场景中，调用此函数能初始化服务器（必须要调用此函数来构造路由）
func (gft *GoFast) BuildRouters() {
	gft.readyOnce.Do(func() {
		gft.checkDefaultHandler()
		gft.Fit(gft.serveHTTPWithCtx)      // 全局中间件过滤之后加入下一级的处理函数
		gft.execAppHandlers(gft.eReadyHds) // 依次执行 onReady 事件处理函数
	})
	gft.regAllRouters()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// http服务器，所有请求的入口，底层是用 goroutine 发起的一个协程任务
// 也就是说主线程，获取到任何请求事件（数据）之后，通过goroutine调用这个接口方法来并行处理
// 这里的代码就是在一个协程中运行的
// Note:
// 1. 这是请求进来之后的第一级上下文，为了节省内存空间，第一级的拦截器通过之后，会进入第二级更丰富的Context上下文（占用内存更多）
func (gft *GoFast) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cRes := gft.resPool.Get().(*GFResponse)
	cRes.ResWrap.Reset(w)
	cRes.reset()

	// 开始依次执行全局拦截器，第二级的handler函数的入口是这里的最后一个Fit函数
	cRes.NextFit(r)

	// 开始回收资源
	if cRes.Ctx != nil {
		cRes.Ctx.GFResponse = nil
		gft.ctxPool.Put(cRes.Ctx)
		cRes.Ctx = nil
	}
	gft.resPool.Put(cRes)
}

// 全局拦截器过了之后，接下来就是查找路由进入下一阶段生命周期。
func (gft *GoFast) serveHTTPWithCtx(res *GFResponse, req *http.Request) {
	c := gft.ctxPool.Get().(*Context)
	res.Ctx = c
	c.GFResponse = res
	c.ReqRaw = req
	c.reset()
	gft.handleHTTPRequest(c)
}

// TODO: 还有一些特殊情况的处理，需要在这里继续完善
// 处理所有请求，匹配路由并执行指定的路由处理函数
func (gft *GoFast) handleHTTPRequest(c *Context) {
	httpMethod := c.ReqRaw.Method
	rPath := c.ReqRaw.URL.Path
	//unescape := false
	//if gft.UseRawPath && len(c.ReqRaw.URL.RawPath) > 0 {
	//	rPath = c.ReqRaw.URL.RawPath
	//	unescape = gft.UnescapePathValues
	//}

	if gft.RemoveExtraSlash {
		rPath = httpx.CleanPath(rPath)
	}

	// 看能不能找到 http method 对应的路由树
	miniRoot := gft.getMethodMiniRoot(httpMethod)
	if miniRoot != nil {
		// 开始在路由树中匹配 url path
		miniRoot.matchRoute(gft.fstMem, rPath, &c.matchRst)

		// 如果能匹配到路径
		if c.matchRst.ptrNode != nil {
			c.Params = c.matchRst.params
			//c.ParseHttpParams() // 先解析 POST | GET 参数

			// 第一种方案（默认）：两种不用的事件队列结构，看执行那一个
			c.execHandlers(c.matchRst.ptrNode)
			// 第二种方案
			//c.execHandlersMini(nodeVal.ptrNode)

			c.ResWrap.WriteHeaderNow()
			return
			//} else {
			//c.ParseHttpParamsNoRoute()
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
				// TODO: 需要返回错误
				//c.handlers = engine.allNoMethod
				//serveError(c, http.StatusMethodNotAllowed, default405Body)
				//return

				c.execHandlers(gft.miniNode405)
				return
			}
		}
	}

	// 如果没有匹配到任何路由，需要执行: 全局中间件 + noRoute handler
	c.execHandlers(gft.miniNode404)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 第二步：启动端口监听
// 说明：第一步和第二步之间，需要做所有的工作，主要就是初始化参数，设置所有的路由和处理函数
func (gft *GoFast) Listen(addr ...string) (err error) {
	// listen接受请求之前，必须调用这个来生成最终的路由树
	gft.BuildRouters()

	defer logx.DebugPrintError(err)
	// 只要 gft 实现了接口 ServeHTTP(ResponseWriter, *Request) 即可处理所有请求
	gft.srv = &http.Server{Addr: httpx.ResolveAddress(addr), Handler: gft}

	// 设置关闭前等待时间
	go func() {
		err = gft.srv.ListenAndServe()
	}()
	gft.GracefulShutdown()
	return
}

// 创建日志模板
func (gft *GoFast) CreateMetrics() *stat.Metrics {
	var metrics *stat.Metrics

	if len(gft.Name) > 0 {
		metrics = stat.NewMetrics(gft.Name)
	} else {
		metrics = stat.NewMetrics(gft.Addr)
	}

	return metrics
}

// 优雅关闭
func (gft *GoFast) GracefulShutdown() {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// 执行 onClose 事件订阅函数
	gft.execAppHandlers(gft.eCloseHds)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gft.SecondsBeforeShutdown)*time.Millisecond)
	defer cancel()
	if err := gft.srv.Shutdown(ctx); err != nil {
		fmt.Sprintln("Server Shutdown Error: ", err)
	}
	<-ctx.Done()
}
