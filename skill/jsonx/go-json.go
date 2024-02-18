//go:build gojson

package jsonx

// Note(add by chende 20230424):
// go-json 性能不错，可以考虑使用
import (
	goj "github.com/goccy/go-json"
	"github.com/qinchende/gofast/skill/lang"
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
