package performance

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
	addRoutes(routersLevel, func(url string) {
		gftApp.Handle(http.MethodGet, url, gftHandle2)
	})
	gftApp.ReadyToListen()
}

func gftMiddlewareHandle(ctx *fst.Context) {
}
func gftHandle2(_ *fst.Context) {
}

//func gftHandleTest(c *fst.Context) {
//	io.WriteString(c.Reply, c.Request.RequestURI)
//}
//func gftHandleWrite(c *fst.Context) {
//	io.WriteString(c.Reply, c.Params.ByName("name"))
//}

func gftAddMiddlewareHandlers(ct int) {
	for i := 0; i < ct; i++ {
		gftApp.Before(gftMiddlewareHandle)
	}
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, gftApp)
}
