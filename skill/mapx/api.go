package mapx

import (
	"github.com/qinchende/gofast/cst"
	"io"
)

// cst.KV
func ApplyKVByName(dst any, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, true, false)
}

func ApplyKVByNameWithDef(dst any, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, true, true)
}

func ApplyKVByTag(dst any, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, false, false)
}

func ApplyKVByTagWithDef(dst any, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, false, true)
}

// JSON
func ApplyJsonReader(dst any, reader io.Reader) error {
	return decodeJsonReader(dst, reader)
}

func ApplyJsonBytes(dst any, content []byte) error {
	return decodeJsonBytes(dst, content)
}

// Yaml
func ApplyYamlReader(dst any, reader io.Reader) error {
	return decodeYamlReader(dst, reader)
}

func ApplyYamlBytes(dst any, content []byte) error {
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
