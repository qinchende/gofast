// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"reflect"
)

func checkDestType(dest any) (reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		cst.PanicString("Target object must be pointer.")
	}
	sliceType := dTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		cst.PanicString("Target object must be slice.")
	}

	isPtr := false
	isKV := false
	recordType := sliceType.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if recordType.Kind() == reflect.Ptr {
		isPtr = true
		recordType = recordType.Elem()
	} else {
		typName := recordType.Name()
		if typName == "cst.KV" || typName == "fst.KV" || typName == "KV" {
			isKV = true
		}
	}

	return sliceType, recordType, isPtr, isKV
}

func checkDestSchema(dst any, bindOpts *BindOptions) (*reflect.Value, *StructSchema, error) {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return nil, nil, errors.New("Target object must be pointer.")
	}

	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	if dstVal.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("%T not like struct.", dst)
	}

	if bindOpts.CacheSchema {
		return &dstVal, SchemaOfType(dstVal.Type(), bindOpts), nil
	} else {
		return &dstVal, SchemaNoCacheOfType(dstVal.Type(), bindOpts), nil
	}
}
