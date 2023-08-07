// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/validx"
	"reflect"
	"unsafe"
)

//func checkDestSchema(dest any, bindOpts *BindOptions) (*reflect.Value, *StructSchema, error) {
//	dstTyp := reflect.TypeOf(dest)
//	if dstTyp.Kind() != reflect.Pointer {
//		return nil, nil, errors.New("Target object must be pointer.")
//	}
//
//	dstVal := reflect.Indirect(reflect.ValueOf(dest))
//	if dstVal.Kind() != reflect.Struct {
//		return nil, nil, fmt.Errorf("%T not like struct.", dest)
//	}
//
//	return &dstVal, SchemaByType(dstVal.Type(), bindOpts), nil
//}
//
//func isInitialValue(dst reflect.Value) bool {
//	switch dst.Kind() {
//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
//		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
//		return dst.Int() == 0
//	case reflect.Float64, reflect.Float32:
//		return dst.Float() == 0
//	case reflect.String:
//		return dst.String() == ""
//	}
//	return false
//}

//func sdxSetValue(dstVal reflect.Value, src any, fOpt *fieldOptions, bindOpts *BindOptions) error {
//	return nil
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// object:
// 用传入的hash数据源，赋值目标对象，并可以做数据校验
func bindKVToStruct(dst any, kvs cst.SuperKV, bindOpts *BindOptions) (err error) {
	// 数据源和目标对象只要有一个为nil，啥都不做，也不返回错误
	if dst == nil || kvs == nil || kvs.Len() == 0 || bindOpts == nil {
		return nil
	}

	// +++++++++++++++++++++++++++++++
	dstTyp := reflect.TypeOf(dst)
	if dstTyp.Kind() != reflect.Pointer {
		return errors.New("Target object must be pointer value.")
	}
	if dstTyp.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%T not like struct.", dst)
	}

	sm := SchemaByType(dstTyp, bindOpts)
	ptr := (*rt.AFace)(unsafe.Pointer(&dst)).DataPtr
	// +++++++++++++++++++++++++++++++

	var fls []string
	if bindOpts.UseFieldName {
		fls = sm.fields
	} else {
		fls = sm.columns
	}
	flsOpts := sm.fieldsOpts

	// 两种循环方式。1：目标结构的字段  2：源字段（一般情况下，这种更好）
	for i := 0; i < len(fls); i++ {
		fOpt := flsOpts[i] // 这个肯定不为 nil
		vOpt := fOpt.valid // 这个可能是 nil
		fName := fls[i]
		sv, ok := kvs.Get(fName)

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

		//dstVal := reflect.ValueOf(dst)
		//fv := sm.RefValueByIndex(&dstVal, int8(i))
		//// 如果字段是结构体类型
		//fvType := fv.Type()
		//if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
		//	// 如果sv 无法转换成 cst.KV 类型，将要抛出异常
		//	switch sv.(type) {
		//	case map[string]any:
		//		sv = cst.KV(sv.(map[string]any))
		//	}
		//	if err = bindKVToStruct(fv.Addr().Interface(), sv.(cst.KV), bindOpts); err != nil {
		//		return err
		//	}
		//	continue
		//}

		fa := sm.FieldsAttr[i]
		fPtr := unsafe.Pointer(uintptr(ptr) + fa.Offset)
		fa.kvBinder(fPtr, sv)

		// 是否需要验证字段数据的合法性
		if bindOpts.UseValid && vOpt != nil {
			if err = validx.ValidateFieldSmart(fPtr, fa.Kind, vOpt); err != nil {
				return err
			}
		}
	}
	return nil
}

