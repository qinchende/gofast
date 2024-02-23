//go:build jsoniter

package jsonx

// Note(add by chende 20230315):
// 我简单测试过，在绝大多数情况下 json-iterator 并不比标准库提升多少性能。替换标准库的意义不大。
// 真要大幅提升Decode性能，要找其它方案。
import (
	jsonIterator "github.com/json-iterator/go"
	"github.com/qinchende/gofast/aid/lang"
	"io"
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
	return Unmarshal(v, lang.STB(str))
}

func UnmarshalFromReader(v any, reader io.Reader) error {
	decoder := NewDecoder(reader)
	return decoder.Decode(v)
}
