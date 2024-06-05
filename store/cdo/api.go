// Copyright 2024 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package cdo

import (
	"github.com/qinchende/gofast/aid/lang"
	"io"
	"net/http"
)

// Warning!!!
// This package just for small terminal machine.

// cdo (Compact data of object) (紧凑数据对象)
//
//	编码成Cdo字符串
//
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func EncodeToBytes(v any) ([]byte, error) {
	return cdoEncode(v)
}

func EncodeToString(v any) (string, error) {
	b, err := cdoEncode(v)
	return lang.BTS(b), err
}

func EncodeToBytesIndent(v any, prefix, indent string) ([]byte, error) {
	return nil, nil
}

// 解码 Cdo 数据
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func DecodeString(v any, source string) error {
	return decodeFromString(v, source)
}

func DecodeBytes(v any, source []byte) error {
	return decodeFromString(v, lang.BTS(source))
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
	return decodeFromString(v, lang.BTS(newMem))
}

func DecodeBytesCopy(v any, source []byte) error {
	return decodeFromString(v, string(source))
}
