package mapx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/mapx/valid"
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

	sm := SchemaOfType(dstVal.Type(), applyOpts)
	var fls []string
	if applyOpts.FieldDirect {
		fls = sm.fields
	} else {
		fls = sm.columns
	}
	flsOpts := sm.fieldsOpts

	var err error
	for i := 0; i < len(fls); i++ {
		opt := flsOpts[i]
		fName := fls[i]
		sv, ok := kvs[fName]

		if ok {
		} else if opt != nil {
			if opt.Required {
				return fmt.Errorf("field %s requied", fName)
			} else if applyOpts.NotDefValue != true {
				sv = opt.DefValue
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

		if err = sdxSetValue(fv, sv); err != nil {
			return err
		}

		// 是否需要验证字段数据的合法性
		if !applyOpts.NotValid && opt != nil {
			if err = valid.ValidateField(fv, opt); err != nil {
				return err
			}
		}
	}
	return nil
}

// src 只能是 array,slice 类型
func applyList(dst any, src any) error {
	dstV := reflect.Indirect(reflect.ValueOf(dst))
	srcV := reflect.Indirect(reflect.ValueOf(src))

	if srcV.Kind() == reflect.String {
		var srcNew []any
		// TODO：这种情况下只支持JSON解析，YAML是不支持的
		if err := jsonx.UnmarshalFromString(&srcNew, src.(string)); err != nil {
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
		} else {
			dstNew := reflect.MakeSlice(dstType, srcV.Len(), srcV.Len())
			dstV.Set(dstNew)
		}
		for i := 0; i < srcV.Len(); i++ {
			if err := sdxSetValue(dstV.Index(i), srcV.Index(i).Interface()); err != nil {
				return err
			}
		}
	default:
		return errors.New("only array-like value supported")
	}
	return nil
}
