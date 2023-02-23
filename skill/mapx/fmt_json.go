// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"io"
)

// +++ JSON Bytes
func DecodeJsonBytes(dst any, content []byte, like int8) error {
	return DecodeJsonBytesX(dst, content, matchOptions(like))
}

func DecodeJsonBytesX(dst any, content []byte, opts *BindOptions) error {
	var kv map[string]any
	if err := jsonx.Unmarshal(&kv, content); err != nil {
		return err
	}
	return BindKVX(dst, kv, opts)
}

// +++ JSON Reader
func DecodeJsonReader(dst any, reader io.Reader, like int8) error {
	return DecodeJsonReaderX(dst, reader, matchOptions(like))
}

func DecodeJsonReaderX(dst any, reader io.Reader, opts *BindOptions) error {
	var kv map[string]any
	if err := jsonx.UnmarshalFromReader(&kv, reader); err != nil {
		return err
	}
	return BindKVX(dst, kv, opts)
}
