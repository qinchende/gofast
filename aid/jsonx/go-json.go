//go:build gojson

package jsonx

// Note(add by chende 20230424):
// go-json 性能不错，可以考虑使用
import (
	goj "github.com/goccy/go-json"
	"github.com/qinchende/gofast/core/lang"
	"io"
)

var (
	Marshal       = goj.Marshal
	MarshalIndent = goj.MarshalIndent
	NewDecoder    = goj.NewDecoder
	NewEncoder    = goj.NewEncoder
)

func Unmarshal(v any, data []byte) error {
	return goj.Unmarshal(data, v)
}

func UnmarshalFromString(v any, str string) error {
	return Unmarshal(v, lang.STB(str))
}

func UnmarshalFromReader(v any, reader io.Reader) error {
	decoder := NewDecoder(reader)
	return decoder.Decode(v)
}

// ext +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
const maxJsonStrLen = math.MaxInt32 - 1 // 最大解析2GB JSON字符串

func DecodeRequest(v any, req *http.Request) error {
	return decodeFromReader(v, req.Body, req.ContentLength)
}

func DecodeReader(v any, reader io.Reader, bufSize int64) error {
	return decodeFromReader(v, reader, bufSize)
}

func decodeFromReader(dst any, reader io.Reader, ctSize int64) error {
	// 一次性读取完成，或者遇到EOF标记或者其它错误
	if ctSize > maxJsonStrLen {
		ctSize = maxJsonStrLen
	}
	bytes, err1 := iox.ReadAll(reader, ctSize)
	if err1 != nil {
		return err1
	}
	return Unmarshal(dst, bytes)
}
