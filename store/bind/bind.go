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
		return errors.New("Dest object must be pointer value.")
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
func bindList(ptr unsafe.Pointer, dstT reflect.Type, src any, opts *dts.BindOptions) (err error) {
	// 因为绑定数据来源于JSON，YAML等数据的解析，这类数据在遇到数组时候，几乎都是用 []any 表示
	list, ok := src.([]any)
	if !ok {
		return errors.New("dts: value type must be []any")
	}

	var dstSize uintptr
	listLen := len(list)

	if dstT.Kind() == reflect.Array {
		if dstT.Len() != listLen {
			return errors.New("dts: array length not match.")
		}

		dstT = dstT.Elem()
		dstSize = dstT.Size()
	} else {
		dstT = dstT.Elem()
		dstSize = dstT.Size()

		sh := (*rt.SliceHeader)(ptr)
		if sh.Cap < listLen {
			newMem := make([]byte, int(dstSize)*listLen)
			sh.DataPtr = (*rt.SliceHeader)(unsafe.Pointer(&newMem)).DataPtr
			sh.Len, sh.Cap = listLen, listLen
		} else {
			sh.Len = listLen
		}
		ptr = sh.DataPtr
	}

	// TODO: 完善这里可能出现的情况
	itKind := dstT.Kind()
	switch itKind {
	case reflect.Struct:
		for i := 0; i < listLen; i++ {
			itPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*dstSize)
			if err = bindStruct(itPtr, dstT, list[i], opts); err != nil {
				return
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < listLen; i++ {
			itPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*dstSize)
			if err = bindList(itPtr, dstT, list[i], opts); err != nil {
				return
			}
		}
	case reflect.Map:
		// TODO: bindMap
	case reflect.Pointer:
		// TODO: bindPointer
	default:
		for i := 0; i < listLen; i++ {
			itPtr := unsafe.Pointer(uintptr(ptr) + uintptr(i)*dstSize)
			dts.BindBaseValueAsConfig(itKind, itPtr, list[i])
		}
	}
	return
}

func bindStruct(ptr unsafe.Pointer, dstT reflect.Type, src any, opts *dts.BindOptions) error {
	switch v := src.(type) {
	case map[string]any:
		return bindKVToStructInner(ptr, dstT, cst.KV(v), opts)
	default:
		return nil
	}
}

func bindMap(ptr unsafe.Pointer, val any) (err error) {
	return
}

func bindPointer(ptr unsafe.Pointer, val any) (err error) {
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主要用于给dst加上默认值，然后执行下字段验证
func optimizeStruct(dst any, opts *dts.BindOptions) (err error) {
	if dst == nil || opts == nil {
		return nil
	}

	// 以下是必要的检查
	dstTyp := reflect.TypeOf(dst)
	if dstTyp.Kind() != reflect.Pointer {
		return errors.New("Dest object must be pointer value.")
	}
	dstTyp = dstTyp.Elem()
	if dstTyp.Kind() != reflect.Struct {
		return errors.New(dstTyp.String() + " must be struct.")
	}

	ptr := (*rt.AFace)(unsafe.Pointer(&dst)).DataPtr
	return optimizeStructInner(ptr, dstTyp, opts)
}

func optimizeStructInner(ptr unsafe.Pointer, dstT reflect.Type, opts *dts.BindOptions) (err error) {
	sm := dts.SchemaByType(dstT, opts)

	for i := 0; i < len(sm.Fields); i++ {
		fa := &sm.FieldsAttr[i] // 这个肯定不能为 nil
		fKind := fa.Kind
		fPtr := unsafe.Pointer(uintptr(ptr) + fa.Offset)

		// 如果字段是结构体类型
		if fKind == reflect.Struct && fa.Type != cst.TypeTime {
			if err = optimizeStructInner(fPtr, dstT, opts); err != nil {
				return
			}
			continue
		}

		// 如果字段值看上去像变量刚生成后默认初始化值，那么就加载默认信息
		vOpt := fa.Valid
		if isInitialValue(fKind, fPtr) && opts.UseDefValue && vOpt != nil {
			if vOpt.DefValue == "" || fa.KVBinder == nil {
				continue
			}
			fa.KVBinder(fPtr, vOpt.DefValue)
		}

		// Check: 是否需要验证字段数据的合法性
		if opts.UseValid && vOpt != nil {
			if err = validx.ValidateFieldSmart(fPtr, fKind, vOpt); err != nil {
				return
			}
		}
	}
	return nil
}

func isInitialValue(kd reflect.Kind, p unsafe.Pointer) bool {
	switch kd {
	case reflect.Int:
		return *(*int)(p) == 0
	case reflect.Int8:
		return *(*int8)(p) == 0
	case reflect.Int16:
		return *(*int16)(p) == 0
	case reflect.Int32:
		return *(*int32)(p) == 0
	case reflect.Int64:
		return *(*int64)(p) == 0

	case reflect.Uint:
		return *(*uint)(p) == 0
	case reflect.Uint8:
		return *(*uint8)(p) == 0
	case reflect.Uint16:
		return *(*uint16)(p) == 0
	case reflect.Uint32:
		return *(*uint32)(p) == 0
	case reflect.Uint64:
		return *(*uint64)(p) == 0

	case reflect.Float32:
		return *(*float32)(p) == 0
	case reflect.Float64:
		return *(*float64)(p) == 0

	case reflect.String:
		return *(*string)(p) == ""
	case reflect.Bool:
		return *(*bool)(p) == false
	case reflect.Interface:
		return *(*any)(p) == nil

	default:
		return false
	}
}
