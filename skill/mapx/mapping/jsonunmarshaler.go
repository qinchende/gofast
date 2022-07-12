package mapping

import (
	"io"

	"github.com/qinchende/gofast/skill/jsonx"
)

const jsonTagKey = "cnf"

var jsonUnmarshaler = NewUnmarshaler(jsonTagKey)

func UnmarshalJsonBytes(content []byte, v any) error {
	return unmarshalJsonBytes(content, v, jsonUnmarshaler)
}

func UnmarshalJsonReader(reader io.Reader, v any) error {
	return unmarshalJsonReader(reader, v, jsonUnmarshaler)
}

func unmarshalJsonBytes(content []byte, v any, unmarshaler *Unmarshaler) error {
	var m map[string]any
	if err := jsonx.Unmarshal(&m, content); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}

func unmarshalJsonReader(reader io.Reader, v any, unmarshaler *Unmarshaler) error {
	var m map[string]any
	if err := jsonx.UnmarshalFromReader(&m, reader); err != nil {
		return err
	}

	return unmarshaler.Unmarshal(m, v)
}
