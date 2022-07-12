package mapx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"reflect"
)

// 只用传入的值赋值对象
func applyKVToStruct(dest any, kvs cst.KV, opts *ApplyOptions) error {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}

	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not like struct", dest)
	}

	sm := SchemaOfType(dstVal.Type())
	var fls []string
	if opts.FieldDirect {
		fls = sm.fields
	} else {
		fls = sm.columns
	}
	flsOpts := sm.fieldsOpts

	for i := 0; i < len(fls); i++ {
		opt := flsOpts[i]
		fk := fls[i]
		sv, ok := kvs[fk]
		// 只要找到相应的Key，不管Value是啥，不能使用默认值，不能转换就抛异常，说明数据源错误
		if ok {
		} else if opts.NotDefValue != true && opt != nil && opt.DefExist {
			sv = opt.Default
		} else {
			continue
		}

		fv := sm.RefValueByIndex(&dstVal, int8(i))

		// 如果字段是结构体类型
		fvType := fv.Type()
		if fvType.Kind() == reflect.Struct && fvType.String() != "time.Time" {
			// 如果sv 无法转换成 cst.KV 类型，将要抛出异常
			err := applyKVToStruct(fv.Addr().Interface(), sv.(map[string]any), opts)
			if err != nil {
				return err
			}
			continue
		}

		err := sdxSetValue(fv, sv, opts)
		errPanic(err)
	}
	return nil
}

// src 只能是 array,slice 类型
func applyList(dst any, src any, opts *ApplyOptions) error {
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
			if err := sdxSetValue(dstV.Index(i), srcV.Index(i).Interface(), opts); err != nil {
				return err
			}
		}
	default:
		return errors.New("only array-like value supported")
	}
	return nil
}
