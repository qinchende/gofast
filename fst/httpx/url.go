package httpx

import (
	"github.com/qinchende/gofast/cst"
	"net/url"
	"strings"
)

// 解析url中的query参数
func ParseQuery(pms cst.SuperKV, query string) {
	for query != "" {
		var key string

		// k1=v1&k1=v2
		idx := strings.IndexByte(query, '&')
		if idx < 0 {
			key = query
			query = ""
		} else {
			key, query = query[:idx], query[idx+1:]
		}

		//key, query, _ = strings.Cut(query, "&")
		if key == "" {
			continue
		}
		if strings.IndexByte(key, ';') >= 0 {
			continue
		}

		// k = v
		idx = strings.IndexByte(key, '=')
		if idx < 0 {
			continue
		}
		key, value := key[:idx], key[idx+1:]
		if len(key) == 0 || len(value) == 0 {
			continue
		}

		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			continue
		}
		value, err2 := url.QueryUnescape(value)
		if err2 != nil {
			continue
		}

		pms.Set(key, value)
	}
}
