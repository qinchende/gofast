// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build !nomsgpack

package test

import (
	"bytes"
	"github.com/qinchende/gofast/fst/binding"
	"github.com/qinchende/gofast/fst/cst"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

func TestBindingMsgPack(t *testing.T) {
	test := FooStruct{
		Foo: "bar",
	}

	h := new(codec.MsgpackHandle)
	assert.NotNil(t, h)
	buf := bytes.NewBuffer([]byte{})
	assert.NotNil(t, buf)
	err := codec.NewEncoder(buf, h).Encode(test)
	assert.NoError(t, err)

	data := buf.Bytes()

	testMsgPackBodyBinding(t,
		binding.MsgPack, "msgpack",
		"/", "/",
		string(data), string(data[1:]))
}

func testMsgPackBodyBinding(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := FooStruct{}
	req := requestWithBody("POST", path, body)
	req.Header.Add("Content-Type", cst.MIMEMsgPack)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "bar", obj.Foo)

	obj = FooStruct{}
	req = requestWithBody("POST", badPath, badBody)
	req.Header.Add("Content-Type", cst.MIMEMsgPack)
	err = binding.MsgPack.Bind(req, &obj)
	assert.Error(t, err)
}

func TestBindingDefaultMsgPack(t *testing.T) {
	assert.Equal(t, binding.MsgPack, binding.Default("POST", cst.MIMEMsgPack))
	assert.Equal(t, binding.MsgPack, binding.Default("PUT", cst.MIMEXMsgPack))
}
