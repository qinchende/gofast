package performance

import (
	"github.com/gin-gonic/gin"
	"io"
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
	ginAddRoutes(ginApp, routersLevel, ginHandle2)
}

func ginMiddlewareHandle(ctx *gin.Context) {
	ctx.Next()
}
func ginHandle2(ctx *gin.Context) {
}
func ginHandleTest(c *gin.Context) {
	io.WriteString(c.Writer, c.Request.RequestURI)
}
func ginHandleWrite(c *gin.Context) {
	io.WriteString(c.Writer, c.Params.ByName("name"))
}

// routeCt <= 10 && >= 1
func ginAddRoutes(gp *gin.Engine, routeCt int, hd gin.HandlerFunc) {
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
					gp.Handle(http.MethodGet, d, hd)
				}
			}
		}
	}
}

func ginAddMiddlewareHandlers(ginApp *gin.Engine, ct int) {
	for i := 0; i < ct; i++ {
		ginApp.Use(ginMiddlewareHandle)
	}
}

func BenchmarkGinWebRouter(b *testing.B) {
	benchRequest(b, ginApp)
}
