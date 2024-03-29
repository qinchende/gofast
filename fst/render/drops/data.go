// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package drops

import (
	"net/http"
)

// Data contains ContentType and bytes data.
type Data struct {
	ContentType string
	Data        []byte
}

// Render (Data) writes data with custom ContentType.
func (r Data) Write(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	_, err = w.Write(r.Data)
	return
}

// WriteContentType (Data) writes custom ContentType.
func (r Data) WriteContentType(w http.ResponseWriter) {
	//render.writeContentType(w, []string{r.ContentType})
}
