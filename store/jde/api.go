package jde

import (
	"github.com/qinchende/gofast/skill/lang"
	"io"
	"net/http"
)

// Important: 被解析的数据源 source 必须是只读的，不可在解析后再改写，否则可能造成意想不到的错误
// 如果想要避免这样的问题，请将copy(source)后的source传入

// ++++++++++++++++ 解码到对象
func DecodeString(dst any, source string) error {
	return decodeFromString(dst, source)
}

func DecodeBytes(dst any, source []byte) error {
	return decodeFromString(dst, lang.BTS(source))
}

func DecodeReader(dst any, reader io.Reader, ctSize int64) error {
	return decodeFromReader(dst, reader, ctSize)
}

func DecodeRequest(dst any, req *http.Request) error {
	return decodeFromReader(dst, req.Body, req.ContentLength)
}

//// ++++++++++++++++ 绑定到对象
//func BindString(dst any, source string) error {
//	return decodeFromString(dst, source)
//}
//
//func BindBytes(dst any, source []byte) error {
//	return decodeFromString(dst, lang.BTS(source))
//}
//
//func BindReader(dst any, reader io.Reader, ctSize int64) error {
//	return decodeFromReader(dst, reader, ctSize)
//}
//
//func BindRequest(dst any, req *http.Request) error {
//	return decodeFromReader(dst, req.Body, req.ContentLength)
//}

// ++++++++++++++++ 编码成JSON字符串
func EncodeToBytes(v any) ([]byte, error) {
	return nil, nil
}

func EncodeToString(v any) (string, error) {
	return "", nil
}

func EncodeIndent(v any, prefix, indent string) ([]byte, error) {
	return nil, nil
}