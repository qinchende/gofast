// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package bind

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/validx"
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

// object:
// 用传入的hash数据源，赋值目标对象，并可以做数据校验
func bindKVToStruct(dst any, kvs cst.SuperKV, opts *dts.BindOptions) error {
	// 数据源和目标对象只要有一个为nil，啥都不做，也不返回错误
	if dst == nil || kvs == nil || kvs.Len() == 0 || opts == nil {
		return nil
	}

	// 以下是必要的检查
	dstTyp := reflect.TypeOf(dst)
	if dstTyp.Kind() != reflect.Pointer {
		return errors.New("Target object must be pointer value.")
	}
	dstTyp = dstTyp.Elem()
	if dstTyp.Kind() != reflect.Struct {
		return fmt.Errorf("%T must be struct.", dst)
	}
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	ptr := (*rt.AFace)(unsafe.Pointer(&dst)).DataPtr
	return bindKVToStructInner(ptr, dstTyp, kvs, opts)
}

func bindKVToStructInner(ptr unsafe.Pointer, dstType reflect.Type, kvs cst.SuperKV, opts *dts.BindOptions) (err error) {
	sm := dts.SchemaByType(dstType, opts)

	var fls []string
	if opts.UseFieldName {
		fls = sm.Fields
	} else {
		fls = sm.Columns
	}

	// 两种循环方式。1：目标结构的字段  2：源字段（一般情况下，这种更好）
	for i := 0; i < len(fls); i++ {
		fa := &sm.FieldsAttr[i] // 这个肯定不能为 nil
		vOpt := fa.Valid        // 这个可能是 nil
		fName := fls[i]
		fv, ok := kvs.Get(fName)

		if ok == false || fv == nil {
			if vOpt == nil {
				continue
			}
			if vOpt.Required && opts.UseValid {
				return fmt.Errorf("field %s requied", fName)
			}
			if opts.UseDefValue {
				fv = vOpt.DefValue
				if fv == "" {
					continue
				}
			}
		}

		fPtr := unsafe.Pointer(uintptr(ptr) + fa.Offset)
		// TODO: 完善这里可能出现的情况
		switch fa.Kind {
		case reflect.Struct:
			if err = bindStruct(fPtr, fa.Type, fv, opts); err != nil {
				return
			}
			continue
		case reflect.Array, reflect.Slice:
			if err = bindList(fPtr, fa.Type, fv, opts); err != nil {
				return
			}
			continue
		case reflect.Map:
			continue
		case reflect.Pointer:
			continue
		default:
			if fa.KVBinder == nil {
				continue
			}
			// 绑定基础数据类型（number, string, bool）
			fa.KVBinder(fPtr, fv)
		}

		// 是否需要验证字段数据的合法性
		if opts.UseValid && vOpt != nil {
			if err = validx.ValidateFieldSmart(fPtr, fa.Kind, vOpt); err != nil {
				return
			}
		}
	}
	return
}

func bindList(ptr unsafe.Pointer, dstT reflect.Type, val any, opts *dts.BindOptions) (err error) {
	switch v := val.(type) {
	case []any:
		var itSize uintptr
		var startPtr unsafe.Pointer
		ct := len(v)

		if dstT.Kind() == reflect.Array {
			dstT = dstT.Elem()
			if dstT.Len() != ct {
				return errors.New("dts: array length not match.")
			}

			itSize = dstT.Size()
			startPtr = ptr
		} else {
			dstT = dstT.Elem()
			itSize = dstT.Size()

			sh := (*reflect.SliceHeader)(ptr)
			if sh.Cap < ct {
				newMem := make([]byte, int(itSize)*ct)
				sh.Data = (*reflect.SliceHeader)(unsafe.Pointer(&newMem)).Data
				sh.Len, sh.Cap = ct, ct
			} else {
				sh.Len = ct
			}
			startPtr = unsafe.Pointer(sh.Data)
		}

		// 处理每一项值的绑定
		for i := 0; i < ct; i++ {
			itPtr := unsafe.Pointer(uintptr(startPtr) + uintptr(i)*itSize)
			itVal := v[i]

			// TODO: 完善这里可能出现的情况
			dstKind := dstT.Kind()
			switch dstKind {
			case reflect.Struct:
				if err = bindStruct(itPtr, dstT, itVal, opts); err != nil {
					return
				}
			case reflect.Array, reflect.Slice:
				if err = bindList(itPtr, dstT, itVal, opts); err != nil {
					return
				}
			default:
				dts.BindBaseValueAsConfig(dstKind, itPtr, itVal)
			}
		}
	default:
		return errors.New("dts: only array-like value supported.")
	}
	return
}

func bindStruct(ptr unsafe.Pointer, dstT reflect.Type, val any, opts *dts.BindOptions) (err error) {
	var skv cst.SuperKV

	switch v := val.(type) {
	case map[string]any:
		skv = cst.KV(v)
	default:
		return
	}
	if err = bindKVToStructInner(ptr, dstT, skv, opts); err != nil {
		return
	}
	return
}

// 特别类型
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func setPtr(p unsafe.Pointer, val any) {
	p = *((*unsafe.Pointer)(p))
	if *((*unsafe.Pointer)(p)) == nil {
	}
}

func setMap(p unsafe.Pointer, val any) {

}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
