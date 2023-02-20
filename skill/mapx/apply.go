package mapx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/valid"
	"reflect"
)

// 只用传入的值赋值对象
func applyKVToStruct(dest any, kvs cst.KV, applyOpts *ApplyOptions) error {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}
	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not like struct", dest)
	}
	sm := getSchema(dstVal, applyOpts)

	var fls []string
	if applyOpts.FieldDirect {
		fls = sm.fields
	} else {
		fls = sm.columns
	}
	flsOpts := sm.fieldsOpts

	var err error
	for i := 0; i < len(fls); i++ {
		fOpt := flsOpts[i]
		fName := fls[i]
		sv, ok := kvs[fName]

		if ok {
		} else if fOpt != nil {
			if fOpt.Required && applyOpts.NotValid == false {
				return fmt.Errorf("field %s requied", fName)
			} else if applyOpts.NotDefValue != true {
				sv = fOpt.DefValue
				if sv == "" {
					continue
				}
			}
		} else {
			continue
		}

		fv := sm.RefValueByIndex(&dstVal, int8(i))

		// 如果字段是结构体类型
		fvType := fv.Type()
		if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
			// 如果sv 无法转换成 cst.KV 类型，将要抛出异常
			if err = applyKVToStruct(fv.Addr().Interface(), sv.(map[string]any), applyOpts); err != nil {
				return err
			}
			continue
		}

		if err = sdxSetValue(fv, sv, fOpt, applyOpts); err != nil {
			return err
		}

		// 是否需要验证字段数据的合法性
		if !applyOpts.NotValid && fOpt != nil {
			if err = valid.ValidateField(&fv, fOpt); err != nil {
				return err
			}
		}
	}
	return nil
}

// src 只能是 array,slice 类型
func applyList(dst any, src any, fOpt *valid.FieldOpts, applyOpts *ApplyOptions) error {
	dstV := reflect.Indirect(reflect.ValueOf(dst))
	srcV := reflect.Indirect(reflect.ValueOf(src))

	var err error
	if srcV.Kind() == reflect.String {
		var srcNew []any
		// TODO：这种情况下只支持JSON解析，YAML是不支持的
		if err = jsonx.UnmarshalFromString(&srcNew, src.(string)); err != nil {
			return err
		}
		src = srcNew
		srcV = reflect.Indirect(reflect.ValueOf(src))
	}

	dstType := dstV.Type()
	dstKind := dstV.Kind()
	srcKind := dstV.Kind()

	switch {
	case (dstKind == reflect.Slice || dstKind == reflect.Array) && (srcKind == reflect.Slice || srcKind == reflect.Array):
		if dstKind == reflect.Array && dstV.Len() != srcV.Len() {
			return errors.New("array length not match")
		}

		sliceTyp, itemType, isPtr, _ := checkDestType(dst)
		dstNew := reflect.MakeSlice(sliceTyp, srcV.Len(), srcV.Len())
		dstV.Set(dstNew)
		for i := 0; i < srcV.Len(); i++ {
			fv := dstV.Index(i)
			if isPtr {
				fv.Set(reflect.New(itemType))
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				if err = applyKVToStruct(fv.Addr().Interface(), srcV.Index(i).Interface().(map[string]any), applyOpts); err != nil {
					return err
				}
				continue
			}

			if err = sdxSetValue(fv, srcV.Index(i).Interface(), fOpt, applyOpts); err != nil {
				return err
			}
			// 是否需要验证字段数据的合法性
			if !applyOpts.NotValid && fOpt != nil {
				if err = valid.ValidateField(&fv, fOpt); err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("only array-like value supported")
	}

	// 数组不能为空
	if !applyOpts.NotValid && fOpt != nil && fOpt.Required && dstV.Len() == 0 {
		return fmt.Errorf("list field %s requied", dstType.String())
	}

	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 主要用于给dest加上默认值，然后执行下字段验证
func optimizeStruct(dst any, applyOpts *ApplyOptions) error {
	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not like struct", dst)
	}
	sm := getSchema(dstVal, applyOpts)

	var err error
	for i := 0; i < len(sm.fields); i++ {
		fv := sm.RefValueByIndex(&dstVal, int8(i))

		// 如果字段是结构体类型
		fvType := fv.Type()
		if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
			if err = optimizeStruct(fv.Addr().Interface(), applyOpts); err != nil {
				return err
			}
			continue
		}

		// 如果字段值看上去像变量刚生成后默认初始化值，那么就加载默认信息
		fOpt := sm.fieldsOpts[i]
		if isInitialValue(fv) && applyOpts.NotDefValue == false {
			if fOpt.DefValue == "" {
				continue
			}
			if err = sdxSetValue(fv, fOpt.DefValue, fOpt, applyOpts); err != nil {
				return err
			}
		}
		// 是否需要验证字段数据的合法性
		if applyOpts.NotValid == false {
			if err = valid.ValidateField(&fv, fOpt); err != nil {
				return err
			}
		}
	}
	return nil
}
