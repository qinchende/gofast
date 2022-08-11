// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"fmt"
	"io"
	"net/http"
)

// String contains the given interface object slice and its format.
type Text struct {
	Format string
	Data   []any
}

var plainContentType = []string{"text/plain; charset=utf-8"}

// Render (String) writes data with custom ContentType.
func (r Text) Write(w http.ResponseWriter) error {
	return WriteString(w, r.Format, r.Data)
}

// WriteContentType (String) writes Plain ContentType.
func (r Text) WriteContentType(w http.ResponseWriter) {
	setContentType(w, plainContentType)
}

// WriteString writes data according to its format and write custom ContentType.
func WriteString(w http.ResponseWriter, format string, data []any) (err error) {
	setContentType(w, plainContentType)
	if len(data) > 0 {
		_, err = fmt.Fprintf(w, format, data...)
		return
	}
	_, err = io.WriteString(w, format)
	return
}
