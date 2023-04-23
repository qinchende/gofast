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
	// refType reflect.Type
	// refVal   reflect.Value // 反射值

	arrPtr  uintptr // 数组首值地址
	arrSize int     // item类型对应的内存字节大小
	arrLen  int     // 数组长度
	arrIdx  int     // 数组索引

	// arrIntFunc bindIntFunc
	// arrFloatFunc bindFloatFunc
	// arrStrFunc   bindStrFunc
	// arrBoolFunc  bindBoolFunc

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
func (sd *subDecode) initListStruct(dst any) (err error) {
	rfVal := reflect.ValueOf(dst)
	if rfVal.Kind() != reflect.Pointer {
		return errValueMustPtr
	}
	rfVal = rfVal.Elem()

	sd.getPool() // 需要用到缓冲池

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
		sd.pl.arr.dst = dst

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
	sd.isList = true

	a := &sd.pl.arr
	a.listType = rfVal.Type()
	a.itemType = a.listType.Elem()
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

	sd.arr = a
	return nil
}

func (sd *subDecode) initArrayMeta(rfVal reflect.Value) {
	sd.isArray = true

	a := sd.arr
	a.arrLen = rfVal.Len()
	a.arrIdx = 0

	if a.isPtr {
		return
	}

	a.arrSize = int(a.itemType.Size())
	a.arrPtr = uintptr((*emptyInterface)(unsafe.Pointer(&a.dst)).ptr)
	//a.arrIntFunc = kindIntFunc[a.itemKind]
}
