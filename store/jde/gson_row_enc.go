package jde

import (
	"errors"
	"fmt"
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

	se.newBytesBuf()
	se.encGsonRowJustValues()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	se.freeBytesBuf()

	return
}

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
