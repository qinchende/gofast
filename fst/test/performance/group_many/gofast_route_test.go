package group_many

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
	"testing"
	"time"
)

func init() {
	initGoFastServer()
}

var myWeb *fst.GoFast

func initGoFastServer() {
	// 新建Server
	myWeb = fst.CreateServer(&fst.GfConfig{
		RunMode: fst.ProductMode,
	})

	gftAddMiddlewareHandlers(middlewareNum)
	addRoutes(routersLevel, func(url string) {
		myWeb.Handle(http.MethodGet, url, gftHandle2)
	})
	myWeb.BuildRoutes()
}

func gftMiddlewareHandle(ctx *fst.Context) int {
	// 请求前获取当前时间
	nowTime := time.Now()

	arr := [10000]int{}
	ctLen := len(arr)
	for i := 0; i < ctLen; i++ {
		arr[i] = i * 10
	}

	//return arr[0]
	return int(time.Since(nowTime))
}

func gftHandle2(_ *fst.Context) {
	//print(1)
}

//func gftHandleTest(c *fst.Context) {
//	io.WriteString(c.Res, c.Req.RequestURI)
//}
//func gftHandleWrite(c *fst.Context) {
//	io.WriteString(c.Res, c.UrlParams.Value("name"))
//}

// add GoFast middlewares
func gftAddMiddlewareHandlers(ct int) {
	for i := 0; i < ct; i++ {
		myWeb.Before(func(context *fst.Context) {
			gftMiddlewareHandle(context)
		})
	}
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, myWeb)
}
