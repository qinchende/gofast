//go:build !jsoniter && !gojson

package jsonx

import (
	"bytes"
	"encoding/json"
	"io"
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

func UnmarshalStringToKV(str string) (map[string]any, error) {
	res := make(map[string]any)
	if str == "" {
		return res, nil
	}
	err := UnmarshalFromString(&res, str)
	return res, err
}
