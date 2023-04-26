package jde

import (
	"github.com/qinchende/gofast/store/dts"
	"reflect"
	"unsafe"
)

// 绑定 array
type arrIntFunc func(*listPost, int64)

//type bindFloatFunc func(*listPost, int64)
//type bindStrFunc func(*listPost, int64)
//type bindBoolFunc func(*listPost, int64)

// 绑定Struct
type bindStrFunc func(val string) (err int)
type bindBolFunc func(val bool) (err int)

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

	// intFunc arrIntFunc
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
	objPtr uintptr // 数组首值地址
	//refVal reflect.Value     // 反射值
	ss *dts.StructSchema // 目标值是一个Struct时候
}

//func (sd *structPost) setStringByIndex(idx int, val string) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetString(val)
//}
//
//func (sd *structPost) setBoolByIndex(idx int, val bool) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetBool(val)
//}
//
//func (sd *structPost) setIntByIndex(idx int, val int64) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetInt(val)
//}
//
//func (sd *structPost) setFloatByIndex(idx int, val float64) {
//	rfVal := sd.ss.RefValueByIndex(&sd.refVal, int8(idx))
//	rfVal.SetFloat(val)
//}

// 如果不是map和*GsonRow，只能是 Array|Slice|Struct
func (sd *subDecode) initListStruct(dst any) (err error) {
	rfVal := reflect.ValueOf(dst)
	if rfVal.Kind() != reflect.Pointer {
		return errValueMustPtr
	}
	rfVal = reflect.Indirect(rfVal)

	//sd.getPool() // 需要用到缓冲池

	rfTyp := rfVal.Type()
	switch kd := rfTyp.Kind(); kd {
	case reflect.Struct:
		//if rfTyp.String() == "time.Time" {
		//	return errValueType
		//}
		if err = sd.initStructMeta(rfVal); err != nil {
			return
		}
	case reflect.Array, reflect.Slice:
		if err = sd.initListMeta(rfVal); err != nil {
			return
		}
		sd.arr.dst = dst

		// 进一步初始化数组
		if kd == reflect.Array {
			sd.initArrayMeta(rfVal)
		}
	default:
		return errValueType
	}
	return nil
}

var shareObj = structPost{}

func (sd *subDecode) initStructMeta(rfVal reflect.Value) error {
	sd.isStruct = true
	//o := &sd.pl.obj
	o := &shareObj

	// 先假设用Input模式来解JSON
	//o.refVal = rfVal
	o.ss = dts.SchemaForInputByType(rfVal.Type())
	o.objPtr = rfVal.Addr().Pointer()

	sd.obj = o
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
	// a.intFunc = kindIntFunc[a.itemKind]
}
