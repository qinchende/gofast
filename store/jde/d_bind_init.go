package jde

import (
	"github.com/qinchende/gofast/store/dts"
	"reflect"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 基础数据类型的Array或Slice
type listMeta struct {
	dst any // 原始值
	//refType reflect.Type
	//refVal   reflect.Value // 反射值

	arrSize  int // 数组长度
	memSize  int // item类型对应的内存字节大小
	listType reflect.Type
	itemType reflect.Type
	itemKind reflect.Kind

	isPtr    bool
	ptrLevel uint8
}

// 解析Struct对象当前的meta信息
type structMeta struct {
	refVal reflect.Value     // 反射值
	sm     *dts.StructSchema // 目标值是一个Struct时候
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) initListStruct(rfVal reflect.Value) (err error) {
	rfTyp := rfVal.Type()

	switch rfTyp.Kind() {
	case reflect.Struct:
		if rfTyp.String() == "time.Time" {
			return errValueType
		}
		sd.isStruct = true
		sd.obj = &sd.pl.obj
	case reflect.Array:
		if err = sd.initListMeta(rfVal); err != nil {
			return
		}
		sd.isArray = true
		sd.isList = true
		sd.arr.arrSize = rfVal.Len()
		sd.arr.memSize = int(sd.arr.itemType.Size())
	case reflect.Slice:
		if err = sd.initListMeta(rfVal); err != nil {
			return
		}
		sd.isList = true
	default:
		return errValueType
	}
	return nil
}

func (sd *subDecode) initListMeta(rfVal reflect.Value) error {
	sd.arr = &sd.pl.arr
	//sd.arr.dst = rfVal.Addr().Interface()

	sd.arr.listType = rfVal.Type()
	sd.arr.itemType = sd.arr.listType.Elem()
	sd.arr.itemKind = sd.arr.itemType.Kind()

peelPtr:
	if sd.arr.itemKind == reflect.Pointer {
		sd.arr.itemType = sd.arr.itemType.Elem()
		sd.arr.itemKind = sd.arr.itemType.Kind()
		sd.arr.isPtr = true
		sd.arr.ptrLevel++
		// TODO：指针嵌套不能超过3层
		if sd.arr.ptrLevel > 3 {
			return errPtrLevel
		}
		goto peelPtr
	}
	return nil
}
