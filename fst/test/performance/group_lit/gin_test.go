package group_lit

import (
	"github.com/gin-gonic/gin"
	"github.com/qinchende/gofast/fst/test"
	"github.com/stretchr/testify/assert"
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

	// 添加一定数量的中间件函数
	for i := 0; i < middlewareNum; i++ {
		ginApp.Use(func(ctx *gin.Context) {})
	}

	addRoutes(func(url string) {
		ginApp.Handle(http.MethodGet, url, func(c *gin.Context) {})
	})
}

func BenchmarkGinWebRouter(b *testing.B) {
	benchRequest(b, ginApp)
}

func TestRequestGin(t *testing.T) {
	w := test.ExecRequest(ginApp, http.MethodGet, reqPool[0].URL.Path)
	assert.Equal(t, http.StatusOK, w.Code)
}
