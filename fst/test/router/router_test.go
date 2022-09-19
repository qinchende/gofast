package router

import (
	"fmt"
	"github.com/qinchende/gofast/fst"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type header struct {
	Key   string
	Value string
}

func performRequest(app *fst.GoFast, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	return w
}

func testRouteOK(method string, t *testing.T) {
	passed := false
	passedAny := false

	router := fst.Default()
	router.All("/test2", func(c *fst.Context) {
		passedAny = true
	})
	router.Handle(method, "/test", func(c *fst.Context) {
		passed = true
	})
	router.BuildRoutes()

	w := performRequest(router, method, "/test")
	assert.True(t, passed)
	assert.Equal(t, http.StatusOK, w.Code)

	performRequest(router, method, "/test2")
	assert.True(t, passedAny)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK(method string, t *testing.T) {
	passed := false
	router := fst.Default()
	router.Handle(method, "/test_2", func(c *fst.Context) {
		passed = true
	})
	router.BuildRoutes()

	w := performRequest(router, method, "/test")
	assert.False(t, passed)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestSingleRouteOK tests that POST route is correctly invoked.
func testRouteNotOK2(method string, t *testing.T) {
	passed := false
	router := fst.Default()
	router.DefNotAllowedHandler = true
	var methodRoute string
	if method == http.MethodPost {
		methodRoute = http.MethodGet
	} else {
		methodRoute = http.MethodPost
	}
	router.Handle(methodRoute, "/test", func(c *fst.Context) {
		passed = true
	})
	router.BuildRoutes()

	w := performRequest(router, method, "/test")
	assert.False(t, passed)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func TestRouterMethod(t *testing.T) {
	router := fst.Default()
	router.Put("/hey2", func(c *fst.Context) {
		c.String(http.StatusOK, "sup2")
	})
	router.Put("/hey", func(c *fst.Context) {
		c.String(http.StatusOK, "called")
	})
	router.Put("/hey3", func(c *fst.Context) {
		c.String(http.StatusOK, "sup3")
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodPut, "/hey")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "called", w.Body.String())

	w = performRequest(router, http.MethodPut, "/hey3")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "sup3", w.Body.String())

	w = performRequest(router, http.MethodPut, "/hey2")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "sup2", w.Body.String())
}

func TestRouterGroupRouteOK(t *testing.T) {
	testRouteOK(http.MethodGet, t)
	testRouteOK(http.MethodPost, t)
	testRouteOK(http.MethodPut, t)
	testRouteOK(http.MethodPatch, t)
	testRouteOK(http.MethodHead, t)
	testRouteOK(http.MethodOptions, t)
	testRouteOK(http.MethodDelete, t)
	testRouteOK(http.MethodConnect, t)
	testRouteOK(http.MethodTrace, t)
}

func TestRouteNotOK(t *testing.T) {
	testRouteNotOK(http.MethodGet, t)
	testRouteNotOK(http.MethodPost, t)
	testRouteNotOK(http.MethodPut, t)
	testRouteNotOK(http.MethodPatch, t)
	testRouteNotOK(http.MethodHead, t)
	testRouteNotOK(http.MethodOptions, t)
	testRouteNotOK(http.MethodDelete, t)
	testRouteNotOK(http.MethodConnect, t)
	testRouteNotOK(http.MethodTrace, t)
}

func TestRouteNotOK2(t *testing.T) {
	testRouteNotOK2(http.MethodGet, t)
	testRouteNotOK2(http.MethodPost, t)
	testRouteNotOK2(http.MethodPut, t)
	testRouteNotOK2(http.MethodPatch, t)
	testRouteNotOK2(http.MethodHead, t)
	testRouteNotOK2(http.MethodOptions, t)
	testRouteNotOK2(http.MethodDelete, t)
	testRouteNotOK2(http.MethodConnect, t)
	testRouteNotOK2(http.MethodTrace, t)
}

func TestRouteRedirectTrailingSlash(t *testing.T) {
	router := fst.Default()
	router.RedirectTrailingSlash = true
	router.Get("", func(c *fst.Context) { c.String(200, "home") })
	router.Get("/path", func(c *fst.Context) { c.String(200, "non slash") })
	router.Get("/path2/", func(c *fst.Context) {})
	router.Get("/path2others", func(c *fst.Context) {})
	router.Post("/path3", func(c *fst.Context) {})
	router.Put("/path4/", func(c *fst.Context) {})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "home", w.Body.String())

	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "non slash", w.Body.String())

	w = performRequest(router, http.MethodGet, "/path2")
	assert.Equal(t, "/path2/", w.Header().Get("Location"))
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	w = performRequest(router, http.MethodPost, "/path3/")
	assert.Equal(t, "/path3", w.Header().Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	w = performRequest(router, http.MethodPut, "/path4")
	assert.Equal(t, "/path4/", w.Header().Get("Location"))
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodGet, "/path2/")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodPost, "/path3")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodPut, "/path4/")
	assert.Equal(t, http.StatusOK, w.Code)

	w = performRequest(router, http.MethodGet, "/path2", header{Key: "X-Forwarded-Prefix", Value: "/api"})
	assert.Equal(t, "/api/path2/", w.Header().Get("Location"))
	assert.Equal(t, 301, w.Code)

	w = performRequest(router, http.MethodGet, "/path2/", header{Key: "X-Forwarded-Prefix", Value: "/api/"})
	assert.Equal(t, 200, w.Code)

	router.RedirectTrailingSlash = false

	w = performRequest(router, http.MethodGet, "/path/")
	assert.Equal(t, http.StatusNotFound, w.Code)
	w = performRequest(router, http.MethodGet, "/path2")
	assert.Equal(t, http.StatusNotFound, w.Code)
	w = performRequest(router, http.MethodPost, "/path3/")
	assert.Equal(t, http.StatusNotFound, w.Code)
	w = performRequest(router, http.MethodPut, "/path4")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestContextParamsGet tests that a parameter can be parsed from the URL.
