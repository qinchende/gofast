//go:build jsoniter

package jsonx

import (
	jsonIterator "github.com/json-iterator/go"
	"io"
	"strings"
)

var (
	jit = jsonIterator.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}.Froze()

	Marshal       = jit.Marshal
	MarshalIndent = jit.MarshalIndent
	NewDecoder    = jit.NewDecoder
	NewEncoder    = jit.NewEncoder
)

func Unmarshal(v any, data []byte) error {
	return jit.Unmarshal(data, v)
}

func UnmarshalFromString(v any, str string) error {
	decoder := NewDecoder(strings.NewReader(str))
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

func UnmarshalFromReader(v any, reader io.Reader) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := NewDecoder(teeReader)
	return decoder.Decode(v)
}
