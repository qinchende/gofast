package jde

import (
	"github.com/qinchende/gofast/store/dts"
	"reflect"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type listDest struct {
	dst      any // 原始值
	itemLen  int32
	itemSize int32
	//refVal   reflect.Value // 反射值
	listType reflect.Type
	itemType reflect.Type
	itemKind reflect.Kind
	isPtr    bool
	ptrLevel uint8
}

type structDest struct {
	refVal reflect.Value     // 反射值
	sm     *dts.StructSchema // 目标值是一个Struct时候
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) initListStruct(rfVal reflect.Value) error {
	a := &sd.pl.arr
	//a.refVal = rfVal
	a.dst = rfVal.Addr().Interface()
	rfTyp := rfVal.Type()

	switch rfTyp.Kind() {
	case reflect.Struct:
		sd.isStruct = true
		if rfTyp.String() == "time.Time" {
			return errValueType
		}
		sd.obj = &sd.pl.obj
	case reflect.Array:
		sd.isArray = true
		fallthrough
	case reflect.Slice:
		sd.isList = true

		a.listType = rfTyp
		a.itemType = rfTyp.Elem()
		a.itemKind = a.itemType.Kind()

	peelPtr:
		if a.itemKind == reflect.Pointer {
			a.itemType = a.itemType.Elem()
			a.itemKind = a.itemType.Kind()
			a.isPtr = true
			a.ptrLevel++
			// TODO：指针嵌套不能超过3层
			if a.ptrLevel > 3 {
				return errPtrLevel
			}
			goto peelPtr
		}

		a.itemLen = int32(rfVal.Len())
		a.itemSize = int32(a.itemType.Size())

		sd.arr = a
	default:
		return errValueType
	}
	return nil
}
