// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"context"
	"fmt"
	"github.com/qinchende/gofast/fst/door"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/httpx"
	"github.com/qinchende/gofast/skill/timex"
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
	// 有两个特殊 RouteItem： 1. noRoute  2. noMethod
	// 这两个节点不参与构建路由树
	routerItem404 *RouteItem
	miniNode404   *radixMiniNode
	routerItem405 *RouteItem
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
	gft.resPool.New = func() interface{} {
		return &GFResponse{gftApp: gft, fitIdx: -1, ResWrap: &ResWriterWrap{}}
	}
	gft.ctxPool.New = func() interface{} {
		c := &Context{}
		// c.Pms = make(map[string]string)
		// c.match.needRTS = gft.RedirectTrailingSlash
		// c.GFResponse = &GFResponse{gftApp: gft}
		return c
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

	// TODO: 这里可以加入对全局路由的中间件函数（这里是已经匹配过路由的中间件）
	// TODO: 因为Server初始化之后就执行了这里，所以这里的中间件在客户自定义中间件之前
	// 加入路由的访问性能统计，但是这里还没有匹配路由，无法准确分路由统计
	if gft.EnableRouteMonitor {
		// 初始化全局的keeper变量
		door.InitKeeper(gft.FullPath)

		// 加入全局路由的中间件
		gft.Before(theFirstBeforeHandler)
		gft.After(theLastAfterHandler)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE：重构路由树。（重要！重要！重要！必须调用这个方法初始化路由树和中间件）
// 在不执行真正Listen的场景中，调用此函数能初始化服务器（必须要调用此函数来构造路由）
func (gft *GoFast) BuildRouters() {
	gft.readyOnce.Do(func() {
		gft.checkDefaultHandler()
		// TODO: 下面可以加入框架默认的Fits，用户自定义的fit只能在这些之前执行。

		// 这必须是最后一个Fit函数，由此进入下一级的 handlers
		gft.Fit(gft.serveHTTPWithCtx)
		// 依次执行 onReady 事件处理函数
		gft.execAppHandlers(gft.eReadyHds)
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

	// 记录请求进入的时间
	cRes.EnterTime = timex.Now()

	// 开始依次执行全局拦截器
	// 第二级的handler函数 (serveHTTPWithCtx) 的入口是这里的最后一个Fit函数
	cRes.NextFit(r)

	// 请求处理完成，开始回收资源
	if cRes.Ctx != nil {
		cRes.Ctx.GFResponse = nil
		gft.ctxPool.Put(cRes.Ctx)
		cRes.Ctx = nil
	}
	gft.resPool.Put(cRes)
}

// 承上启下：从全局拦截器 过度 到路由 handlers
// 全局拦截器过了之后，接下来就是查找路由进入下一阶段生命周期。
func (gft *GoFast) serveHTTPWithCtx(res *GFResponse, req *http.Request) {
	c := gft.ctxPool.Get().(*Context)
	res.Ctx = c
	c.GFResponse = res
	c.ReqRaw = req
	c.reset()
	gft.handleHTTPRequest(c)
}

// 处理所有请求，匹配路由并执行指定的路由处理函数
func (gft *GoFast) handleHTTPRequest(c *Context) {
	reqPath := c.ReqRaw.URL.Path
	unescape := false
	if gft.UseRawPath && len(c.ReqRaw.URL.RawPath) > 0 {
		reqPath = c.ReqRaw.URL.RawPath
		unescape = gft.UnescapePathValues
	}

	// 是否需要规范请求过来的URL，默认不需要
	if gft.RemoveExtraSlash {
		reqPath = httpx.CleanPath(reqPath)
	}

	// 看能不能找到 http method 对应的路由树
	miniRoot := gft.getMethodMiniRoot(c.ReqRaw.Method)
	if miniRoot != nil {
		// 开始在路由树中匹配 url path
		miniRoot.matchRoute(gft.fstMem, reqPath, &c.match, unescape)
		c.Params = c.match.params

		// 如果能匹配到路径
		if c.match.ptrNode != nil {
			// 第一种方案（默认）：两种不用的事件队列结构，看执行那一个
			c.execHandlers()
			// 第二种方案
			//c.execHandlersMini(nodeVal.ptrNode)

			c.ResWrap.WriteHeaderNow()
			return
		}

		// 匹配不到路由 先考虑 重定向
		// c.ReqRaw.Method != CONNECT && reqPath != [home index]
		if c.match.rts && c.ReqRaw.Method[0] != 'C' && reqPath != "/" {
			redirectTrailingSlash(c)
			return
		}
	}

	// 如果需要查找非本Method中的路由匹配，就尝试去找。
	// 找到了：就给出Method错误提示
	// 找不到：就走后面路由没匹配的逻辑
	if gft.HandleMethodNotAllowed {
		for _, tree := range gft.treeAll {
			if tree.method == c.ReqRaw.Method || tree.miniRoot == nil {
				continue
			}
			// 在别的 Method 路由树中匹配到了当前路径，返回提示 当前请求的 Method 错了。
			if tree.miniRoot.matchRoute(gft.fstMem, reqPath, &c.match, unescape); c.match.ptrNode != nil {
				c.match.ptrNode = gft.miniNode405
				c.Params = c.match.params
				c.execHandlers()
				return
			}
		}
	}

	c.match.ptrNode = gft.miniNode404
	// 如果没有匹配到任何路由，需要执行: 全局中间件 + noRoute handler
	c.execHandlers()
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 第二步：启动端口监听
// 说明：第一步和第二步之间，需要做所有的工作，主要就是初始化参数，设置所有的路由和处理函数
func (gft *GoFast) Listen(addr ...string) (err error) {
	// listen接受请求之前，必须调用这个来生成最终的路由树
	gft.BuildRouters()

	defer logx.DebugPrintError(err)
	// 只要 gft 实现了接口 ServeHTTP(ResponseWriter, *Request) 即可处理所有请求
	if addr == nil && gft.Addr != "" {
		addr = []string{gft.Addr}
	}
	gft.srv = &http.Server{Addr: httpx.ResolveAddress(addr), Handler: gft}

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
