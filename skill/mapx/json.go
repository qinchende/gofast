package mapx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"io"
)

func DecodeJsonReader(dst any, reader io.Reader, opts *ApplyOptions) error {
	var kv map[string]any
	if err := jsonx.UnmarshalFromReader(&kv, reader); err != nil {
		return err
	}

	return ApplyKV(dst, kv, opts)
}

func DecodeJsonBytes(dst any, content []byte, opts *ApplyOptions) error {
	var kv map[string]any
	if err := jsonx.Unmarshal(&kv, content); err != nil {
		return err
	}

	return ApplyKV(dst, kv, opts)
}
