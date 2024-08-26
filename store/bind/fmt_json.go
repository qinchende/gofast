// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package bind

import (
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/dts"
	"github.com/qinchende/gofast/store/jde"
	"io"
)

// +++ JSON Bytes
func BindJsonBytes(dst any, content []byte, like int8) error {
	return BindJsonBytesX(dst, content, dts.AsOptions(like))
}

func BindJsonBytesX(dst any, content []byte, opts *dts.BindOptions) error {
	var kv cst.KV
	if err := jde.DecodeBytes(&kv, content); err != nil {
		return err
	}
	return BindKVX(dst, kv, opts)
}

// +++ JSON Reader
func BindJsonReader(dst any, reader io.Reader, like int8) error {
	return BindJsonReaderX(dst, reader, dts.AsOptions(like))
}

func BindJsonReaderX(dst any, reader io.Reader, opts *dts.BindOptions) error {
	var kv cst.KV
	if err := jde.DecodeReader(&kv, reader, 0); err != nil {
		return err
	}
	return BindKVX(dst, kv, opts)
}
