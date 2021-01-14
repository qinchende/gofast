package test

import (
	"testing"
)

// 单元测试案例
//func fib(n int) int {
//	if n == 0 || n == 1 {
//		return n
//	}
//	return fib(n-2) + fib(n-1)
//}

//func fullName(first string, second string) string {
//	return fmt.Sprintf("%s %s", first, second)
//}
//
//func TestFullName(t *testing.T) {
//	full := fullName("chen", "de")
//	assert.Equal(t, "chen de", full)
//}


func TestCover(t *testing.T) {
	Cover(1)
	Cover(2)
	Cover(3)
}

//
//// 模拟网络请求
//func execReq(app *fst.GoFast, method, path string) *httptest.ResponseRecorder {
//	app.ReadyToListen()
//	req := httptest.NewRequest(method, path, nil)
//	resW := httptest.NewRecorder()
//	app.ServeHTTP(resW, req)
//	return resW
//}
//
//func TestRouterMethod(t *testing.T) {
//	router := fst.Default()
//	router.Put("/hey", func(c *fst.Context) {
//		c.String(http.StatusOK, "chen de")
//	})
//	w := execReq(router, http.MethodPut, "/hey")
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Equal(t, "chen de", w.Body.String())
//}
//
//// 模拟网络请求2
//func newTstServer() *httptest.Server {
//	handler := func (rw http.ResponseWriter, r *http.Request) {
//		u := struct {
//			Name string
//		}{
//			Name: "闪电侠",
//		}
//
//		rw.Header().Set("Content-Type", "application/json")
//		rw.WriteHeader(http.StatusOK)
//		_ = json.NewEncoder(rw).Encode(u)
//	}
//	return httptest.NewServer(http.HandlerFunc(handler))
//}
//
//func TestSendJSONData(t *testing.T) {
//	server := newTstServer()
//	defer server.Close()
//
//	res, err := http.Get(server.URL)
//	if err != nil {
//		t.Fatal("Get请求失败")
//	}
//	defer res.Body.Close()
//
//	log.Println("Status code: ", res.StatusCode)
//	jsonData, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	log.Printf("Body: %s", jsonData)
//}
