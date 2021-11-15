// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"net/http"
)

type funcServeHTTP func(w http.ResponseWriter, r *http.Request)

type FitHelper struct {
	nextHandler funcServeHTTP
}

func (fh *FitHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.nextHandler(w, r)
}
