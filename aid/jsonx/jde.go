//go:build jde

package jsonx

// Note(add by sdx 20241019):
// 自己实现的JSON编解码库，性能不错，可以考虑使用
import (
	"github.com/qinchende/gofast/core/lang"
	"github.com/qinchende/gofast/store/jde"
	"io"
)

var (
	Marshal       = jde.EncodeToBytes
	MarshalIndent = jde.EncodeToBytesIndent
	NewDecoder    = jde.NewDecoder
	NewEncoder    = jde.NewEncoder

	DecodeRequest = jde.DecodeRequest
	DecodeReader  = jde.DecodeReader
)

func Unmarshal(v any, data []byte) error {
	return jde.DecodeBytes(data, v)
}

func UnmarshalFromString(v any, str string) error {
	return Unmarshal(v, lang.S2B(str))
}

func UnmarshalFromReader(v any, reader io.Reader) error {
	decoder := NewDecoder(reader)
	return decoder.Decode(v)
}
