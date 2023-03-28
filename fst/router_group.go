// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

//func (gp *RouteGroup) Attrs(attrs *GAttrs) *RouteGroup {
//	if gp.attrs == nil {
//		gp.attrs = attrs
//	} else {
//		gp.attrs.PmsFields = append(gp.attrs.PmsFields, attrs.PmsFields...)
//	}
//	return gp
//}

// Note：如果分组已经存在，需要报错。 或者不报错。
// GoFast选择不报错，允许添加相同路径的不同分组，区别应用不同的特性
func (gp *RouteGroup) Group(relPath string) *RouteGroup {
	gpNew := &RouteGroup{
		prefix: gp.fixAbsolutePath(relPath),
		myApp:  gp.myApp,
		hdsIdx: -1,
	}
	gp.children = append(gp.children, gpNew)
	return gpNew
}

// Prefix returns the base path of router gp.
// For example, if v := router.Group("/rest/n/v1/api"), v.Prefix() is "/rest/n/v1/api".
func (gp *RouteGroup) Prefix() string {
	return gp.prefix
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// idx 0. not match any 1. 404 handlers 2. 405 handlers
func (gft *GoFast) regSpecialHandlers(hds []CtxHandler, idx int) {
	cst.PanicIf(len(hds) <= 0, "there must be at least one handler")
	cst.PanicIf(len(gft.allRoutes[idx].eHds) > 0, "handlers already exists.")
	gft.allRoutes[idx].eHds = addCtxHandlers(gft.fstMem, hds)
}

// 在当前分组注册一个新节点
func (gp *RouteGroup) register(httpMethod, relPath string, hds []CtxHandler) *RouteItem {
	// 最终的路由绝对路径
	absPath := gp.fixAbsolutePath(relPath)
	cst.PanicIf(absPath[0] != '/', "Path must begin with '/'")
	cst.PanicIf(len(absPath) > math.MaxUint8, "The path is more than 255 chars")
	cst.PanicIf(len(httpMethod) == 0, "HTTP method can not be empty")

	// 新添加一个 GroupItem，记录所有的处理函数
	ri := &RouteItem{
		method:   httpMethod,
		fullPath: absPath,
		group:    gp,
		routeIdx: 0,
	}

	// 可以 RouteItem 只创建对象，不指定处理函数。等后面再设置
	if len(hds) == 0 {
		return ri
	} else {
		return ri.Handle(hds...)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 有个问题，httpMethod参数没有做枚举校验，可以创建任意名称的method路由数，真要这么自由吗???
func (gp *RouteGroup) Handle(httpMethod, relPath string, hds ...CtxHandler) *RouteItem {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		panic("http method " + httpMethod + " is not valid")
	}
	return gp.register(httpMethod, relPath, hds)
}

func (gp *RouteGroup) Get(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodGet, relPath, hds)
}

func (gp *RouteGroup) Post(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPost, relPath, hds)
}

func (gp *RouteGroup) Delete(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodDelete, relPath, hds)
}

func (gp *RouteGroup) Patch(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPatch, relPath, hds)
}

func (gp *RouteGroup) Put(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodPut, relPath, hds)
}

func (gp *RouteGroup) Options(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodOptions, relPath, hds)
}

func (gp *RouteGroup) Head(relPath string, hds ...CtxHandler) *RouteItem {
	return gp.register(http.MethodHead, relPath, hds)
}

// 特殊类型
func (gp *RouteGroup) GetPost(relPath string, hds ...CtxHandler) RouteItems {
	get := gp.register(http.MethodGet, relPath, hds)
	post := gp.register(http.MethodPost, relPath, hds)
	return RouteItems{get, post}
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (gp *RouteGroup) All(relPath string, hds ...CtxHandler) RouteItems {
	get := gp.register(http.MethodGet, relPath, hds)
	post := gp.register(http.MethodPost, relPath, hds)
	put := gp.register(http.MethodPut, relPath, hds)
	patch := gp.register(http.MethodPatch, relPath, hds)
	head := gp.register(http.MethodHead, relPath, hds)
	opt := gp.register(http.MethodOptions, relPath, hds)
	del := gp.register(http.MethodDelete, relPath, hds)
	conn := gp.register(http.MethodConnect, relPath, hds)
	trace := gp.register(http.MethodTrace, relPath, hds)
	return RouteItems{get, post, put, patch, head, opt, del, conn, trace}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// StaticFile 	-> 指定路由到某个具体的磁盘文件
// Static 		-> 指定URL映射到某个磁盘目录，不打印当前路径下的文件列表
// StaticFS 	-> 和Static一样，但是显示当前目录文件列表（类似FTP），需要自定义磁盘路径http.FileSystem

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (gp *RouteGroup) StaticFile(relPath, filepath string) *RouteGroup {
	if strings.ContainsAny(relPath, ":*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *Context) {
		c.File(filepath)
	}
	gp.Get(relPath, handler)
	gp.Head(relPath, handler)
	return gp
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func (gp *RouteGroup) Static(relPath, root string) *RouteGroup {
	return gp.StaticFS(relPath, Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (gp *RouteGroup) StaticFS(relPath string, fs http.FileSystem) *RouteGroup {
	if strings.ContainsAny(relPath, ":*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := gp.createStaticHandler(relPath, fs)
	urlPattern := path.Join(relPath, "/*filepath")

	// Register GET and HEAD handlers
	gp.Get(urlPattern, handler)
	gp.Head(urlPattern, handler)
	return gp
}

func (gp *RouteGroup) createStaticHandler(relPath string, fs http.FileSystem) CtxHandler {
	absPath := gp.fixAbsolutePath(relPath)
	fileServer := http.StripPrefix(absPath, http.FileServer(fs))

	return func(c *Context) {
		if _, noListing := fs.(*onlyFilesFS); noListing {
			c.Res.WriteHeader(http.StatusNotFound)
		}

		file := c.UrlParam("filepath")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			// 没有匹配到静态文件，用系统中 404 （NoRoute handler）做响应处理
			c.Res.WriteHeader(http.StatusNotFound)
			c.route.ptrNode = gp.myApp.miniNode404
			c.execHandlers()
			return
		}
		_ = f.Close()
		fileServer.ServeHTTP(c.Res, c.Req.Raw)
	}
}

func (gp *RouteGroup) fixAbsolutePath(relPath string) string {
	return joinPaths(gp.prefix, relPath)
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type onlyFilesFS struct {
	fs http.FileSystem
}

type neuteredReaddirFile struct {
	http.File
}

// Dir returns a http.Filesystem that can be used by http.FileServer(). It is used internally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	}
	return &onlyFilesFS{fs}
}

// Open conforms to http.Filesystem.
func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

// Readdir overrides the http.File default implementation.
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
