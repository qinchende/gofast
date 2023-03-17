//go:build !jsoniter

package jsonx

import (
	"encoding/json"
	"github.com/qinchende/gofast/skill/lang"
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
	// return json.Unmarshal(data, v)
	// 为了统一设置 decoder.UseNumber() 这里转换成字符串使用
	return UnmarshalFromString(v, lang.BTS(data))
}

func UnmarshalFromString(v any, str string) error {
	decoder := NewDecoder(strings.NewReader(str))
	decoder.UseNumber()
	return decoder.Decode(v)
}

func UnmarshalFromReader(v any, reader io.Reader) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := NewDecoder(teeReader)
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
