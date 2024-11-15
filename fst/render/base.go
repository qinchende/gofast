// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import (
	"github.com/qinchende/gofast/aid/logx"
	"net/http"
)

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	// Render writes data with custom ContentType.
	Write(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(http.ResponseWriter)
}

var (
	_ Render = Text{}
	_ Render = JSON{}
	_ Render = IndentedJSON{}
	_ Render = SecureJSON{}
	_ Render = JsonpJSON{}
	_ Render = AsciiJSON{}
)

//var (
//	_ Render           = drops.XML{}
//	_ Render           = drops.Redirect{}
//	_ Render           = drops.Data{}
//	_ Render           = drops.HTML{}
//	_ drops.HTMLRender = drops.HTMLDebug{}
//	_ drops.HTMLRender = drops.HTMLProduction{}
//	_ Render           = drops.YAML{}
//	_ Render           = drops.Reader{}
//	_ Render           = drops.ProtoBuf{}
//)

func setContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	// 第一次设置时生效，后面再设置无效
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	} else {
		logx.Info().Msg("The content-type already set.")
	}
}
