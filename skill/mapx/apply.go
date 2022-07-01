package mapx

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"reflect"
)

// 只用传入的值赋值对象
func applyKVToStruct(dest any, kvs cst.KV, useName bool, useDef bool) error {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}

	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	if dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("%T is not like struct", dest)
	}

	sm := SchemaOfType(dstVal.Type())
	var fls []string
	if useName {
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
		} else if useDef && opt != nil && opt.DefExist {
			sv = opt.Default
		} else {
			continue
		}

		fv := sm.RefValueByIndex(&dstVal, int8(i))
		err := sdxSetValue(fv, sv, opt, useName, useDef)
		errPanic(err)
	}
	return nil
}

// src 只能是 array,slice 类型
func applyList(dst any, src any, useName bool, useDef bool) error {
	srcV := reflect.Indirect(reflect.ValueOf(src))

	switch srcV.Kind() {
	//case reflect.Slice:
	//
	//case reflect.Array:
	//
	//	vs := []string{src}
	//	if len(vs) != dst.Len() {
	//		return fmt.Errorf("%q is not valid value for %s", vs, dst.Type().String())
	//	}

	default:
		return errNotArrayType
	}
	return nil
}

func checkDestType(dest any) (reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		panic("dest must be pointer.")
	}
	dSliceTyp := dTyp.Elem()
	if dSliceTyp.Kind() != reflect.Slice {
		panic("dest must be slice.")
	}

	isPtr := false
	isKV := false
	dItemType := dSliceTyp.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if dItemType.Kind() == reflect.Ptr {
		isPtr = true
		dItemType = dItemType.Elem()
	} else if dItemType.Name() == "KV" {
		isKV = true
	}

	return dSliceTyp, dItemType, isPtr, isKV
}
