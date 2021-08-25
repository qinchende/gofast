// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrorType is an unsigned 64-bit error code as defined in the gin spec.
type ErrorType uint64

const (
	// ErrorTypeBind is used when Context.Bind() fails.
	ErrorTypeBind ErrorType = 1 << 63
	// ErrorTypeRender is used when Context.Render() fails.
	ErrorTypeRender ErrorType = 1 << 62
	// ErrorTypePrivate indicates a private error.
	ErrorTypePrivate ErrorType = 1 << 0
	// ErrorTypePublic indicates a public error.
	ErrorTypePublic ErrorType = 1 << 1
	// ErrorTypeAny indicates any other error.
	ErrorTypeAny ErrorType = 1<<64 - 1
	// ErrorTypeNu indicates any other error.
	ErrorTypeNu = 2
)

// Error represents a error's specification.
type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

type errMessages []*Error

var _ error = &Error{}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
/************************************/
/********* error management *********/
/************************************/

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) CollectError(err error) *Error {
	if err == nil {
		RaisePanic("err is nil")
	}
	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}
	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

func (w *GFResponse) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	w.Errors = append(w.Errors, parsedError)
	return parsedError
}

func (w *GFResponse) ErrorN(err interface{}) {
	//_ = w.Error(err)
}

func (w *GFResponse) ErrorF(format string, v ...interface{}) {
	_ = w.Error(errors.New(fmt.Sprintf(format, v...)))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// SetType sets the error's type.
func (msg *Error) SetType(flags ErrorType) *Error {
	msg.Type = flags
	return msg
}

// SetMeta sets the error's meta data.
func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

// JSON creates a properly formatted JSON
func (msg *Error) JSON() interface{} {
	hash := KV{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				hash[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			hash["meta"] = msg.Meta
		}
	}
	if _, ok := hash["error"]; !ok {
		hash["error"] = msg.Error()
	}
	return hash
}

// MarshalJSON implements the json.Marshaller interface.
func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

// Error implements the error interface.
func (msg Error) Error() string {
	return msg.Err.Error()
}

// IsType judges one error.
func (msg *Error) IsType(flags ErrorType) bool {
	return (msg.Type & flags) > 0
}

// ByType returns a readonly copy filtered the byte.
// ie ByType(gin.ErrorTypePublic) returns a slice of errors with type=ErrorTypePublic.
func (a errMessages) ByType(typ ErrorType) errMessages {
	if len(a) == 0 {
		return nil
	}
	if typ == ErrorTypeAny {
		return a
	}
	var result errMessages
	for _, msg := range a {
		if msg.IsType(typ) {
			result = append(result, msg)
		}
	}
	return result
}

// Last returns the last error in the slice. It returns nil if the array is empty.
// Shortcut for errors[len(errors)-1].
func (a errMessages) Last() *Error {
	if length := len(a); length > 0 {
		return a[length-1]
	}
	return nil
}

// Errors returns an array will all the error messages.
// Example:
// 		c.Error(errors.New("first"))
// 		c.Error(errors.New("second"))
// 		c.Error(errors.New("third"))
// 		c.Errors.Errors() // == []string{"first", "second", "third"}
func (a errMessages) Errors() []string {
	if len(a) == 0 {
		return nil
	}
	errorStrings := make([]string, len(a))
	for i, err := range a {
		errorStrings[i] = err.Error()
	}
	return errorStrings
}

func (a errMessages) JSON() interface{} {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a.Last().JSON()
	default:
		json := make([]interface{}, len(a))
		for i, err := range a {
			json[i] = err.JSON()
		}
		return json
	}
}

// MarshalJSON implements the json.Marshaller interface.
func (a errMessages) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.JSON())
}

func (a errMessages) String() string {
	if len(a) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range a {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
