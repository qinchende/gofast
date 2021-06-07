package group_lit

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func init() {
	initGoFastServer()
}

var gftApp *fst.GoFast

func initGoFastServer() {
	gftApp = fst.Default()

	// 添加一定数量的中间件函数
	for i := 0; i < middlewareNum; i++ {
		gftApp.Before(func(ctx *fst.Context) {})
	}
	addRoutes(func(url string) {
		gftApp.Handle(http.MethodGet, url, func(c *fst.Context) {})
	})
	gftApp.BuildRouters()
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, gftApp)
}

func TestRequestGoFast(t *testing.T) {
	w := test.ExecRequest(gftApp, http.MethodGet, reqPool[0].URL.Path)
	assert.Equal(t, http.StatusOK, w.Code)
}
