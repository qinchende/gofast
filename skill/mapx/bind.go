// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/validx"
	"reflect"
)

// 只用传入的值赋值对象
func bindKVToStruct(dst any, kvs cst.KV, bindOpts *BindOptions) (err error) {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}
	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not like struct", dst)
	}
	sm := getSchema(dstVal, bindOpts)

	var fls []string
	if bindOpts.UseFieldName {
		fls = sm.fields
	} else {
		fls = sm.columns
	}
	flsOpts := sm.fieldsOpts

	for i := 0; i < len(fls); i++ {
		fOpt := flsOpts[i] // 这个肯定不为 nil
		vOpt := fOpt.valid // 这个可能是 nil
		fName := fls[i]
		sv, ok := kvs[fName]

		if ok == false {
			if vOpt == nil {
				continue
			}
			if vOpt.Required && bindOpts.UseValid {
				return fmt.Errorf("field %s requied", fName)
			}
			if bindOpts.UseDefValue {
				sv = vOpt.DefValue
				if sv == "" {
					continue
				}
			}
		}
		fv := sm.RefValueByIndex(&dstVal, int8(i))

		// 如果字段是结构体类型
		fvType := fv.Type()
		if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
			// 如果sv 无法转换成 cst.KV 类型，将要抛出异常
			if err = bindKVToStruct(fv.Addr().Interface(), sv.(map[string]any), bindOpts); err != nil {
				return err
			}
			continue
		}

		if err = sdxSetValue(fv, sv, fOpt, bindOpts); err != nil {
			return err
		}

		// 是否需要验证字段数据的合法性
		if bindOpts.UseValid && vOpt != nil {
			if err = validx.ValidateField(&fv, vOpt); err != nil {
				return err
			}
		}
	}
	return nil
}

// src 只能是 array, slice 类型。如果是 string ，先按照JSON格式解析成数组
func bindList(dst any, src any, fOpt *fieldOptions, bindOpts *BindOptions) (err error) {
	if fOpt == nil {
		return fmt.Errorf("field options can't nil.")
	}

	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	srcVal := reflect.Indirect(reflect.ValueOf(src))
	// 如果数据源是字符串，先按照JSON解析成数组
	if srcVal.Kind() == reflect.String {
		var srcNew []any
		if err = jsonx.UnmarshalFromString(&srcNew, src.(string)); err != nil {
			return err
		}
		src = srcNew
		srcVal = reflect.Indirect(reflect.ValueOf(src))
	}

	dstKind := dstVal.Kind()
	srcKind := dstVal.Kind()

	switch {
	case (dstKind == reflect.Slice || dstKind == reflect.Array) && (srcKind == reflect.Slice || srcKind == reflect.Array):
		// NOTE: 这里可能 dstVal.Len() > srcVal.Len() 也应该支持
		if dstKind == reflect.Array && dstVal.Len() != srcVal.Len() {
			return errors.New("array length not match")
		}

		sliceTyp, itemType, isPtr, _ := checkDestType(dst)
		dstNew := reflect.MakeSlice(sliceTyp, srcVal.Len(), srcVal.Len())
		dstVal.Set(dstNew)
		for i := 0; i < srcVal.Len(); i++ {
			fv := dstVal.Index(i)
			if isPtr {
				fv.Set(reflect.New(itemType))
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				if err = bindKVToStruct(fv.Addr().Interface(), srcVal.Index(i).Interface().(map[string]any), bindOpts); err != nil {
					return err
				}
				continue
			}

			if err = sdxSetValue(fv, srcVal.Index(i).Interface(), fOpt, bindOpts); err != nil {
				return err
			}
			// 是否需要验证字段数据的合法性
			if bindOpts.UseValid && fOpt.valid != nil {
				if err = validx.ValidateField(&fv, fOpt.valid); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("only array-like value supported")
	}

	// 数组不能为空
	if bindOpts.UseValid && fOpt.valid != nil && fOpt.valid.Required && dstVal.Len() == 0 {
		return fmt.Errorf("list field %s requied", dstVal.Type().String())
	}

	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主要用于给dst加上默认值，然后执行下字段验证
func optimizeStruct(dst any, bindOpts *BindOptions) (err error) {
	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not like struct", dst)
	}
	sm := getSchema(dstVal, bindOpts)

	for i := 0; i < len(sm.fields); i++ {
		fv := sm.RefValueByIndex(&dstVal, int8(i))

		// 如果字段是结构体类型
		fvType := fv.Type()
		if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
			if err = optimizeStruct(fv.Addr().Interface(), bindOpts); err != nil {
				return err
			}
			continue
		}

		// 如果字段值看上去像变量刚生成后默认初始化值，那么就加载默认信息
		fOpt := sm.fieldsOpts[i]
		vOpt := fOpt.valid
		if isInitialValue(fv) && bindOpts.UseDefValue && vOpt != nil {
			if vOpt.DefValue == "" {
				continue
			}
			if err = sdxSetValue(fv, vOpt.DefValue, fOpt, bindOpts); err != nil {
				return err
			}
		}
		// 是否需要验证字段数据的合法性
		if bindOpts.UseValid && fOpt != nil {
			if err = validx.ValidateField(&fv, vOpt); err != nil {
				return err
			}
		}
	}
	return nil
}
