package mapx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"io"
)

func DecodeJsonReaderOfData(dst any, reader io.Reader) error {
	return DecodeJsonReaderX(dst, reader, dbStructOptions)
}

func DecodeJsonBytesOfConfig(dst any, content []byte) error {
	return DecodeJsonBytesX(dst, content, configStructOptions)
}

func DecodeJsonReaderX(dst any, reader io.Reader, opts *ApplyOptions) error {
	var kv map[string]any
	if err := jsonx.UnmarshalFromReader(&kv, reader); err != nil {
		return err
	}
	return ApplyKVX(dst, kv, opts)
}

func DecodeJsonBytesX(dst any, content []byte, opts *ApplyOptions) error {
	var kv map[string]any
	if err := jsonx.Unmarshal(&kv, content); err != nil {
		return err
	}
	return ApplyKVX(dst, kv, opts)
}
