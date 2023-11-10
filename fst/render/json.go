// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"bytes"
	"fmt"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/jde"
	"html/template"
	"net/http"
)

var jsonContentType = []string{"application/json; charset=utf-8"}
var jsonpContentType = []string{"application/javascript; charset=utf-8"}
var jsonAsciiContentType = []string{"application/json"}

// JSON contains the given interface object.
type JSON struct {
	Data any
}

// Render (JSON) writes data with custom ContentType.
func (r JSON) Write(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	setContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj any) error {
	setContentType(w, jsonContentType)
	jsonBytes, err := jde.EncodeToBytes(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// other kind of json

// IndentedJSON contains the given interface object.
type IndentedJSON struct {
	Data any
}

// SecureJSON contains the given interface object and its prefix.
type SecureJSON struct {
	Prefix string
	Data   any
}

// JsonpJSON contains the given interface object its callback.
type JsonpJSON struct {
	Callback string
	Data     any
}

// AsciiJSON contains the given interface object.
type AsciiJSON struct {
	Data any
}

// SecureJSONPrefix is a string which represents SecureJSON prefix.
type SecureJSONPrefix string

// PureJSON contains the given interface object.
type PureJSON struct {
	Data any
}

// Render (IndentedJSON) marshals the given interface object and writes it with custom ContentType.
func (r IndentedJSON) Write(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := jde.EncodeToBytesIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

// WriteContentType (IndentedJSON) writes JSON ContentType.
func (r IndentedJSON) WriteContentType(w http.ResponseWriter) {
	setContentType(w, jsonContentType)
}

// Render (SecureJSON) marshals the given interface object and writes it with custom ContentType.
func (r SecureJSON) Write(w http.ResponseWriter) error {
	r.WriteContentType(w)
	jsonBytes, err := jde.EncodeToBytes(r.Data)
	if err != nil {
		return err
	}
	// if the jsonBytes is array values
	if bytes.HasPrefix(jsonBytes, lang.STB("[")) && bytes.HasSuffix(jsonBytes,
		lang.STB("]")) {
		_, err = w.Write(lang.STB(r.Prefix))
		if err != nil {
			return err
		}
	}
	_, err = w.Write(jsonBytes)
	return err
}

// WriteContentType (SecureJSON) writes JSON ContentType.
func (r SecureJSON) WriteContentType(w http.ResponseWriter) {
	setContentType(w, jsonContentType)
}

// Render (JsonpJSON) marshals the given interface object and writes it and its callback with custom ContentType.
func (r JsonpJSON) Write(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := jde.EncodeToBytes(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		_, err = w.Write(ret)
		return err
	}

	callback := template.JSEscapeString(r.Callback)
	_, err = w.Write(lang.STB(callback))
	if err != nil {
		return err
	}
	_, err = w.Write(lang.STB("("))
	if err != nil {
		return err
	}
	_, err = w.Write(ret)
	if err != nil {
		return err
	}
	_, err = w.Write(lang.STB(");"))
	if err != nil {
		return err
	}

	return nil
}

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (r JsonpJSON) WriteContentType(w http.ResponseWriter) {
	setContentType(w, jsonpContentType)
}

// Render (AsciiJSON) marshals the given interface object and writes it with custom ContentType.
func (r AsciiJSON) Write(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := jde.EncodeToBytes(r.Data)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	for _, r := range lang.BTS(ret) {
		cvt := string(r)
		if r >= 128 {
			cvt = fmt.Sprintf("\\u%04x", int64(r))
		}
		buffer.WriteString(cvt)
	}

	_, err = w.Write(buffer.Bytes())
	return err
}

// WriteContentType (AsciiJSON) writes JSON ContentType.
func (r AsciiJSON) WriteContentType(w http.ResponseWriter) {
	setContentType(w, jsonAsciiContentType)
}

//// Render (PureJSON) writes custom ContentType and encodes the given interface object.
//func (r PureJSON) Write(w http.ResponseWriter) error {
//	r.WriteContentType(w)
//	encoder := jsonx.NewEncoder(w)
//	encoder.SetEscapeHTML(false)
//	return encoder.Encode(r.Data)
//}

// WriteContentType (PureJSON) writes custom ContentType.
func (r PureJSON) WriteContentType(w http.ResponseWriter) {
	setContentType(w, jsonContentType)
}
