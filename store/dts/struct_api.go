// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
)

// fetch StructSchema
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

// ++++++++++++++++++++++++++
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

// reflect apis
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *StructSchema) ValueByIndex(rVal *reflect.Value, idx int8) any {
	return rVal.FieldByIndex(ss.fieldsIndex[idx]).Interface()
}

func (ss *StructSchema) AddrByIndex(rVal *reflect.Value, idx int8) any {
	return rVal.FieldByIndex(ss.fieldsIndex[idx]).Addr().Interface()
}

func (ss *StructSchema) RefValueByIndex(rVal *reflect.Value, idx int8) reflect.Value {
	idxArr := ss.fieldsIndex[idx]
	if len(idxArr) == 1 {
		return rVal.Field(idxArr[0])
	}
	tmpVal := *rVal
	for _, x := range idxArr {
		tmpVal = tmpVal.Field(x)
	}
	return tmpVal
}

// Quick search for structure fields
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *StructSchema) ColumnIndex(k string) int {
	kv := ss.cTips
	if idx := lang.SearchSortedSkip(kv.items, int(kv.lenOff[len(k)]), k); idx < 0 {
		return -1
	} else {
		return int(kv.idxes[idx])
	}
}

func (ss *StructSchema) FieldIndex(k string) int {
	kv := ss.fTips
	if idx := lang.SearchSortedSkip(kv.items, int(kv.lenOff[len(k)]), k); idx < 0 {
		return -1
	} else {
		return int(kv.idxes[idx])
	}
}
