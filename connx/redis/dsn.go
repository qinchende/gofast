package redis

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/bind"
	"strings"
)

func ParseDsn(connStr string) *ConnCnf {
	rdsCnf := &ConnCnf{}

	kvs := cst.KV{}
	items := strings.Split(connStr, "&")
	for _, item := range items {
		idx := strings.IndexByte(item, '=')
		kvs[item[:idx]] = item[idx+1:]
	}

	_ = bind.BindKV(rdsCnf, kvs, bind.AsConfig)
	return rdsCnf
}
