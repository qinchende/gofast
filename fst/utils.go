// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"log"
	"runtime"
	"strconv"
	"strings"
)

func getRuntimeMainVer(v string) (float64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseFloat(v[first-1:], 64)
	}
	return strconv.ParseFloat(v[first-1:last], 64)
}

func checkRuntimeVer() {
	if v, e := getRuntimeMainVer(runtime.Version()); e == nil && v < gftSupportMinGoVer {
		log.Fatal("[Error] Now GoFast requires Go 1.18 or later and Go 1.20 will be required soon.")
	}
}
