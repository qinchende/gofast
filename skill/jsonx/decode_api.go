package jsonx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/iox"
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"net/http"
)

// Important: 被解析的数据源 source 必须是只读的，不可在解析后再改写，否则可能造成意想不到的错误
// 如果想要避免这样的问题，请将copy(source)后的source传入

func DecodeString(dst cst.SuperKV, source string) error {
	return decodeFromString(dst, source)
}

func DecodeBytes(dst cst.SuperKV, source []byte) error {
	return decodeFromString(dst, lang.BTS(source))
}

func DecodeReader(dst cst.SuperKV, reader io.Reader, ctSize int64) error {
	return decodeFromReader(dst, reader, ctSize)
}

func DecodeRequest(dst cst.SuperKV, req *http.Request) error {
	return decodeFromReader(dst, req.Body, req.ContentLength)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decodeFromReader(dst cst.SuperKV, reader io.Reader, ctSize int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	bytes, err1 := iox.ReadAll(reader, ctSize)
	if err1 != nil {
		return err1
	}
	str := lang.BTS(bytes)
	return decodeFromString(dst, str)
}

func decodeFromString(dst cst.SuperKV, source string) error {
	decode := fastDecode{}
	if err := decode.init(dst, source); err != nil {
		return err
	}
	return decode.warpError(decode.parseJson())
}
