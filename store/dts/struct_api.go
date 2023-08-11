// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"strings"
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
	return Schema(obj, reqStructOptions)
}

func SchemaForConfig(obj any) *StructSchema {
	return Schema(obj, cfgStructOptions)
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
	return SchemaByType(rTyp, reqStructOptions)
}

func SchemaForConfigByType(rTyp reflect.Type) *StructSchema {
	return SchemaByType(rTyp, cfgStructOptions)
}

// reflect apis
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *StructSchema) ValueByIndex(rVal *reflect.Value, idx int8) any {
	return rVal.FieldByIndex(ss.FieldsAttr[idx].rIndex).Interface()
}

func (ss *StructSchema) AddrByIndex(rVal *reflect.Value, idx int8) any {
	return rVal.FieldByIndex(ss.FieldsAttr[idx].rIndex).Addr().Interface()
}

func (ss *StructSchema) RefValueByIndex(rVal *reflect.Value, idx int8) reflect.Value {
	idxArr := ss.FieldsAttr[idx].rIndex
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

func (ss *StructSchema) ColumnName(idx int) string {
	return ss.columns[idx]
}

func (ss *StructSchema) FieldName(idx int) string {
	return ss.fields[idx]
}

// Gson
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *StructSchema) CTips() (string, []uint8) {
	return strings.Join(ss.cTips.items, ","), ss.cTips.idxes
}

func (ss *StructSchema) FTips() (string, []uint8) {
	return strings.Join(ss.fTips.items, ","), ss.fTips.idxes
}

func (ss *StructSchema) CIndexes(cls []string) (ret []uint8) {
	ret = make([]uint8, len(cls))
	for i := range cls {
		ret[i] = uint8(ss.ColumnIndex(cls[i]))
	}
	return
}

func (ss *StructSchema) FIndexes(fls []string) (ret []uint8) {
	ret = make([]uint8, len(fls))
	for i := range fls {
		ret[i] = uint8(ss.FieldIndex(fls[i]))
	}
	return
}
