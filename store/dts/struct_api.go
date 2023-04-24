// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
)

func Schema(obj any, opts *BindOptions) *StructSchema {
	if opts.CacheSchema {
		return fetchSchemaCache(reflect.TypeOf(obj), opts)
	} else {
		return buildStructSchema(reflect.TypeOf(obj), opts)
	}
}

func SchemaForDB(obj any) *StructSchema {
	return Schema(obj, dbStructOptions)
}

func SchemaForInput(obj any) *StructSchema {
	return Schema(obj, inputStructOptions)
}

func SchemaForConfig(obj any) *StructSchema {
	return Schema(obj, configStructOptions)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func SchemaByType(rTyp reflect.Type, opts *BindOptions) *StructSchema {
	if opts.CacheSchema {
		return fetchSchemaCache(rTyp, opts)
	} else {
		return buildStructSchema(rTyp, opts)
	}
}

func SchemaForDBByType(rTyp reflect.Type) *StructSchema {
	return SchemaByType(rTyp, dbStructOptions)
}

func SchemaForInputByType(rTyp reflect.Type) *StructSchema {
	return SchemaByType(rTyp, inputStructOptions)
}

func SchemaForConfigByType(rTyp reflect.Type) *StructSchema {
	return SchemaByType(rTyp, configStructOptions)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// reflect
func (ss *StructSchema) ValueByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ss.fieldsIndex[index]).Interface()
}

func (ss *StructSchema) AddrByIndex(rVal *reflect.Value, index int8) any {
	return rVal.FieldByIndex(ss.fieldsIndex[index]).Addr().Interface()
}

func (ss *StructSchema) RefValueByIndex(rVal *reflect.Value, index int8) reflect.Value {
	idxArr := ss.fieldsIndex[index]
	if len(idxArr) == 1 {
		return rVal.Field(idxArr[0])
	}
	tmpVal := *rVal
	for _, x := range idxArr {
		tmpVal = tmpVal.Field(x)
	}
	return tmpVal
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *StructSchema) KeyIndex(k string) int {
	idx := lang.SearchSortStrings(ss.columnsKV.items, k)
	if idx < 0 {
		return -1
	}
	return ss.columnsKV.idxes[idx]
}
