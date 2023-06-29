package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"strconv"
	"unsafe"
)

func encGsonRows(v any) (bs []byte, err error) {
	if v == nil {
		return nullBytes, nil
	}

	defer func() {
		if pic := recover(); pic != nil {
			if err1, ok := pic.(error); ok {
				err = err1
			} else {
				err = errors.New(fmt.Sprint(pic))
			}
		}
	}()

	se := subEncode{}
	se.getEncMeta(reflect.TypeOf(v), (*rt.AFace)(unsafe.Pointer(&v)).DataPtr)

	se.newBytesBuf()
	se.encStart()
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	se.freeBytesBuf()

	return
}

func encToGsonRowsString() {

}

func (se *subEncode) encListGson(size int) {
	pet := (*gson.RowsPet)(se.srcPtr)

	tp := *se.bf
	tp = append(tp, '[')

	// 0. 当前记录数量
	tp = append(tp, strconv.FormatInt(pet.Ct, 10)...)
	tp = append(tp, ',')
	// 1. 总记录数量
	tp = append(tp, strconv.FormatInt(pet.Tt, 10)...)
	tp = append(tp, ",["...)

	// 2. 字段
	tp = append(tp, "],["...)
	se.em.ss.FieldsAttr

	// 3. 记录值

	//*se.bf = append(*se.bf, '[')
	//
	//*bf = append(*bf, strconv.FormatUint(uint64(*((*T)(ptr))), 10)...)

	for i := 0; i < size; i++ {
		se.em.itemEnc(se.bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemRawSize)), se.em.itemType)
	}
	if size > 0 {
		*se.bf = (*se.bf)[:len(*se.bf)-1]
	}
	//*se.bf = append(*se.bf, ']')

	tp = append(tp, "],["...)
}

func (se *subEncode) encGsonRecode() {
	fls := se.em.ss.FieldsAttr
	size := len(fls)

	tp := *se.bf
	tp = append(tp, '{')
	for i := 0; i < size; i++ {
		tp = append(tp, '"')
		tp = append(tp, se.em.ss.ColumnName(i)...)
		tp = append(tp, "\":"...)

		ptr := unsafe.Pointer(uintptr(se.srcPtr) + fls[i].Offset)
		ptrCt := fls[i].PtrLevel
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
		se.em.fieldsEnc[i](se.bf, ptr, fls[i].Type)
		tp = *se.bf
	}
	if size > 0 {
		tp = tp[:len(tp)-1]
	}
	*se.bf = append(tp, '}')
}
