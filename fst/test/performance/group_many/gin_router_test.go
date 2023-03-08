package group_many

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"time"
)

func init() {
	initGinServer()
}

var ginApp *gin.Engine

func initGinServer() {
	gin.SetMode(gin.ReleaseMode)

	ginApp = gin.New()
	ginAddMiddlewareHandlers(ginApp, middlewareNum)
	addRoutes(routersLevel, func(url string) {
		ginApp.Handle(http.MethodGet, url, ginHandle2)
	})
}

func ginMiddlewareHandle(ctx *gin.Context) int {
	// 请求前获取当前时间
	nowTime := time.Now()

	arr := [10000]int{}
	ctLen := len(arr)
	for i := 0; i < ctLen; i++ {
		arr[i] = i * 10
	}

	ctx.Next()

	//time := time.Since(nowTime)
	return int(time.Since(nowTime))
	//return arr[0]
}
func ginHandle2(ctx *gin.Context) {
}

//func ginHandleTest(c *gin.Context) {
//	io.WriteString(c.Writer, c.Req.RequestURI)
//}
//func ginHandleWrite(c *gin.Context) {
//	io.WriteString(c.Writer, c.Params.Value("name"))
//}

// add gin middlewares
func ginAddMiddlewareHandlers(ginApp *gin.Engine, ct int) {
	for i := 0; i < ct; i++ {
		ginApp.Use(func(context *gin.Context) {
			ginMiddlewareHandle(context)
		})
	}
}

// start benchmark
func BenchmarkGinWebRouter(b *testing.B) {
	benchRequest(b, ginApp)
}
