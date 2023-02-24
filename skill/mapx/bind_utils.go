// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import "reflect"

func checkDestType(dest any) (reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		panic("dest must be pointer.")
	}
	dSliceTyp := dTyp.Elem()
	if dSliceTyp.Kind() != reflect.Slice {
		panic("dest must be slice.")
	}

	isPtr := false
	isKV := false
	dItemType := dSliceTyp.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if dItemType.Kind() == reflect.Ptr {
		isPtr = true
		dItemType = dItemType.Elem()
	} else if dItemType.Name() == "KV" {
		isKV = true
	}

	return dSliceTyp, dItemType, isPtr, isKV
}

func getSchema(dstVal reflect.Value, bindOpts *BindOptions) *GfStruct {
	if bindOpts.CacheSchema {
		return SchemaOfType(dstVal.Type(), bindOpts)
	} else {
		return SchemaNoCacheOfType(dstVal.Type(), bindOpts)
	}
}
