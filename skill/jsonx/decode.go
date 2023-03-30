package jsonx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"net/http"
)

func UnmarshalGsonRequest(dst cst.SuperKV, req *http.Request) error {
	return gsonDecodeEnter(dst, req.Body, req.ContentLength)
}

func UnmarshalGsonReader(dst cst.SuperKV, reader io.Reader, size int64) error {
	return gsonDecodeEnter(dst, reader, 4096)
}

func gsonDecodeEnter(dst cst.SuperKV, reader io.Reader, size int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	bytes, err1 := iox.ReadAll(reader, size)
	if err1 != nil {
		return err1
	}
	str := lang.BTS(bytes)

	//fmt.Println(str)
	decode := gsonDecode{}
	if err2 := decode.init(dst, str); err2 != nil {
		return err2
	}
	return decode.parse()
}
