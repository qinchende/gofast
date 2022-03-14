// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bind

import (
	"bytes"
	"github.com/qinchende/gofast/skill/json"
	"io"
	"net/http"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	//if req == nil || req.Body == nil {
	//	return fmt.Errorf("invalid request")
	//}
	return decodeJSON(req.Body, obj)
}

//func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
//	bodyStr := []byte("{}")
//	_, _ = req.Body.Read(bodyStr)
//	return decodeJSON(bytes.NewReader(bodyStr), obj)
//}

func (jsonBinding) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	// Modify by sdx on 20220311 如果body为空，不抛异常
	// if err := decoder.Decode(obj); err != nil {
	if err := decoder.Decode(obj); err != nil && err != io.EOF {
		return err
	}
	return validate(obj)
}
