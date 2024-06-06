// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"fmt"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/core/cst"
	"reflect"
	"unsafe"
)

// fetch StructSchema
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Schema(obj any, opts *BindOptions) *StructSchema {
	return SchemaByType(reflect.TypeOf(obj), opts)
}

func SchemaAsDB(obj any) *StructSchema {
	return Schema(obj, dbStructOptions)
}

func SchemaAsReq(obj any) *StructSchema {
	return Schema(obj, reqStructOptions)
}

func SchemaAsConfig(obj any) *StructSchema {
	return Schema(obj, cfgStructOptions)
}

func SchemaByType(typ reflect.Type, opts *BindOptions) *StructSchema {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		cst.PanicString(fmt.Sprintf("%T is not like struct", typ))
	}

	if opts.CacheSchema {
		return fetchSchemaCache(typ, opts)
	} else {
		return buildStructSchema(typ, opts)
	}
}

func SchemaAsDBByType(typ reflect.Type) *StructSchema {
	return SchemaByType(typ, dbStructOptions)
}

func SchemaAsReqByType(typ reflect.Type) *StructSchema {
	return SchemaByType(typ, reqStructOptions)
}

func SchemaAsConfigByType(typ reflect.Type) *StructSchema {
	return SchemaByType(typ, cfgStructOptions)
}

// reflect apis
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *StructSchema) ValueByIndex(rVal *reflect.Value, idx int8) any {
	return rVal.FieldByIndex(ss.FieldsAttr[idx].RefIndex).Interface()
}

func (ss *StructSchema) AddrByIndex(rVal *reflect.Value, idx int8) any {
	return rVal.FieldByIndex(ss.FieldsAttr[idx].RefIndex).Addr().Interface()
}

func (ss *StructSchema) RefValueByIndex(rVal *reflect.Value, idx int8) reflect.Value {
	idxArr := ss.FieldsAttr[idx].RefIndex
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
	return ss.Columns[idx]
}

func (ss *StructSchema) FieldName(idx int) string {
	return ss.Fields[idx]
}

// Gson useful apis
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (ss *StructSchema) CTips() (string, []uint8) {
//	return strings.Join(ss.cTips.items, ","), ss.cTips.idxes
//}
//
//func (ss *StructSchema) FTips() (string, []uint8) {
//	return strings.Join(ss.fTips.items, ","), ss.fTips.idxes
//}

func (ss *StructSchema) CTips() ([]string, []uint8) {
	return ss.cTips.items, ss.cTips.idxes
}

func (ss *StructSchema) FTips() ([]string, []uint8) {
	return ss.fTips.items, ss.fTips.idxes
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

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//

func (fa *fieldAttr) MyPtr(structPtr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Add(structPtr, fa.Offset)
}

func PeelPtr(ptr unsafe.Pointer, level uint8) unsafe.Pointer {
	if level == 0 || ptr == nil {
		return ptr
	}
	for level > 0 {
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			break
		}
		level--
	}
	return ptr
}
