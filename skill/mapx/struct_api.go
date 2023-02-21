// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"reflect"
)

// no cache
func SchemaNoCache(obj any, opts *ApplyOptions) *GfStruct {
	return structSchema(reflect.TypeOf(obj), opts)
}

func SchemaNoCacheOfType(rTyp reflect.Type, opts *ApplyOptions) *GfStruct {
	return structSchema(rTyp, opts)
}

// cached
func Schema(obj any, opts *ApplyOptions) *GfStruct {
	return fetchSchemaCache(reflect.TypeOf(obj), opts)
}

func SchemaOfType(rTyp reflect.Type, opts *ApplyOptions) *GfStruct {
	return fetchSchemaCache(rTyp, opts)
}

// reflect
func (ms *GfStruct) ValueByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Interface()
}

func (ms *GfStruct) AddrByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Addr().Interface()
}

func (ms *GfStruct) RefValueByIndex(rVal *reflect.Value, index int8) reflect.Value {
	idxArr := ms.fieldsIndex[index]
	if len(idxArr) == 1 {
		return rVal.Field(idxArr[0])
	}
	tmpVal := *rVal
	for _, x := range idxArr {
		tmpVal = tmpVal.Field(x)
	}
	return tmpVal
}