func TestRouteParamsByName(t *testing.T) {
	name := ""
	lastName := ""
	wild := ""
	router := fst.Default()
	router.Get("/test/:name/:last_name/*wild", func(c *fst.Context) {
		name = c.UrlParams.ByName("name")
		lastName = c.UrlParams.ByName("last_name")
		var ok bool
		wild, ok = c.UrlParams.Get("wild")

		assert.True(t, ok)
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, lastName, c.Param("last_name"))

		assert.Empty(t, c.Param("wtf"))
		assert.Empty(t, c.UrlParams.ByName("wtf"))

		wtf, ok := c.UrlParams.Get("wtf")
		assert.Empty(t, wtf)
		assert.False(t, ok)
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/test/john/smith/is/super/great")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "john", name)
	assert.Equal(t, "smith", lastName)
	assert.Equal(t, "is/super/great", wild)
}

// TestContextParamsGet tests that a parameter can be parsed from the URL even with extra slashes.
func TestRouteParamsByNameWithExtraSlash(t *testing.T) {
	name := ""
	lastName := ""
	wild := ""
	router := fst.Default()
	router.RemoveExtraSlash = true
	router.Get("/test/:name/:last_name/*wild", func(c *fst.Context) {
		name = c.UrlParams.ByName("name")
		lastName = c.UrlParams.ByName("last_name")
		var ok bool
		wild, ok = c.UrlParams.Get("wild")

		assert.True(t, ok)
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, lastName, c.Param("last_name"))

		assert.Empty(t, c.Param("wtf"))
		assert.Empty(t, c.UrlParams.ByName("wtf"))

		wtf, ok := c.UrlParams.Get("wtf")
		assert.Empty(t, wtf)
		assert.False(t, ok)
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "//test//john//smith//is//super//great")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "john", name)
	assert.Equal(t, "smith", lastName)
	assert.Equal(t, "is/super/great", wild)
}

// TestHandleStaticFile - ensure the static file handles properly
func TestRouteStaticFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := ioutil.TempFile(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString("Gin Web Framework")
	assert.NoError(t, err)
	_ = f.Close()

	dir, filename := filepath.Split(f.Name())

	// SETUP gin
	router := fst.Default()
	router.Static("/using_static", dir)
	router.StaticFile("/result", f.Name())
	router.BuildRoutes()

	w1 := performRequest(router, http.MethodGet, "/using_static/"+filename)
	w2 := performRequest(router, http.MethodGet, "/result")

	assert.Equal(t, w1, w2)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "Gin Web Framework", w1.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w1.Header().Get("Content-Type"))

	w3 := performRequest(router, http.MethodHead, "/using_static/"+filename)
	w4 := performRequest(router, http.MethodHead, "/result")

	assert.Equal(t, w3, w4)
	assert.Equal(t, http.StatusOK, w3.Code)
}

// TestHandleStaticDir - ensure the root/sub dir handles properly
func TestRouteStaticListingDir(t *testing.T) {
	router := fst.Default()
	router.StaticFS("/", fst.Dir("./", true))
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "helper_test.go")
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

