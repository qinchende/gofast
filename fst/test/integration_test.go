package test

import (
	"gofast/fst"
	"io"
	"net/http"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(1)
	initServer()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func fstHandle(_ *fst.Context) {}
func fstHandleTest(c *fst.Context) {
	io.WriteString(c.Reply, c.Request.RequestURI)
}
func fstHandleWrite(c *fst.Context) {
	io.WriteString(c.Reply, c.Params.ByName("name"))
}

var firstSeg = [10]string{"first", "wang_zhi", "admin", "yes", "parents", "good", "boys", "health", "count", "right_now"}
var secondSeg = [10]string{"second", "xin", "kankan", "yes", "child", "tidy", "admin", "others", "parents", "testb"}
//var thirdSeg = [10]string{"third", "chende", "shuo", "no", "finished", "xin", "kankan", "yes", "child", "tidy"}
var thirdSeg = [1]string{":third"}
var forthSeg = [10]string{"third2", "chende2", "shuo2", "no2", "finished2", "xin2", "kankan2", "yes2", "child2", "tidy2"}

// routeCt <= 10 && >= 1
func addRoutes(gp *fst.HomeSite, routeCt int, hd fst.CtxHandler) {
	var a, b, c, d string
	for i := 0; i < routeCt; i++ {
		a = "/" + firstSeg[i]
		for j := 0; j < len(secondSeg); j++ {
			b = a + "/" + secondSeg[j]
			for k := 0; k < len(thirdSeg); k++ {
				c = b + "/" + thirdSeg[k]
				for n := 0; n < len(forthSeg); n++ {
					d = c + "/" + forthSeg[n]
					gp.Method(http.MethodGet, d, hd)
				}
			}
		}
	}
}

func addUseHandlers(home *fst.HomeSite, ct int) {
	for i := 0; i < ct; i++ {
		home.Before(fstHandle)
	}
}

var gft *fst.GoFast

func initServer() {
	// 新建Server
	app2, home := fst.CreateServer(&fst.AppConfig{
		RunMode: fst.ProductMode,
	})
	gft = app2
	addUseHandlers(home, 10)
	addRoutes(home, 10, fstHandle)
	gft.ReadyToListen()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func BenchmarkGoFastWeb(b *testing.B) {
	r, _ := http.NewRequest("GET", "/first/xin/no/chende2", nil)
	benchRequest(b, gft, r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func benchRequest(b *testing.B, router http.Handler, req *http.Request) {
	res := new(mockResponseWriter)
	u := req.URL
	rq := u.RawQuery
	req.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(res, req)
	}
}

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
