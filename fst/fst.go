// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"context"
	"fmt"
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
	*GfConfig              // 引用配置
	srv       *http.Server // WebServer
	appEvents              // 应用级事件
	readyOnce sync.Once    // WebServer初始化只能执行一次

	// 根路由组相关属性
	*HomeRouter                  // 根路由组（Root Group）
	fitHandlers []FitFunc        // 全局中间件处理函数，incoming request handlers
	fitEnter    http.HandlerFunc // fit系列中间件函数的入口，请求进入之后第一个接收函数
	ctxPool     sync.Pool        // 第二级：Handler context pools (第一级是标准形式，不需要缓冲池)
}

// 站点根目录是一个特殊的路由分组，所有其他分组都是他的子孙节点
type HomeRouter struct {
	RouteGroup // HomeRouter 本身就是一个路由分组

	// 有两个特殊 RouteItem： 1. noRoute  2. noMethod
	// 这两个节点不参与构建路由树
	miniNode404 *radixMiniNode
	miniNode405 *radixMiniNode
	allRoutes   []*RouteItem // 记录当前Server所有的路由信息，方便后期重构路由树

	// 虽然支持 RESTFUL 路由规范，但 GET 和 POST 是一等公民。
	// 绝大部分应用Get和Post路由居多，我们能尽快匹配就不需要无用的Method比较选择的过程
	routerTrees methodTrees

	// 主要以数组结构的形式，存储了 Routes & Handlers
	fstMem *fstMemSpace
}

