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

var myWeb *fst.GoFast

func initGoFastServer() {
	myWeb = fst.Default()

	// 添加一定数量的中间件函数
	for i := 0; i < middlewareNum; i++ {
		myWeb.Before(func(ctx *fst.Context) {})
	}
	addRoutes(func(url string) {
		myWeb.Handle(http.MethodGet, url, func(c *fst.Context) {})
	})
	myWeb.BuildRoutes()
}

func BenchmarkGoFastWebRouter(b *testing.B) {
	benchRequest(b, myWeb)
}

func TestRequestGoFast(t *testing.T) {
	w := test.ExecRequest(myWeb, http.MethodGet, reqPool[0].URL.Path)
	assert.Equal(t, http.StatusOK, w.Code)
}
