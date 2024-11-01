package jde

import (
	"github.com/qinchende/gofast/core/lang"
	"io"
	"net/http"
)

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// bind and valid with json data
//// +++ JSON Bytes
//func BindJsonBytes(dst any, content []byte, like int8) error {
//	return BindJsonBytesX(dst, content, dts.AsOptions(like))
//}
//
//func BindJsonBytesX(dst any, content []byte, opts *dts.BindOptions) error {
//	var kv map[string]any
//	if err := DecodeBytes(&kv, content); err != nil {
//		return err
//	}
//	return dts.BindKVX(dst, kv, opts)
//}
//
//// +++ JSON Reader
//func BindJsonReader(dst any, reader io.Reader, like int8) error {
//	return BindJsonReaderX(dst, reader, dts.AsOptions(like))
//}
//
//func BindJsonReaderX(dst any, reader io.Reader, opts *dts.BindOptions) error {
//	var kv map[string]any
//	if err := DecodeReader(&kv, reader, 0); err != nil {
//		return err
//	}
//	return dts.BindKVX(dst, kv, opts)
//}

// 解码到对象
// Important: 被解析的数据源 source 必须是只读的，不可在解析后再改写，否则可能造成意想不到的错误
// 如果想要避免这样的问题，请将copy(source)后的source传入
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DecodeString(v any, source string) error {
	return decodeFromString(v, source)
}

func DecodeBytes(v any, source []byte) error {
	return decodeFromString(v, lang.B2S(source))
}

// 事先精确指定 bufSize ，能有效避免字节数组中途扩容，若未知，传 0 即可，默认初始化内存空间
func DecodeReader(v any, reader io.Reader, bufSize int64) error {
	return decodeFromReader(v, reader, bufSize)
}

func DecodeRequest(v any, req *http.Request) error {
	return decodeFromReader(v, req.Body, req.ContentLength)
}

// +++++++ Copy source content for safe decode
func DecodeStringCopy(v any, source string) error {
	newMem := make([]byte, len(source))
	copy(newMem, source)
	return decodeFromString(v, lang.B2S(newMem))
}

func DecodeBytesCopy(v any, source []byte) error {
	return decodeFromString(v, string(source))
}

//func Unmarshal(str []byte, v any) error {
//	return decodeFromString(v, lang.B2S(str))
//}
//
//func UnmarshalFromString(str string, v any) error {
//	return decodeFromString(v, str)
//}

//	编码成JSON字符串
//
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func EncodeToBytes(v any) ([]byte, error) {
	return startEncode(v)
}

func EncodeToString(v any) (string, error) {
	b, err := startEncode(v)
	return lang.B2S(b), err
}

func EncodeToBytesIndent(v any, prefix, indent string) ([]byte, error) {
	return nil, nil
}

//func Marshal(v any) ([]byte, error) {
//	return startEncode(v)
//}
//
//func MarshalToString(v any) (string, error) {
//	b, err := startEncode(v)
//	return lang.B2S(b), err
//}
//
//func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
//	return nil, nil
//}

// extend 编解码读写流
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type decoder struct {
	dec *subDecode
}

type encoder struct {
	enc *subEncode
}

// Encoder encodes JSON into io.Writer
type Encoder interface {
	// Encode writes the JSON encoding of v to the stream, followed by a newline character.
	Encode(val interface{}) error
	// SetEscapeHTML specifies whether problematic HTML characters
	// should be escaped inside JSON quoted strings.
	// The default behavior NOT ESCAPE
	SetEscapeHTML(on bool)
	// SetIndent instructs the encoder to format each subsequent encoded value
	// as if indented by the package-level function Indent(v, src, prefix, indent).
	// Calling SetIndent("", "") disables indentation
	SetIndent(prefix, indent string)
}

// Decoder decodes JSON from io.Read
type Decoder interface {
	// Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
	Decode(val interface{}) error
	// Buffered returns a reader of the data remaining in the Decoder's buffer.
	// The reader is valid until the next call to Decode.
	Buffered() io.Reader
	// DisallowUnknownFields causes the Decoder to return an error when the destination is a struct
	// and the input contains object keys which do not match any non-ignored, exported fields in the destination.
	DisallowUnknownFields()
	// More reports whether there is another element in the current array or object being parsed.
	More() bool
	// UseNumber causes the Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	UseNumber()
}

// NewEncoder create a Encoder holding writer
func NewEncoder(writer io.Writer) Encoder {
	return nil
}

// NewDecoder create a Decoder holding reader
func NewDecoder(reader io.Reader) Decoder {
	return nil
}

func Valid(data []byte) bool {
	return false
}

// Config params
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Config is a combination of sonic/encoder.Options and sonic/decoder.Options
type Config struct {
	// EscapeHTML indicates encoder to escape all HTML characters
	// after serializing into JSON (see https://pkg.go.dev/encoding/json#HTMLEscape).
	// WARNING: This hurts performance A LOT, USE WITH CARE.
	EscapeHTML bool

	// SortMapKeys indicates encoder that the keys of a map needs to be sorted
	// before serializing into JSON.
	// WARNING: This hurts performance A LOT, USE WITH CARE.
	SortMapKeys bool

	// CompactMarshaler indicates encoder that the output JSON from json.Marshaler
	// is always compact and needs no validation
	CompactMarshaler bool

	// NoQuoteTextMarshaler indicates encoder that the output text from encoding.TextMarshaler
	// is always escaped string and needs no quoting
	NoQuoteTextMarshaler bool

	// NoNullSliceOrMap indicates encoder that all empty Array or Object are encoded as '[]' or '{}',
	// instead of 'null'
	NoNullSliceOrMap bool

	// UseInt64 indicates decoder to unmarshal an integer into an interface{} as an
	// int64 instead of as a float64.
	UseInt64 bool

	// UseNumber indicates decoder to unmarshal a number into an interface{} as a
	// json.Number instead of as a float64.
	UseNumber bool

	// UseUnicodeErrors indicates decoder to return an error when encounter invalid
	// UTF-8 escape sequences.
	UseUnicodeErrors bool

	// DisallowUnknownFields indicates decoder to return an error when the destination
	// is a struct and the input contains object keys which do not match any
	// non-ignored, exported fields in the destination.
	DisallowUnknownFields bool

	// CopyString indicates decoder to decode string values by copying instead of referring.
	CopyString bool

	// ValidateString indicates decoder and encoder to valid string values: decoder will return errors
	// when unescaped control chars(\u0000-\u001f) in the string value of JSON.
	ValidateString bool
}
