package group_lit

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
	"testing"
)

func init() {
	initGoFastServer()
}

var gftApp *fst.GoFast

func initGoFastServer() {
	// 新建Server
	gftApp = fst.CreateServer(&fst.AppConfig{
		RunMode: fst.ProductMode,
	})

	gftAddMiddlewareHandlers(middlewareNum)
	addRoutes(func(url string) {
		gftApp.Handle(http.MethodGet, url, gftHandle2)
	})
	gftApp.BuildRouters()
}

func gftHandle2(c *fst.Context) {
	//println(unsafe.Sizeof(*c))
}

// add GoFast middlewares
func gftAddMiddlewareHandlers(ct int) {
	for i := 0; i < ct; i++ {
		gftApp.Before(func(context *fst.Context) {
			gftMiddlewareHandle(context)
		})
	}
}

func gftMiddlewareHandle(ctx *fst.Context) {
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, gftApp)
}