//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 主要用于给dst加上默认值，然后执行下字段验证
//func optimizeStruct(dst any, bindOpts *BindOptions) (err error) {
//	if dst == nil || bindOpts == nil {
//		return nil
//	}
//	dstVal, sm, err := checkDestSchema(dst, bindOpts)
//	if err != nil {
//		return err
//	}
//
//	for i := 0; i < len(sm.fields); i++ {
//		fv := sm.RefValueByIndex(dstVal, int8(i))
//
//		// 如果字段是结构体类型
//		fvType := fv.Type()
//		if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
//			if err = optimizeStruct(fv.Addr().Interface(), bindOpts); err != nil {
//				return err
//			}
//			continue
//		}
//
//		// 如果字段值看上去像变量刚生成后默认初始化值，那么就加载默认信息
//		fOpt := sm.fieldsOpts[i]
//		vOpt := fOpt.valid
//		if isInitialValue(fv) && bindOpts.UseDefValue && vOpt != nil {
//			if vOpt.DefValue == "" {
//				continue
//			}
//			if err = sdxSetValue(fv, vOpt.DefValue, fOpt, bindOpts); err != nil {
//				return err
//			}
//		}
//		// 是否需要验证字段数据的合法性
//		if bindOpts.UseValid && fOpt != nil {
//			if err = validx.ValidateField(&fv, vOpt); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}

//	特殊值
//
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func setPtr(p unsafe.Pointer, val any) {
	p = *((*unsafe.Pointer)(p))
	if *((*unsafe.Pointer)(p)) == nil {
	}
}

func setStruct(p unsafe.Pointer, val any) {

}

func setMap(p unsafe.Pointer, val any) {

}

func setList(p unsafe.Pointer, val any) {

}

// NOTE：不通用
// 下面的绑定函数，只针对类似Web请求提交的数据。只支持string|number|bool等基础知识，或者 KV|Array
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// int
func setInt(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt(p, v)
	case string:
		BindInt(p, lang.ParseInt(v))
	case *string:
		BindInt(p, lang.ParseInt(*v))
	}
}

func setInt8(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt8(p, v)
	case string:
		BindInt8(p, lang.ParseInt(v))
	case *string:
		BindInt8(p, lang.ParseInt(*v))
	}
}

func setInt16(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt16(p, v)
	case string:
		BindInt16(p, lang.ParseInt(v))
	case *string:
		BindInt16(p, lang.ParseInt(*v))
	}
}

func setInt32(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt32(p, v)
	case string:
		BindInt32(p, lang.ParseInt(v))
	case *string:
		BindInt32(p, lang.ParseInt(*v))
	}
}

func setInt64(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt64(p, v)
	case string:
		BindInt64(p, lang.ParseInt(v))
	case *string:
		BindInt32(p, lang.ParseInt(*v))
	}
}

// uint
func setUint(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint(p, v)
	case string:
		BindUint(p, lang.ParseUint(v))
	case *string:
		BindUint(p, lang.ParseUint(*v))
	}
}

func setUint8(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint8(p, v)
	case string:
		BindUint8(p, lang.ParseUint(v))
	case *string:
		BindUint8(p, lang.ParseUint(*v))
	}
}

func setUint16(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint16(p, v)
	case string:
		BindUint16(p, lang.ParseUint(v))
	case *string:
		BindUint16(p, lang.ParseUint(*v))
	}
}

func setUint32(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint32(p, v)
	case string:
		BindUint32(p, lang.ParseUint(v))
	case *string:
		BindUint32(p, lang.ParseUint(*v))
	}
}

func setUint64(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint64(p, v)
	case string:
		BindUint64(p, lang.ParseUint(v))
	case *string:
		BindUint64(p, lang.ParseUint(*v))
	}
}

// float
func setFloat32(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case float64:
		BindFloat32(p, v)
	case string:
		BindFloat32(p, lang.ParseFloat(v))
	case *string:
		BindFloat32(p, lang.ParseFloat(*v))
	}
}

func setFloat64(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case float64:
		BindFloat64(p, v)
	case string:
		BindFloat64(p, lang.ParseFloat(v))
	case *string:
		BindFloat64(p, lang.ParseFloat(*v))
	}
}

// string
func setString(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case string:
		BindString(p, v)
	case *string:
		BindString(p, *v)
	default:
		BindString(p, lang.ToString(v))
	}
}

func setBool(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case bool:
		BindBool(p, v)
	case string:
		BindBool(p, lang.ParseBool(v))
	case *string:
		BindBool(p, lang.ParseBool(*v))
	}
}

func setAny(p unsafe.Pointer, val any) {
	BindAny(p, val)
}
