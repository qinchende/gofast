// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package bind

import (
	"errors"
	"github.com/qinchende/gofast/aid/validx"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/dts"
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"unsafe"
)

func checkStructDest(dst any) (dstTyp reflect.Type, ptr unsafe.Pointer, err error) {
	// 以下是必要的检查
	dstTyp = reflect.TypeOf(dst)
	if dstTyp.Kind() != reflect.Pointer {
		return nil, nil, errors.New("Dest object must be pointer value.")
	}
	dstTyp = dstTyp.Elem()
	if dstTyp.Kind() != reflect.Struct {
		return nil, nil, errors.New(dstTyp.String() + " must be struct.")
	}
	ptr = (*rt.AFace)(unsafe.Pointer(&dst)).DataPtr
	return
}

// object:
// 用传入的hash数据源，赋值目标对象，并可以做数据校验
func bindKVToStruct(dst any, kvs cst.SuperKV, opts *dts.BindOptions) error {
	// 数据源和目标对象只要有一个为nil，啥都不做，也不返回错误
	if dst == nil || opts == nil {
		return errors.New("has nil param.")
	}
	if dstType, ptr, err := checkStructDest(dst); err != nil {
		return err
	} else {
		return bindKVToStructIter(ptr, dstType, kvs, opts)
	}
}

func bindKVToStructIter(ptr unsafe.Pointer, dstT reflect.Type, kvs cst.SuperKV, opts *dts.BindOptions) (err error) {
	sm := dts.SchemaByType(dstT, opts)

	var fls []string
	if opts.UseFieldName {
		fls = sm.Fields
	} else {
		fls = sm.Columns
	}

	// 两种循环方式。1：目标结构的字段  2：源字段（一般情况下，这种更好）
	for i := 0; i < len(fls); i++ {
		fa := &sm.FieldsAttr[i]   // 肯定不是nil
		vOpt := fa.Valid          // 可能是nil
		fv, ok := kvs.Get(fls[i]) // 查找字段值
		fPtr := fa.MyPtr(ptr)

		// 没有找到字段，或者值为nil，那么就看看是否有默认值
		if ok == false || fv == nil {
			if opts.UseDefValue && vOpt != nil {
				fv = vOpt.DefValue
				if fv == "" {
					goto validField
				}
			}
		}

		// TODO: 完善这里可能出现的情况, fPtr 可能为 nil
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
		default:
			if fa.KVBinder != nil {
				fa.KVBinder(fPtr, fv) // 绑定基础数据类型（number, string, bool）
			}
		}

	validField:
		// 是否需要验证字段数据的合法性
		if opts.UseValid && vOpt != nil {
			if err = validx.ValidateFieldPtr(fPtr, fa.Kind, vOpt); err != nil {
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

	srcLen := len(list)
	dstKind := dstT.Kind()

	if dstKind == reflect.Array && dstT.Len() != srcLen {
		return errors.New("dts: array length not match.")
	}

	itemType := dstT.Elem()
	itemBytes := int(itemType.Size())

	if dstKind == reflect.Slice {
		ptr = rt.SliceToArray(ptr, itemBytes, srcLen)
	}

	switch itemKind := itemType.Kind(); itemKind {
	case reflect.Struct:
		for i := 0; i < srcLen; i++ {
			itPtr := unsafe.Add(ptr, i*itemBytes)
			if err = bindStruct(itPtr, itemType, list[i], opts); err != nil {
				return
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < srcLen; i++ {
			itPtr := unsafe.Add(ptr, i*itemBytes)
			if err = bindList(itPtr, itemType, list[i], opts); err != nil {
				return
			}
		}
	case reflect.Map:
		// TODO: bindMap
	case reflect.Pointer:
		// TODO: bindPointer
	default:
		for i := 0; i < srcLen; i++ {
			itPtr := unsafe.Add(ptr, i*itemBytes)
			dts.BindBaseValueAsConfig(itemKind, itPtr, list[i])
		}
	}
	return
}

func bindStruct(ptr unsafe.Pointer, dstT reflect.Type, src any, opts *dts.BindOptions) error {
	switch v := src.(type) {
	case map[string]any:
		return bindKVToStructIter(ptr, dstT, cst.KV(v), opts)
	default:
		// TODO: 其它数据源暂时忽略，不执行绑定，也不报错
		return nil
	}
}

func bindMap(ptr unsafe.Pointer, val any) (err error) {
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主要用于给dst加上默认值，然后执行每个字段规则验证
func optimizeStruct(dst any, opts *dts.BindOptions) error {
	if dst == nil || opts == nil {
		return errors.New("has nil param.")
	}
	if dstType, ptr, err := checkStructDest(dst); err != nil {
		return err
	} else {
		return optimizeStructIter(ptr, dstType, opts)
	}
}

func optimizeStructIter(ptr unsafe.Pointer, dstT reflect.Type, opts *dts.BindOptions) (err error) {
	sm := dts.SchemaByType(dstT, opts)

	for i := 0; i < len(sm.Fields); i++ {
		fa := &sm.FieldsAttr[i] // 肯定不会是nil
		fKind := fa.Kind
		fPtr := fa.MyPtr(ptr)

		// 如果字段是结构体类型
		if fKind == reflect.Struct && fa.Type != cst.TypeTime {
			// Note：指向其它结构体的字段，不做处理，防止嵌套循环验证
			if fa.PtrLevel > 0 {
				continue
			}
			if err = optimizeStructIter(fPtr, fa.Type, opts); err != nil {
				return
			}
			continue
		}

		vOpt := fa.Valid
		if vOpt == nil {
			continue
		}

		fPtr = dts.PeelPtr(fPtr, fa.PtrLevel)
		// 如果字段是初始化值，尝试设置默认值
		if opts.UseDefValue && vOpt.DefValue != "" && isInitialValue(fKind, fPtr) && fa.KVBinder != nil && fPtr != nil {
			fa.KVBinder(fPtr, vOpt.DefValue)
		}

		// Check: 是否需要验证字段数据的合法性
		if opts.UseValid {
			if err = validx.ValidateFieldPtr(fPtr, fKind, vOpt); err != nil {
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
