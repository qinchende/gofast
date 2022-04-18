package gfrds

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/bind"
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

	_ = bind.Pms.BindPms(rdsCnf, kvs)
	return rdsCnf
}