// 一个快速创建Server的函数，使用默认配置参数，方便调用。
// 记住：使用之前一定要先调用 ReadyToListen方法。
func Default() *GoFast {
	app := CreateServer(&GfConfig{
		RunningMode: ProductMode,
	})
	return app
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 第一步：初始化一个 WebServer , 配置各种参数
func CreateServer(cfg *GfConfig) *GoFast {
	// 初始化当前环境变量
	app := new(GoFast)
	if cfg == nil {
		app.GfConfig = &GfConfig{}
	} else {
		app.GfConfig = cfg
	}
	app.initServerConfig()
	app.initResourcePool()
	app.initHomeRouter()
	return app
}

// 初始化资源池
func (gft *GoFast) initResourcePool() {
	gft.ctxPool.New = func() any {
		c := &Context{myApp: gft, ResWrap: &ResponseWrap{}}
		// c.Pms = make(map[string]string)
		// c.route.needRTS = gft.RedirectTrailingSlash
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
	gft.myApp = gft

	gft.allRoutes = make([]*RouteItem, 0)
	// 默认为空的节点
	gft.allRoutes = append(gft.allRoutes, &RouteItem{
		group:    nil,
		fullPath: "*",
		routeIdx: 0,
	})
	// 404 Default Route
	gft.allRoutes = append(gft.allRoutes, &RouteItem{
		group:    &gft.RouteGroup,
		fullPath: "/404",
		routeIdx: 1,
	})
	// 405 Default Route
	gft.allRoutes = append(gft.allRoutes, &RouteItem{
		group:    &gft.RouteGroup,
		fullPath: "/405",
		routeIdx: 2,
	})
	gft.fstMem = &fstMemSpace{myApp: gft}

	//// TODO: 这里可以加入对全局路由的中间件函数（这里是已经匹配过路由的中间件）
	//// TODO: 因为Server初始化之后就执行了这里，所以这里的中间件在客户自定义中间件之前
	//// 加入路由的访问性能统计，但是这里还没有匹配路由，无法准确分路由统计
	//if gft.EnableRouteMonitor {
	//	// 初始化全局的keeper变量
	//	door.InitKeeper(gft.FullPath)
	//
	//	// 加入全局路由的中间件
	//	gft.Before(theFirstBeforeHandler)
	//	gft.After(theLastAfterHandler)
	//}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// http服务器，所有请求的入口，底层是用 goroutine 发起的一个协程任务
// 也就是说主线程，获取到任何请求事件（数据）之后，通过goroutine调用这个接口方法来并行处理
// 这里的代码就是在一个协程中运行的
// Note:
// 1. 这是请求进来之后的第一级上下文，为了节省内存空间，第一级的拦截器通过之后，会进入第二级更丰富的Context上下文（占用内存更多）
func (gft *GoFast) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 开始依次执行全局拦截器，开始是第一级的fits系列
	// 第二级的handler函数 (serveHTTPWithCtx) 的入口是这里的最后一个Fit函数
	gft.fitEnter(w, r)
}

// 这里是第二级执行链
// 承上启下：从全局拦截器 过度 到路由 handlers
// 全局拦截器过了之后，接下来就是查找路由进入下一阶段生命周期。
func (gft *GoFast) serveHTTPWithCtx(w http.ResponseWriter, r *http.Request) {
	c := gft.ctxPool.Get().(*Context)
	c.EnterTime = timex.Now() // 请求开始进入上下文阶段，开始计时
	c.ResWrap.Reset(w)
	c.ReqRaw = r
	c.reset()
	gft.handleHTTPRequest(c)
	// 超时引发的对象不能放回缓存池
	if !c.IsTimeout {
		gft.ctxPool.Put(c)
	}
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

	// 以下分A、B、C三步走
	// A. 看能不能找到 http method 对应的路由树
	miniRoot := gft.getMethodMiniRoot(c.ReqRaw.Method)
	if miniRoot != nil {
		// 开始在路由树中匹配 url path
		miniRoot.matchRoute(gft.fstMem, reqPath, &c.route, unescape)
		c.UrlParams = c.route.params

		// 如果能匹配到路径
		if c.route.ptrNode != nil {
			// 进一步的check，比如可以在这里跳转成404；或者直接AbortDirect
			c.execAfterMatchHandlers()

			c.execHandlers() // 执行处理链
			return
		}

		// 检查重定向
		// c.ReqRaw.Method != CONNECT && reqPath != [home index]
		if c.route.rts && c.ReqRaw.Method[0] != 'C' && reqPath != "/" {
			redirectTrailingSlash(c)
			return
		}
	}

	// B. 可以尝试是否不同的Method中能匹配路由
	// 如果需要查找非本Method中的路由匹配，就尝试去找。
	// 找到了：就给出Method错误提示
	// 找不到：就走后面路由没匹配的逻辑
	if gft.CheckOtherMethodRoute {
		for _, tree := range gft.routerTrees {
			if tree.method == c.ReqRaw.Method || tree.miniRoot == nil {
				continue
			}
			// 在别的 Method 路由树中匹配到了当前路径，返回提示 当前请求的 Method 错了。
			if tree.miniRoot.matchRoute(gft.fstMem, reqPath, &c.route, unescape); c.route.ptrNode != nil {
				c.route.ptrNode = gft.miniNode405
				c.UrlParams = c.route.params

				c.execHandlers()
				return
			}
		}
	}

	// C. 以上都无法匹配，就走404逻辑
	c.route.ptrNode = gft.miniNode404
	c.execHandlers()
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE：重构路由树。（重要！重要！重要！必须调用这个方法初始化路由树和中间件）
// 在不执行真正Listen的场景中，调用此函数能初始化服务器（必须要调用此函数来构造路由）
func (gft *GoFast) BuildRoutes() {
	gft.readyOnce.Do(func() {
		gft.initDefaultHandlers()
		// TODO: 下面可以加入框架默认的Fits，用户自定义的fit只能在这些之前执行。
		// 这必须是最后一个Fit函数，由此进入下一级的 handlers
		gft.bindContextFit(gft.serveHTTPWithCtx)
	})
	gft.execAppHandlers(gft.eBeforeBuildRoutesHds) // before build routes
	gft.buildAllRoutes()
	gft.execAppHandlers(gft.eAfterBuildRoutesHds) // after build routes
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 第二步：启动端口监听
// 说明：第一步和第二步之间，需要做所有的工作，主要就是初始化参数，设置所有的路由和处理函数
func (gft *GoFast) Listen(addr ...string) (err error) {
	// listen接受请求之前，必须调用这个来生成最终的路由树
	gft.BuildRoutes()

	// 依次执行 onReady 事件处理函数
	gft.execAppHandlers(gft.eReadyHds)

	defer logx.Stacks(err)
	// 只要 gft 实现了接口 ServeHTTP(ResponseWriter, *Request) 即可处理所有请求
	if addr == nil && gft.ListenAddr != "" {
		addr = []string{gft.ListenAddr}
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
	// kill -9 is syscall. SIGKILL but can't catch
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// 执行 onClose 事件订阅函数
	gft.execAppHandlers(gft.eCloseHds)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gft.BeforeShutdownMS)*time.Millisecond)
	defer cancel()
	if err := gft.srv.Shutdown(ctx); err != nil {
		fmt.Sprintln("Server Shutdown Error: ", err)
	}
	<-ctx.Done()
}
