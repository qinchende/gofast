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
	dstTyp := reflect.TypeOf(dest)
	if dstTyp.Kind() != reflect.Pointer {
		cst.PanicString("Target object must be pointer.")
	}
	sliceType := dstTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		cst.PanicString("Target object must be slice.")
	}

	isPtr := false
	isKV := false
	recordType := sliceType.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if recordType.Kind() == reflect.Pointer {
		isPtr = true
		recordType = recordType.Elem()
	} else {
		typName := recordType.Name()
		if typName == "cst.KV" || typName == "KV" {
			isKV = true
		}
	}

	return sliceType, recordType, isPtr, isKV
}

func checkDestSchema(dest any, bindOpts *BindOptions) (*reflect.Value, *StructSchema, error) {
	dstTyp := reflect.TypeOf(dest)
	if dstTyp.Kind() != reflect.Pointer {
		return nil, nil, errors.New("Target object must be pointer.")
	}

	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	if dstVal.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("%T not like struct.", dest)
	}

	return &dstVal, SchemaOfType(dstVal.Type(), bindOpts), nil
}
