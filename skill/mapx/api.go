package mapx

import "github.com/qinchende/gofast/cst"

// 不设置Default值
func ApplyKVByName(dst interface{}, kvs cst.KV) error {
	return applyKV(dst, kvs, true, false)
}

func ApplyKVByNameWithDef(dst interface{}, kvs cst.KV) error {
	return applyKV(dst, kvs, true, true)
}

func ApplyKVByTag(dst interface{}, kvs cst.KV) error {
	return applyKV(dst, kvs, false, false)
}

func ApplyKVByTagWithDef(dst interface{}, kvs cst.KV) error {
	return applyKV(dst, kvs, false, true)
}

//func BindKVDef(dst interface{}, kvs cst.KV) error {
//	return mapKVApplyDefault(dst, kvs)
//}
//
//func BindKVValid(dst interface{}, kvs cst.KV) error {
//	if err := mapKVApplyDefault(dst, kvs); err != nil {
//		return err
//	}
//	return Validate(dst)
//}
