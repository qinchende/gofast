// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package httpx

import (
	"github.com/qinchende/gofast/core/cst"
	"net/http"
)

type HttpClient struct {
	http.Client
}

type RequestPet struct {
	ProxyUrl   string    // 代理服务器地址
	Method     string    // GET or POST
	Url        string    // http(s)地址
	Headers    cst.WebKV // 请求头
	QueryArgs  cst.WebKV // url上的参数
	BodyArgs   cst.WebKV // body带的参数
	BodyFormat int8      // body数据的格式，比如 json|url-encoding|xml
}

var myClient = &HttpClient{} // 框架默认client对象

func ClientByPet(pet *RequestPet) *HttpClient {
	if pet.ProxyUrl == "" {
		return myClient
	}
	t := &HttpClient{}
	t.Transport = getTransport(pet.ProxyUrl)
	return t
}

func (cli *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return cli.Client.Do(req)
}

func (cli *HttpClient) DoGetKV(req *http.Request) (cst.KV, error) {
	return parseJsonResponse(cli.Client.Do(req))
}
