// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package httpx

import (
	"errors"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

// Note：这里的解析函数主要来自标准库，为配合框架需要，做了适当的改动
// add by sdx on 20230327
var (
	ErrMissingBoundary = errors.New("no multipart boundary param in Content-Type")
	ErrNotMultipart    = errors.New("request Content-Type isn't multipart/form-data")
)

func ParseMultipartForm(pms cst.SuperKV, r *http.Request, maxMemory int64) error {
	// 看是否已经解析过，同时如果有上传文件，文件的信息将被解析在r.MultipartForm中
	if r.MultipartForm != nil {
		return errors.New("http: multipart already parsed")
	}
	ct := r.Header.Get("Content-Type")

	parseFormErr := parseForm(pms, r, ct)

	// 流式数据解析 ++++++++++++++++++++++
	mr, err := multipartReader(r, ct, false)
	if err != nil {
		return err
	}
	f, err := mr.ReadForm(maxMemory)
	if err != nil {
		return err
	}
	// 流式数据解析结果就在MultipartForm
	r.MultipartForm = f
	// END +++++++++++++++++++++++++++++++

	for k := range f.Value {
		pms.Set(k, f.Value[k][0]) // 提取流式解析得到的键值对
	}
	return parseFormErr
}

// +++++++++++++++++++++++++++++++++++++
// 我们不需要解析结果放入r.Form或者r.PostForm对象中，直接放入cst.SuperKV
func parseForm(pms cst.SuperKV, r *http.Request, ct string) error {
	var err error
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		err = parsePostForm(pms, r, ct)
	}
	ParseQuery(pms, r.URL.RawQuery) // url中参数优先级高于post提交，相同字段则覆盖
	return err
}

const (
	//maxPostFormSize    = int64(1<<63 - 1)
	maxLimitReaderSize = int64(10<<20) + 1 // 10 MB is a lot of text.
)

// 解决：application/x-www-form-urlencoded
func parsePostForm(pms cst.SuperKV, r *http.Request, ct string) error {
	if r.Body == nil {
		return errors.New("missing form body")
	}

	if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		// TODO：这种超大Body情况先不考虑
		//var reader io.Reader = r.Body
		//maxFormSize := maxPostFormSize
		//if _, ok := r.Body.(*maxBytesReader); !ok {
		//	maxFormSize = maxLimitReaderSize
		//	reader = io.LimitReader(r.Body, maxFormSize)
		//}

		// content-type: application/x-www-form-urlencoded
		// 最多读取10MB（10x1024x1024B）的内容，多了就丢弃了
		// 几乎http协议中指定了Content-Length这个header的请求，其Body都是LimitReader了
		// 参考源代码:  net/http/transfer.go 565行 case realLength > 0:
		reader := io.LimitReader(r.Body, maxLimitReaderSize)
		// 一次性读取完成，知道读取maxLimitReaderSize字节，或者遇到EOF标记
		// add by sdx on 20230329
		bytes, err := ReadAll(reader, r.ContentLength)
		if err != nil {
			return err
		}
		if int64(len(bytes)) > maxLimitReaderSize {
			return errors.New("http: POST too large")
		}

		ParseQuery(pms, lang.BTS(bytes))
	}
	return nil
}

// 流式数据体的处理，比如上传文件
// 解决："multipart/form-data"
func multipartReader(r *http.Request, ct string, allowMixed bool) (*multipart.Reader, error) {
	if ct == "" {
		return nil, ErrNotMultipart
	}
	// TODO: 这里可能有大小写的问题
	if allowMixed == false && !strings.HasPrefix(ct, "multipart/form-data") {
		return nil, ErrNotMultipart
	}

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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//type maxBytesReader struct {
//	w   http.ResponseWriter
//	r   io.ReadCloser // underlying reader
//	n   int64         // max bytes remaining
//	err error         // sticky error
//}
//
//func (l *maxBytesReader) Read(p []byte) (n int, err error) {
//	return 0, nil
//}
//
//func (l *maxBytesReader) Close() error {
//	return nil
//}

//func copyValues(dst, src url.Values) {
//	for k, vs := range src {
//		dst[k] = append(dst[k], vs...)
//	}
//}

// Copy from io/io.go 638行的函数，用最有可能的[]byte长度申请内存空间，防止动态扩容
// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r io.Reader, size int64) ([]byte, error) {
	// 内存空间尽量一次性分配到位
	b := make([]byte, 0, size)
	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)] // Add more capacity (let append pick how much).
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}
