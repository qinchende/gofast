//go:build jsoniter

package jsonx

// Note(add by chende 20230315):
// 我简单测试过，在绝大多数情况下 json-iterator 并不比标准库提升多少性能。替换标准库的意义不大。
// 真要大幅提升Decode性能，要找其它方案。
import (
	jsonIterator "github.com/json-iterator/go"
	"github.com/qinchende/gofast/skill/lang"
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
	return Unmarshal(v, lang.StringToBytes(str))
}

func UnmarshalFromReader(v any, reader io.Reader) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := NewDecoder(teeReader)
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
