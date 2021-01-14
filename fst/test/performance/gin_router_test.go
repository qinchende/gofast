package performance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
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

func ginMiddlewareHandle(ctx *gin.Context) {
	//ctx.Next()
}
func ginHandle2(ctx *gin.Context) {
}

//func ginHandleTest(c *gin.Context) {
//	io.WriteString(c.Writer, c.Request.RequestURI)
//}
//func ginHandleWrite(c *gin.Context) {
//	io.WriteString(c.Writer, c.Params.ByName("name"))
//}

func ginAddMiddlewareHandlers(ginApp *gin.Engine, ct int) {
	for i := 0; i < ct; i++ {
		ginApp.Use(ginMiddlewareHandle)
	}
}

func BenchmarkGinWebRouter(b *testing.B) {
	benchRequest(b, ginApp)
}
