package fst

import "net/http"

// 请求从net/http包传出来之后，需要在框架中转换成我们自己的Request对象
type Request struct {
	rawReq *http.Request
	gftApp *GoFast
	fitIdx int
}
