package performance

import (
	"net/http"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(2)
	if differentReqNum > routersSum {
		panic("differentReqNum > routersSum, This is not allowed.")
	}
}

//var rtStrings []string
var reqPool []*http.Request
var routersLevel = 1                 // 路由数量的基数，实际值=routersSum
var routersSum = 1000 * routersLevel // 1000*routersNum
var middlewareNum = 10               // 中间件函数的数量
var reqPoolSize = routersSum         // 内置请求对象，用于模拟发起的不同Router请求
var differentReqNum = 1              // 用多少个不同路由的请求来测试

var firstSeg = [10]string{"first", "wang_zhi", "admin", "yes", "parents", "good", "boys", "health", "count", "right_now"}
var secondSeg = [10]string{"second", "xin", "kankan", "yes", "child", "tidy", "admin", "others", "parents", "testb"}
var thirdSeg = [10]string{"third", "chende", "shuo", "no", "finished", "xin", "kankan", "yes", "child", "tidy"}

// var thirdSeg = [1]string{":third"}
var forthSeg = [10]string{"third2", "chende2", "shuo2", "no2", "finished2", "xin2", "kankan2", "yes2", "child2", "tidy2"}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func benchRequest(b *testing.B, router http.Handler) {
	res := new(mockResponseWriter)
	//u := req.URL
	//rq := u.RawQuery
	//req.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	var req *http.Request
	for i := 0; i < b.N; i++ {
		req = reqPool[i%differentReqNum]
		//u.RawQuery = rq
		router.ServeHTTP(res, req)
	}
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
