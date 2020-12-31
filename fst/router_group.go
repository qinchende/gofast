package fst

import (
	"net/http"
	"path"
	"strings"
)

func (gp *RouterGroup) AddGroup(relPath string, handlers ...CHandler) *RouterGroup {
	gpNew := &RouterGroup{
		prefix: gp.fixAbsolutePath(relPath),
		faster: gp.faster,
	}
	gp.children = append(gp.children, gpNew)
	return gpNew
}

// Prefix returns the base path of router gp.
// For example, if v := router.Group("/rest/n/v1/api"), v.Prefix() is "/rest/n/v1/api".
func (gp *RouterGroup) Prefix() string {
	return gp.prefix
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (gp *RouterGroup) StaticFile(relPath, filepath string) *RouterGroup {
	if strings.Contains(relPath, ":") || strings.Contains(relPath, "*") {
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
func (gp *RouterGroup) Static(relPath, root string) *RouterGroup {
	return gp.StaticFS(relPath, Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (gp *RouterGroup) StaticFS(relPath string, fs http.FileSystem) *RouterGroup {
	if strings.Contains(relPath, ":") || strings.Contains(relPath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := gp.createStaticHandler(relPath, fs)
	urlPattern := path.Join(relPath, "/*filepath")

	// Register GET and HEAD handlers
	gp.Get(urlPattern, handler)
	gp.Head(urlPattern, handler)
	return gp
}

func (gp *RouterGroup) createStaticHandler(relPath string, fs http.FileSystem) CHandler {
	absPath := gp.fixAbsolutePath(relPath)
	fileServer := http.StripPrefix(absPath, http.FileServer(fs))

	return func(c *Context) {
		if _, noListing := fs.(*onlyFilesFS); noListing {
			c.Reply.WriteHeader(http.StatusNotFound)
		}

		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Reply.WriteHeader(http.StatusNotFound)
			//c.handlers = gp.faster.noRoute
			// Reset index
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Reply, c.Request)
	}
}

func (gp *RouterGroup) fixAbsolutePath(relPath string) string {
	return joinPaths(gp.prefix, relPath)
}

//
//func (ft *Faster) rebuild404Handlers() {
//	ft.allNoRoute = ft.combineHandlers(ft.noRoute)
//}
//
//func (ft *Faster) rebuild405Handlers() {
//	ft.allNoMethod = ft.combineHandlers(ft.noMethod)
//}
