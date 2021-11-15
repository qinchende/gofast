// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"compress/gzip"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/httpx"
	"net/http"
	"strings"
)

func Gunzip(ctx *fst.Context) {
	if strings.Contains(ctx.ReqRaw.Header.Get(httpx.ContentEncoding), "gzip") {
		reader, err := gzip.NewReader(ctx.ReqRaw.Body)
		if err != nil {
			ctx.ResWrap.WriteHeader(http.StatusBadRequest)
			ctx.AbortChain()
		}
		ctx.ReqRaw.Body = reader
	}
}
