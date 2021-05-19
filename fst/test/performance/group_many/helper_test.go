package group_many

import (
	"net/http"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(4)
	if differentReqNum > routersSum {
		panic("differentReqNum > routersSum, This is not allowed.")
	}
}

// 调节这些参数，来模拟不同路由场景下，Gin和GoFast的性能
// var rtStrings []string
var reqPool []*http.Request          // 模拟请求的对象数组（伪造并缓存请求对象）
var routersLevel = 10                // 路由数量的基数，实际值=routersSum
var routersSum = 1000 * routersLevel // 1000 * routersNum
var middlewareNum = 2                // 中间件函数的数量
var reqPoolSize = routersSum         // 内置请求对象，用于模拟发起的不同Router请求
var differentReqNum = 1000           // 用多少个不同路由的请求来测试

// 模拟四段组成的 Url 路由
var (
	firstSeg  = [10]string{"first", "wang_zhi", "admin", "yes", "parents", "good", "boys", "health", "count", "right_now"}
	secondSeg = [10]string{"second", "xin", "kankan", "yes", "child", "tidy", "admin", "others", "parents", "testb"}
	thirdSeg  = [10]string{"third", "chende", "shuo", "no", "finished", "xin", "kankan", "yes", "child", "tidy"}
	forthSeg  = [10]string{"third2", "chende2", "shuo2", "no2", "finished2", "xin2", "kankan2", "yes2", "child2", "tidy2"}
	// var thirdSeg = [1]string{":third"}
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func benchRequest(b *testing.B, router http.Handler) {
	res := new(mockResponseWriter)
	//u := req.URL
	//rq := u.RawQuery
	//req.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	// 并发测试模式
	b.SetParallelism(20000)
	b.RunParallel(func(pb *testing.PB) {
		var req *http.Request
		i := -1
		for pb.Next() {
			i++
			req = reqPool[i%differentReqNum]
			router.ServeHTTP(res, req)
		}
	})

	//var req *http.Request
	//for i := 0; i < b.N; i++ {
	//	req = reqPool[i%differentReqNum]
	//	//u.RawQuery = rq
	//	router.ServeHTTP(res, req)
	//}
}

type regRouteFun func(url string)

// routeCt <= 10 && >= 1
func addRoutes(routeCt int, regRoute regRouteFun) {
	// rtStrings = make([]string, 0 , reqPoolSize)
	reqPool = make([]*http.Request, 0, reqPoolSize)

	var a, b, c, d string
	for i := 0; i < routeCt; i++ {
		a = "/" + firstSeg[i]
		for j := 0; j < len(secondSeg); j++ {
			b = a + "/" + secondSeg[j]
			for k := 0; k < len(thirdSeg); k++ {
				c = b + "/" + thirdSeg[k]
				for n := 0; n < len(forthSeg); n++ {
					// 模拟请求的 url
					d = c + "/" + forthSeg[n]
					// rtStrings = append(rtStrings, d)
					req, _ := http.NewRequest("GET", d, nil)
					reqPool = append(reqPool, req)

					// 注册路由
					regRoute(d)
				}
			}
		}
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
