// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package httpx

import (
	"github.com/qinchende/gofast/core/cst"
	"net/url"
	"strings"
)

// 解析url中的query参数
// 发生error，意味着解析不到合适的值，啥都不做了。
func ParseQuery(pms cst.SuperKV, query string) {
	for query != "" {
		var key, value string

		// k1=v1&k1=v2
		idx := strings.IndexByte(query, '&')
		if idx < 0 {
			key = query
			query = ""
		} else {
			key, query = query[:idx], query[idx+1:]
		}

		if key == "" {
			continue
		}
		if strings.IndexByte(key, ';') >= 0 {
			continue
		}

		// k1=v1
		idx = strings.IndexByte(key, '=')
		if idx <= 0 {
			continue
		}
		key, value = key[:idx], key[idx+1:]

		// check
		var err error
		if key, err = url.QueryUnescape(key); err != nil {
			continue
		}
		if value, err = url.QueryUnescape(value); err != nil {
			continue
		}

		pms.SetString(key, value)
	}
}