// TestHandleHeadToDir - ensure the root/sub dir handles properly
func TestRouteStaticNoListing(t *testing.T) {
	router := fst.Default()
	router.Static("/", "./")
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.NotContains(t, w.Body.String(), "helper_test.go")
}

func TestRouterMiddlewareAndStatic(t *testing.T) {
	router := fst.Default()
	static := router.Group("/")
	static.Before(func(c *fst.Context) {
		c.ResWrap.Header().Add("Last-Modified", "Mon, 02 Jan 2006 15:04:05 MST")
		c.ResWrap.Header().Add("Expires", "Mon, 02 Jan 2006 15:04:05 MST")
		c.ResWrap.Header().Add("X-GIN", "GoFast Framework")
	})
	static.Static("/", "./")
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/helper_test.go")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "package router")
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.NotEqual(t, w.Header().Get("Last-Modified"), "Mon, 02 Jan 2006 15:04:05 MST")
	assert.Equal(t, "Mon, 02 Jan 2006 15:04:05 MST", w.Header().Get("Expires"))
	assert.Equal(t, "GoFast Framework", w.Header().Get("x-GIN"))
}

func TestRouteNotAllowedEnabled(t *testing.T) {
	router := fst.Default()
	router.DefNotAllowedHandler = true
	router.Post("/path", func(c *fst.Context) {})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	router.NoMethod(func(c *fst.Context) {
		c.String(http.StatusTeapot, "responseText")
	})
	router.BuildRoutes()
	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, "responseText", w.Body.String())
	assert.Equal(t, http.StatusTeapot, w.Code)
}

