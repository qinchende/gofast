package mapx

import "github.com/qinchende/gofast/cst"

// 不设置Default值
func BindKV(dest interface{}, kvs cst.KV) error {
	return mapPms(dest, kvs)
}

func BindKVDef(dest interface{}, kvs cst.KV) error {
	return mapPms(dest, kvs)
}

func BindKVValid(dest interface{}, kvs cst.KV) error {
	if err := mapPms(dest, kvs); err != nil {
		return err
	}
	return Validate(dest)
}
