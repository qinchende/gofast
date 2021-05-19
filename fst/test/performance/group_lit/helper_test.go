package group_lit

import (
	"net/http"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(2)
}

var reqPool []*http.Request // 模拟请求的对象数组（伪造并缓存请求对象）
var middlewareNum = 10      // 中间件函数的数量

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func benchRequest(b *testing.B, router http.Handler) {
	res := new(mockResponseWriter)
	b.ReportAllocs()
	b.ResetTimer()

	// 并发测试模式
	b.SetParallelism(110000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			router.ServeHTTP(res, reqPool[0])
		}
	})
}

//func benchRequest(b *testing.B, router http.Handler) {
//	res := new(mockResponseWriter)
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		router.ServeHTTP(res, reqPool[0])
//	}
//}

type regRouteFun func(url string)

// routeCt <= 10 && >= 1
func addRoutes(regRoute regRouteFun) {
	reqPool = make([]*http.Request, 0, 1)

	d := "/first/second/third"
	req, _ := http.NewRequest("GET", d, nil)
	reqPool = append(reqPool, req)

	regRoute(d)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}
