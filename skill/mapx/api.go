package mapx

import "github.com/qinchende/gofast/cst"

// 不设置Default值
func BindKV(dst interface{}, kvs cst.KV) error {
	return mapKVJust(dst, kvs)
}

func BindKVDef(dst interface{}, kvs cst.KV) error {
	return mapKVApplyDefault(dst, kvs)
}

func BindKVValid(dst interface{}, kvs cst.KV) error {
	if err := mapKVApplyDefault(dst, kvs); err != nil {
		return err
	}
	return Validate(dst)
}