func TestRouteNotAllowedEnabled2(t *testing.T) {
	router := fst.Default()
	router.DefNotAllowedHandler = true
	// add one methodTree to trees
	router.Handle(http.MethodPost, "/", func(_ *fst.Context) {})
	router.Get("/path2", func(c *fst.Context) {})
	router.BuildRoutes()

	w := performRequest(router, http.MethodPost, "/path2")
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestRouteNotAllowedDisabled(t *testing.T) {
	router := fst.Default()
	router.DefNotAllowedHandler = false
	router.Post("/path", func(c *fst.Context) {})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, http.StatusNotFound, w.Code)

	router.NoMethod(func(c *fst.Context) {
		c.String(http.StatusTeapot, "responseText")
	})
	router.BuildRoutes()

	w = performRequest(router, http.MethodGet, "/path")
	assert.Equal(t, "404 (page not found)", w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterNotFoundWithRemoveExtraSlash(t *testing.T) {
	router := fst.Default()
	router.RemoveExtraSlash = true
	router.Get("/path", func(c *fst.Context) {})
	router.Get("/", func(c *fst.Context) {})

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/../path", http.StatusOK, ""},    // CleanPath
		{"/nope", http.StatusNotFound, ""}, // NotFound
	}
	router.BuildRoutes()
	for _, tr := range testRoutes {
		w := performRequest(router, "GET", tr.route)
		assert.Equal(t, tr.code, w.Code)
		if w.Code != http.StatusNotFound {
			assert.Equal(t, tr.location, fmt.Sprint(w.Header().Get("Location")))
		}
	}
}

func TestRouterNotFound(t *testing.T) {
	router := fst.Default()
	router.RedirectTrailingSlash = true

	router.Get("/path", func(c *fst.Context) {})
	router.Get("/dir/", func(c *fst.Context) {})
	router.Get("/", func(c *fst.Context) {})

	testRoutes := []struct {
		route    string
		code     int
		location string
	}{
		{"/path/", http.StatusMovedPermanently, "/path"}, // TSR -/
		{"/dir", http.StatusMovedPermanently, "/dir/"},   // TSR +/
		{"/nope", http.StatusNotFound, ""},               // NotFound
	}
	router.BuildRoutes()
	for _, item := range testRoutes {
		w := performRequest(router, http.MethodGet, item.route)
		assert.Equal(t, item.code, w.Code)
		if w.Code != http.StatusNotFound {
			assert.Equal(t, item.location, fmt.Sprint(w.Header().Get("Location")))
		}
	}

	// Test custom not found handler
	var notFound bool
	router.NoRoute(func(c *fst.Context) {
		//c.AbortAndRender(http.StatusNotFound)
		notFound = true
	})
	router.BuildRoutes()
	w := performRequest(router, http.MethodGet, "/nope")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.True(t, notFound)

	// Test other method than GET (want 307 instead of 301)
	router.Patch("/path", func(c *fst.Context) {})
	router.BuildRoutes()
	w = performRequest(router, http.MethodPatch, "/path/")
	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "map[Location:[/path]]", fmt.Sprint(w.Header()))

	// Test special case where no node for the prefix "/" exists
	router = fst.Default()
	router.BuildRoutes()
	router.Get("/a", func(c *fst.Context) {})
	w = performRequest(router, http.MethodGet, "/")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRouterStaticFSNotFound(t *testing.T) {
	router := fst.Default()
	//router.SetMode(fst.DebugMode)

	router.StaticFS("/", http.FileSystem(http.Dir("/thisreallydoesntexist/")))
	router.NoRoute(func(c *fst.Context) {
		c.String(404, "non existent")
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/nonexistent")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "non existent", w.Body.String())

	w = performRequest(router, http.MethodHead, "/nonexistent")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "non existent", w.Body.String())
}

func TestRouterStaticFSFileNotFound(t *testing.T) {
	router := fst.Default()
	router.StaticFS("/", http.FileSystem(http.Dir(".")))
	router.BuildRoutes()

	assert.NotPanics(t, func() {
		performRequest(router, http.MethodGet, "/nonexistent")
	})
}

// Reproduction test for the bug of issue #1805
func TestMiddlewareCalledOnceByRouterStaticFSNotFound(t *testing.T) {
	router := fst.Default()

	// Middleware must be called just only once by per request.
	middlewareCalledNum := 0
	router.Before(func(c *fst.Context) {
		middlewareCalledNum++
	})

	router.StaticFS("/", http.FileSystem(http.Dir("/thisreallydoesntexist/")))
	router.BuildRoutes()

	// First access
	performRequest(router, http.MethodGet, "/nonexistent")
	assert.Equal(t, 1, middlewareCalledNum)

	// Second access
	performRequest(router, http.MethodHead, "/nonexistent")
	assert.Equal(t, 2, middlewareCalledNum)
}

func TestRouteRawPath(t *testing.T) {
	router := fst.Default()
	router.UseRawPath = true
	router.UnescapePathValues = true

	router.Post("/project/:name/build/:num", func(c *fst.Context) {
		name := c.UrlParams.ByName("name")
		num := c.UrlParams.ByName("num")

		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, num, c.Param("num"))

		assert.Equal(t, "Some/Other/Project", name)
		assert.Equal(t, "222", num)
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodPost, "/project/Some%2FOther%2FProject/build/222")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouteRawPathNoUnescape(t *testing.T) {
	router := fst.Default()
	router.UseRawPath = true
	router.UnescapePathValues = false
	//router.DisableDefNoRoute = true

	router.Post("/project/:name/build/:num", func(c *fst.Context) {
		name := c.UrlParams.ByName("name")
		num := c.UrlParams.ByName("num")

		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, num, c.Param("num"))

		assert.Equal(t, "Some%2FOther%2FProject", name)
		assert.Equal(t, "333", num)
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodPost, "/project/Some%2FOther%2FProject/build/333")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouteServeErrorWithWriteHeader(t *testing.T) {
	router := fst.Default()
	router.Before(func(c *fst.Context) {
		//c.Status(421)
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/NotFound")
	assert.Equal(t, 421, w.Code)
	assert.Equal(t, 0, w.Body.Len())
}

func TestRouteContextHoldsFullPath(t *testing.T) {
	router := fst.Default()

	// Test routes
	routes := []string{
		"/simple",
		"/project/:name",
		"/",
		"/news/home",
		"/news",
		"/simple-two/one",
		"/simple-two/one-two",
		"/project/:name/build/*params",
		"/project/:name/bui",
		"/user/:id/status",
		"/user/:id",
		"/user/:id/",
		"/user/:id/profile",
	}

	for _, rt := range routes {
		actualRoute := rt
		router.Get(rt, func(c *fst.Context) {
			// For each defined route context should contain its full path
			assert.Equal(t, actualRoute, c.FullPath())
			c.String(http.StatusOK, "")
		})
	}
	router.BuildRoutes()

	for _, rt := range routes {
		w := performRequest(router, http.MethodGet, rt)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// ++++++++++++++++++++++++++++++++++++++++++
	// Test not found
	router.Before(func(c *fst.Context) {
		// For not found routes full path is empty
		assert.Equal(t, "", c.FullPath())
	})
	router.BuildRoutes()

	w := performRequest(router, http.MethodGet, "/not-found")
	assert.Equal(t, http.StatusNotFound, w.Code)
}
