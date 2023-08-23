// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package bind

import (
	"errors"
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
		return errors.New(dstTyp.String() + " must be struct.")
	}
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	ptr := (*rt.AFace)(unsafe.Pointer(&dst)).DataPtr
	return bindKVToStructInner(ptr, dstTyp, kvs, opts)
}

func bindKVToStructInner(ptr unsafe.Pointer, dstT reflect.Type, kvs cst.SuperKV, opts *dts.BindOptions) (err error) {
	sm := dts.SchemaByType(dstT, opts)

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
				return errors.New("the field must required: " + fName)
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
			if err = bindMap(fPtr, fv); err != nil {
				return
			}
			continue
		case reflect.Pointer:
			if err = bindPointer(fPtr, fv); err != nil {
				return
			}
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

// 绑定列表数据
func bindList(ptr unsafe.Pointer, dstT reflect.Type, val any, opts *dts.BindOptions) (err error) {
	// 因为绑定数据来源于JSON，YAML等数据的解析，这类数据在遇到数组时候，几乎都是用 []any 表示
	list, ok := val.([]any)
	if !ok {
		return errors.New("dts: value type must be []any")
	}

	var itSize uintptr
	ct := len(list)

	if dstT.Kind() == reflect.Array {
		if dstT.Len() != ct {
			return errors.New("dts: array length not match.")
		}

		dstT = dstT.Elem()
		itSize = dstT.Size()
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
		ptr = unsafe.Pointer(sh.Data)
	}

	// TODO: 完善这里可能出现的情况
	itKind := dstT.Kind()
	switch itKind {
	case reflect.Struct:
		for i := 0; i < ct; i++ {
			itPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*itSize)
			if err = bindStruct(itPtr, dstT, list[i], opts); err != nil {
				return
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < ct; i++ {
			itPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*itSize)
			if err = bindList(itPtr, dstT, list[i], opts); err != nil {
				return
			}
		}
	case reflect.Map:
		// TODO: bindMap
	case reflect.Pointer:
		// TODO: bindPointer
	default:
		for i := 0; i < ct; i++ {
			itPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*itSize)
			dts.BindBaseValueAsConfig(itKind, itPtr, list[i])
		}
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

func bindMap(ptr unsafe.Pointer, val any) (err error) {
	return
}

func bindPointer(ptr unsafe.Pointer, val any) (err error) {
	return
}
