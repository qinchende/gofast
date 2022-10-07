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

var myApp *fst.GoFast

func initGoFastServer() {
	myApp = fst.Default()

	// 添加一定数量的中间件函数
	for i := 0; i < middlewareNum; i++ {
		myApp.Before(func(ctx *fst.Context) {})
	}
	addRoutes(func(url string) {
		myApp.Handle(http.MethodGet, url, func(c *fst.Context) {})
	})
	myApp.BuildRoutes()
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, myApp)
}

func TestRequestGoFast(t *testing.T) {
	w := test.ExecRequest(myApp, http.MethodGet, reqPool[0].URL.Path)
	assert.Equal(t, http.StatusOK, w.Code)
}
