package jde

import (
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 基础数据类型的Array或Slice
type listPost struct {
	dst any // 原始值
	//refType reflect.Type
	//refVal   reflect.Value // 反射值

	arrSize int // item类型对应的内存字节大小
	arrLen  int // 数组长度
	arrIdx  int // 数组索引
	arrPtr  uintptr

	arrIntFunc bindIntFunc
	//arrFloatFunc bindFloatFunc
	//arrStrFunc   bindStrFunc
	//arrBoolFunc  bindBoolFunc

	listType reflect.Type
	itemType reflect.Type
	itemKind reflect.Kind

	isPtr    bool
	ptrLevel uint8
}

// 解析Struct对象当前的meta信息
type structPost struct {
	refVal reflect.Value     // 反射值
	sm     *dts.StructSchema // 目标值是一个Struct时候
}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) initListStruct(rfVal reflect.Value) (err error) {
	rfTyp := rfVal.Type()

	switch kd := rfTyp.Kind(); kd {
	case reflect.Struct:
		if rfTyp.String() == "time.Time" {
			return errValueType
		}
		sd.isStruct = true
		sd.obj = &sd.pl.obj
	case reflect.Array, reflect.Slice:
		if err = sd.initListMeta(rfVal); err != nil {
			return
		}
		sd.isList = true

		// 进一步初始化数组
		if kd == reflect.Array {
			sd.initArrayMeta(rfVal)
		}
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

func (sd *subDecode) initArrayMeta(rfVal reflect.Value) {
	sd.isArray = true
	sd.arr.arrLen = rfVal.Len()
	sd.arr.arrIdx = 0

	if sd.arr.isPtr {
		return
	}

	sd.arr.arrSize = int(sd.arr.itemType.Size())
	sd.arr.arrPtr = uintptr((*emptyInterface)(unsafe.Pointer(&sd.arr.dst)).ptr)
	sd.arr.arrIntFunc = kindIntFunc[sd.arr.itemKind]
}
