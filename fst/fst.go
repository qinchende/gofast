package fst

import (
	"gofast/fst/render"
	"gofast/skill"
	"net/http"
	"sync"
)

// 应用程序开始的地方
type Faster struct {
	FConfig
	appEvents
	home       *RouterGroup
	pool       sync.Pool
	treeGet    *methodTree
	treePost   *methodTree
	treeOthers methodTrees

	secureJsonPrefix string
	HTMLRender       render.HTMLRender
	//FuncMap          template.FuncMap

	allNoRoute  CHandlers
	allNoMethod CHandlers
	//noRoute     CHandlers
	//noMethod    CHandlers

}

// 第一步：初始化一个 WebServer , 配置各种参数
func CreateServer(cfg *FConfig) (*Faster, *RouterGroup) {
	cfg.initServerEnv()
	ft := new(Faster)

	ft.home = &RouterGroup{
		prefix: "/",
		faster: ft,
	}
	fstMem.hdsGroupCt++

	ft.treeGet = &methodTree{method: http.MethodGet}
	ft.treePost = &methodTree{method: http.MethodPost}
	ft.treeOthers = make(methodTrees, 0, 7)

	ft.pool.New = func() interface{} {
		return &Context{faster: ft}
	}
	return ft, ft.home
}

// Ready to listen the ip address
func (ft *Faster) ReadyToListen() {
	if FEnv.PrintRouteTrees {
		ft.printRouteTrees()
	}
	// 这里开始完整的重建整个路由树的数据结构
	ft.rebuildRoutes()
	// 依次执行 onReady 事件处理函数
	ft.execHandlers(ft.eReadyHds)
}

// 第二步：启动端口监听
// 说明：第一步和第二步之间，需要做所有的工作，主要就是初始化参数，设置所有的路由和处理函数
func (ft *Faster) Listen(addr ...string) (err error) {
	ft.ReadyToListen()

	defer func() { skill.DebugPrintError(err) }()
	err = http.ListenAndServe(resolveAddress(addr), ft)
	return
}

// http服务器，所有请求的入口，底层是用 goroutine 发起的一个协程任务
// 也就是说主线程，获取到任何请求事件（数据）之后，通过goroutine调用这个接口方法来并行处理
func (ft *Faster) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c := ft.pool.Get().(*Context)
	c.resW.Reset(res)
	c.Request = req
	c.reset()
	ft.handleHTTPRequest(c)
	ft.pool.Put(c)
}

// 处理所有请求，匹配路由并执行指定的路由处理函数
func (ft *Faster) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	//unescape := false
	//if ft.UseRawPath && len(c.Request.URL.RawPath) > 0 {
	//	rPath = c.Request.URL.RawPath
	//	unescape = ft.UnescapePathValues
	//}

	if ft.RemoveExtraSlash {
		rPath = skill.CleanPath(rPath)
	}
	miniRoot := ft.getMethodMiniRoot(httpMethod)
	if miniRoot != nil {
		nodeVal := miniRoot.matchRoute(rPath, c.Params)
		if nodeVal.ptrNode != nil {
			c.Params = nodeVal.params

			// 第一种方案（默认）：两种不用的事件队列结构，看执行那一个
			c.execHandlers(nodeVal.ptrNode)
			// 第二种方案
			//c.execHandlersMini(nodeVal.ptrNode)

			c.resW.WriteHeaderNow()
			return
		}
		//if httpMethod != "CONNECT" && rPath != "/" {
		//	if nodeVal.tsr && ft.RedirectTrailingSlash {
		//		redirectTrailingSlash(c)
		//		return
		//	}
		//	if ft.RedirectFixedPath && redirectFixedPath(c, miniRoot, ft.RedirectFixedPath) {
		//		return
		//	}
		//}
	}
	//return
	//
	//if ft.HandleMethodNotAllowed {
	//	for _, tree := range ft.treeOthers {
	//		if tree.method == httpMethod {
	//			continue
	//		}
	//		//if nodeVal := tree.root.findRoute(rPath, nil, unescape); nodeVal.handlers != nil {
	//		if nodeVal := tree.miniRoot.matchRoute(rPath); nodeVal.handlers != nil {
	//			c.handlers = ft.allNoMethod
	//			serveError(c, http.StatusMethodNotAllowed, default405Body)
	//			return
	//		}
	//	}
	//}
	// c.handlers = ft.allNoRoute
	requestError(c, http.StatusNotFound, default404Body)
}

func requestError(c *Context, code int, defaultMessage []byte) {
	c.resW.Status = code
	if c.resW.Written() {
		return
	}
	if c.resW.Status == code {
		c.resW.Header()["Content-Type"] = mimePlain
		_, err := c.Reply.Write(defaultMessage)
		if err != nil {
			skill.DebugPrint("Cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.resW.WriteHeaderNow()
	return
}
