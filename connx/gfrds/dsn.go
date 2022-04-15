package gfrds

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/bind"
	"strings"
)

func ParseDsn(connStr string) *ConnCnf {
	rdsCnf := &ConnCnf{}

	kvs := fst.KV{}
	items := strings.Split(connStr, "&")
	for _, item := range items {
		idx := strings.IndexByte(item, '=')
		kvs[item[:idx]] = item[idx+1:]
	}

	_ = bind.Pms.BindPms(kvs, rdsCnf)
	return rdsCnf
}
