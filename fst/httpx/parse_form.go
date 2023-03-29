// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package httpx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	maxPostContentReadSize = int64(10 << 20) // 10 MB is a lot of text.
)

// Note：这里的解析函数主要来自标准库，为配合框架需要，做了适当的改动
// add by sdx on 20230327
var (
	ErrMissingBoundary = errors.New("no multipart boundary param in Content-Type")
	ErrNotMultipart    = errors.New("request Content-Type isn't multipart/form-data")
)

// 用来占位，不会有任何数据。Important! -> 客户端不要设置r.MultipartForm的值。
// multipartByReader is a sentinel value.
// Its presence in Request.MultipartForm indicates that parsing of the request
// body has been handed off to a MultipartReader instead of ParseMultipartForm.
var multipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

// 解析Http提交的数据，一般来说此方法最多只执行一次
func ParseMultipartForm(pms cst.SuperKV, r *http.Request, ct string, maxMemory int64) error {
	// 看是否已经解析过，同时如果有上传文件，文件的信息将被解析在r.MultipartForm中
	if r.MultipartForm != nil {
		return errors.New("http: multipart form already parsed")
	}
	if r.Body == nil {
		return errors.New("missing form body")
	}

	var parseFormErr error
	// 1. 解析Body：application/x-www-form-urlencoded
	if (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") &&
		strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		// TODO：可以自定义maxBytesReader，防止客户在当前content-type模式下提交大量数据攻击服务器处理能力
		// 但是下面用LimitReader限制了这种缓冲区读取数据量，起到一定保护作用。
		//var reader io.Reader = r.Body
		//maxFormSize := maxPostFormSize
		//if _, ok := r.Body.(*maxBytesReader); !ok {
		//	maxFormSize = maxLimitReaderSize
		//	reader = io.LimitReader(r.Body, maxFormSize)
		//}

		limitSize := maxPostContentReadSize
		if limitSize > maxMemory && maxMemory >= 0 {
			limitSize = maxMemory
		}
		// content-type: application/x-www-form-urlencoded
		// 最多读取10MB（10x1024x1024B）的内容，多了就丢弃了
		// 几乎http协议中指定了Content-Length这个header的请求，其Body都是LimitReader了
		// 参考源代码:  net/http/transfer.go 565行 case realLength > 0:
		reader := io.LimitReader(r.Body, limitSize)
		// 一次性读取完成，知道读取maxLimitReaderSize字节，或者遇到EOF标记
		// add by sdx on 20230329
		bytes, err := iox.ReadAll(reader, r.ContentLength)
		if err != nil {
			parseFormErr = err
		} else {
			if int64(len(bytes)) > limitSize {
				parseFormErr = errors.New("http: POST too large")
			} else {
				ParseQuery(pms, lang.BTS(bytes))
			}
		}
	}
	// 2. 解析URL：query params like k1=b1&k2=v2
	ParseQuery(pms, r.URL.RawQuery) // url中参数优先级高于post提交，相同字段则覆盖
	if parseFormErr != nil {
		r.MultipartForm = multipartByReader
		return parseFormErr
	}

	// 3. 解析Body：multipart/form-data，比如上传文件这种场景
	// TODO: 这里可能有大小写的问题
	if strings.HasPrefix(ct, "multipart/form-data") {
		mr, err := multipartReader(r, ct, false)
		if err != nil {
			return err
		}
		f, err := mr.ReadForm(maxMemory)
		if err != nil {
			return err
		}
		r.MultipartForm = f // 流式数据解析结果就在MultipartForm

		for k := range f.Value {
			pms.Set(k, f.Value[k][0]) // 提取流式解析得到的键值对
		}
	}
	if r.MultipartForm == nil {
		r.MultipartForm = multipartByReader
	}
	return nil
}

// 流式数据体的处理，比如上传文件
// 解决："multipart/form-data"
func multipartReader(r *http.Request, ct string, allowMixed bool) (*multipart.Reader, error) {
	d, params, err := mime.ParseMediaType(ct)
	if err != nil || !(d == "multipart/form-data" || (allowMixed && d == "multipart/mixed")) {
		return nil, ErrNotMultipart
	}
	boundary, ok := params["boundary"]
	if !ok {
		return nil, ErrMissingBoundary
	}
	return multipart.NewReader(r.Body, boundary), nil
}
