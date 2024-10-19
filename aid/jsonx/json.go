//go:build !gojson && !jde

package jsonx

import (
	"bytes"
	"encoding/json"
	"github.com/qinchende/gofast/aid/iox"
	"io"
	"math"
	"net/http"
	"strings"
)

var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

func Unmarshal(v any, data []byte) error {
	//return json.Unmarshal(data, v)
	// 为了统一设置 decoder.UseNumber() 这里转换成字符串使用
	return UnmarshalFromReader(v, bytes.NewReader(data))
}

func UnmarshalFromString(v any, str string) error {
	// 这里无形中带来了字符串到字节数组的copy开销
	// 但是不这么做无法解决UseNumber()的问题，标准库有缺陷吧？
	return UnmarshalFromReader(v, strings.NewReader(str))
}

func UnmarshalFromReader(v any, reader io.Reader) error {
	decoder := NewDecoder(reader)
	decoder.UseNumber()
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
	bs, err1 := iox.ReadAll(reader, ctSize)
	if err1 != nil {
		return err1
	}
	return Unmarshal(dst, bs)
}
