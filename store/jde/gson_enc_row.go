package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"runtime/debug"
	"unsafe"
)

// Encoder +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encGsonRowOnlyValues(obj any) (bs []byte, err error) {
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = errors.New(fmt.Sprintf("error code: %d", code))
			} else {
				// 调试的时候打印错误信息
				fmt.Printf("%s\n%s", pic, debug.Stack())
				err = errors.New(fmt.Sprintf("other panic : %s", pic))
			}
		}
	}()

	af := (*rt.AFace)(unsafe.Pointer(&obj))
	var em *encMeta

	// check target object
	if em = cacheGetGsonEncMeta(af.TypePtr); em == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(obj)
		if dstTyp.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		objType := dstTyp.Elem()
		// TODO: 目前只支持 数据源是 struct 类型
		if objType.Kind() != reflect.Struct {
			panic(errValueMustStruct)
		}

		if em = cacheGetEncMeta(objType); em == nil {
			em = newEncodeMeta(objType)
			cacheSetEncMeta(objType, em)
		}
		cacheSetGsonEncMeta(af.TypePtr, em)
	}

	se := subEncode{}
	se.em = em
	se.srcPtr = af.DataPtr

	se.bf = pool.GetBytesNormal()
	se.encGsonRowJustValues()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	pool.FreeBytes(se.bf)

	return
}

// 单条记录序列化保存，只保存值部分，不用保存字段。结果：[v1,v2,v3,...]
// Note：记录必须是 struct 对象
func (se *subEncode) encGsonRowJustValues() {
	tp := *se.bf
	tp = append(tp, '[')

	fls := se.em.ss.FieldsAttr
	flsSize := len(fls)

	// 循环字段
	for fIndex := 0; fIndex < flsSize; fIndex++ {
		ptr := unsafe.Pointer(uintptr(se.srcPtr) + fls[fIndex].Offset)
		ptrCt := fls[fIndex].PtrLevel
		if ptrCt == 0 {
			goto encObjValue
		}

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			tp = append(tp, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

	encObjValue:
		*se.bf = tp
		se.em.fieldsEnc[fIndex](se.bf, ptr, fls[fIndex].Type)
		tp = *se.bf
	}
	if flsSize > 0 {
		tp = tp[:len(tp)-1]
	}

	*se.bf = append(tp, ']')
}

func encGsonRowFromValues(bf *[]byte, values []any) {
	*bf = append(*bf, '[')
	for _, val := range values {
		encAny(bf, unsafe.Pointer(&val), nil)
	}

	tp := *bf
	if len(values) > 0 {
		tp = tp[:len(tp)-1]
	}
	tp = append(tp, `],`...)
	*bf = tp
}
