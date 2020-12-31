package fst

import (
	"math"
)

// 所有注册路由方法都走这个函数
func (ft *Faster) regRoute(method, path string, ri *RouterItem) {
	ifPanic(path[0] != '/', "Path must begin with '/'")
	ifPanic(len(path) > math.MaxUint8, "The path length is more than 255 chars")
	ifPanic(len(method) == 0, "HTTP method can not be empty")

	mTree := ft.getMethodTree(method)
	if mTree == nil {
		mTree = &methodTree{method: method, root: nil}
		ft.treeOthers = append(ft.treeOthers, *mTree)
	}
	mTree.regRoute(path, ri)
}

//
//// NoRoute adds handlers for NoRoute. It return a 404 code by default.
//func (ft *Faster) NoRoute(handlers ...CHandler) {
//	ft.noRoute = handlers
//	ft.rebuild404Handlers()
//}
//
//// NoMethod sets the handlers called when... TODO.
//func (ft *Faster) NoMethod(handlers ...CHandler) {
//	ft.noMethod = handlers
//	ft.rebuild405Handlers()
//}

func (ft *Faster) getMethodMiniRoot(method string) (tRoot *miniNode) {
	switch method[0] {
	case 'P':
		if method[1] == 'O' {
			tRoot = ft.treePost.miniRoot
		} else {
			tRoot = ft.treeOthers.getTreeMiniRoot(method)
		}
	case 'G':
		tRoot = ft.treeGet.miniRoot
	default:
		tRoot = ft.treeOthers.getTreeMiniRoot(method)
	}
	return
}

func (ft *Faster) getMethodTree(method string) (tree *methodTree) {
	switch method[0] {
	case 'P':
		if method[1] == 'O' {
			tree = ft.treePost
		} else {
			tree = ft.treeOthers.getTree(method)
		}
	case 'G':
		tree = ft.treeGet
	default:
		tree = ft.treeOthers.getTree(method)
	}
	return
}

//func redirectTrailingSlash(c *Context) {
//	req := c.Request
//	p := req.URL.Path
//	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
//		p = prefix + "/" + req.URL.Path
//	}
//	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
//	if req.Method != "GET" {
//		code = http.StatusTemporaryRedirect
//	}
//
//	req.URL.Path = p + "/"
//	if length := len(p); length > 1 && p[length-1] == '/' {
//		req.URL.Path = p[:length-1]
//	}
//	skill.DebugPrint("redirecting request %d: %s --> %s", code, p, req.URL.String())
//	http.Redirect(c.Reply, req, req.URL.String(), code)
//	c.resW.WriteHeaderNow()
//}
