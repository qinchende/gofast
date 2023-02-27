// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"reflect"
)

// no cache
func SchemaNoCache(obj any, opts *BindOptions) *StructSchema {
	return buildStructSchema(reflect.TypeOf(obj), opts)
}

func SchemaNoCacheOfType(rTyp reflect.Type, opts *BindOptions) *StructSchema {
	return buildStructSchema(rTyp, opts)
}

// cached
func Schema(obj any, opts *BindOptions) *StructSchema {
	return fetchSchemaCache(reflect.TypeOf(obj), opts)
}

func SchemaOfType(rTyp reflect.Type, opts *BindOptions) *StructSchema {
	return fetchSchemaCache(rTyp, opts)
}

// reflect
func (ms *StructSchema) ValueByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Interface()
}

func (ms *StructSchema) AddrByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ms.fieldsIndex[index]).Addr().Interface()
}

func (ms *StructSchema) RefValueByIndex(rVal *reflect.Value, index int8) reflect.Value {
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
