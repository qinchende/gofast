// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"compress/gzip"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"net/http"
	"strings"
)

func Gunzip(useGunzip bool) fst.CtxHandler {
	if useGunzip == false {
		return nil
	}

	return func(c *fst.Context) {
		if strings.Contains(c.GetHeader(cst.HeaderContentEncoding), "gzip") {
			reader, err := gzip.NewReader(c.Req.Raw.Body)
			if err != nil {
				c.AbortDirect(http.StatusBadRequest, "Can't unzip body!")
				return
			}
			c.Req.Raw.Body = reader
		}
		c.Next()
	}
}
