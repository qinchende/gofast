package performance

import (
	"gofast/fst"
	"io"
	"net/http"
	"testing"
)

func init() {
	initGoFastServer()
}

var gftApp *fst.GoFast

func initGoFastServer() {
	// 新建Server
	app, home := fst.CreateServer(&fst.AppConfig{
		RunMode: fst.ProductMode,
	})
	gftApp = app

	gftAddMiddlewareHandlers(home, middlewareNum)
	gftAddRoutes(home, routersLevel, gftHandle2)
	gftApp.ReadyToListen()
}

func gftMiddlewareHandle(ctx *fst.Context) {
}
func gftHandle2(_ *fst.Context) {
}
func gftHandleTest(c *fst.Context) {
	io.WriteString(c.Reply, c.Request.RequestURI)
}
func gftHandleWrite(c *fst.Context) {
	io.WriteString(c.Reply, c.Params.ByName("name"))
}

// routeCt <= 10 && >= 1
func gftAddRoutes(gp *fst.HomeSite, routeCt int, hd fst.CtxHandler) {
	//rtStrings = make([]string, 0 , reqPoolSize)
	reqPool = make([]*http.Request, 0, reqPoolSize)

	var a, b, c, d string
	for i := 0; i < routeCt; i++ {
		a = "/" + firstSeg[i]
		for j := 0; j < len(secondSeg); j++ {
			b = a + "/" + secondSeg[j]
			for k := 0; k < len(thirdSeg); k++ {
				c = b + "/" + thirdSeg[k]
				for n := 0; n < len(forthSeg); n++ {
					d = c + "/" + forthSeg[n]
					//rtStrings = append(rtStrings, d)
					r, _ := http.NewRequest("GET", d, nil)
					reqPool = append(reqPool, r)
					gp.Method(http.MethodGet, d, hd)
				}
			}
		}
	}
}

func gftAddMiddlewareHandlers(home *fst.HomeSite, ct int) {
	for i := 0; i < ct; i++ {
		home.Before(gftMiddlewareHandle)
	}
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, gftApp)
}
