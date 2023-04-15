package gson

//import (
//	"github.com/qinchende/gofast/skill/lang"
//	"net/http"
//)
//
//// ++++++++++++++++ 解码Gson文本
//func DecodeGsonString(dst any, source string) error {
//	return decodeFromString(dst, source)
//}
//
//func DecodeGsonBytes(dst any, source []byte) error {
//	return decodeFromString(dst, lang.BTS(source))
//}
//
//func DecodeGsonReader(dst any, reader io.Reader, ctSize int64) error {
//	return decodeFromReader(dst, reader, ctSize)
//}
//
//func DecodeGsonRequest(dst any, req *http.Request) error {
//	return decodeFromReader(dst, req.Body, req.ContentLength)
//}
