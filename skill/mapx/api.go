package mapx

import (
	"github.com/qinchende/gofast/cst"
	"io"
)

// cst.KV
func ApplyKVByName(dst interface{}, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, true, false)
}

func ApplyKVByNameWithDef(dst interface{}, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, true, true)
}

func ApplyKVByTag(dst interface{}, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, false, false)
}

func ApplyKVByTagWithDef(dst interface{}, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, false, true)
}

// JSON
func ApplyJsonReader(dst interface{}, reader io.Reader) error {
	return decodeJsonReader(dst, reader)
}

func ApplyJsonBytes(dst interface{}, content []byte) error {
	return decodeJsonBytes(dst, content)
}

// Yaml
func ApplyYamlReader(dst interface{}, reader io.Reader) error {
	return decodeYamlReader(dst, reader)
}

func ApplyYamlBytes(dst interface{}, content []byte) error {
	return decodeYamlBytes(dst, content)
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
