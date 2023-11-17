package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

var (
	cachedGsonEncMeta sync.Map
)

func cacheSetGsonEncMeta(typAddr *rt.TypeAgent, val *encMeta) {
	cachedGsonEncMeta.Store(typAddr, val)
}

func cacheGetGsonEncMeta(typAddr *rt.TypeAgent) *encMeta {
	if ret, ok := cachedGsonEncMeta.Load(typAddr); ok {
		return ret.(*encMeta)
	}
	return nil
}

// Encoder +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encGsonRows(pet gson.RowsEncPet) (bs []byte, err error) {
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

	af := (*rt.AFace)(unsafe.Pointer(&pet.List))
	var em *encMeta

	// check target object
	if em = cacheGetGsonEncMeta(af.TypePtr); em == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(pet.List)
		if dstTyp.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		sliceType := dstTyp.Elem()
		if sliceType.Kind() != reflect.Slice {
			panic(errValueMustSlice)
		}
		itemType := sliceType.Elem()
		// TODO：只支持struct切片，而不是struct指针切片
		if itemType.Kind() != reflect.Struct {
			panic(errValueMustStruct)
		}

		if em = cacheGetEncMeta(itemType); em == nil {
			em = newEncodeMeta(itemType)
			cacheSetEncMeta(itemType, em)
		}
		cacheSetGsonEncMeta(af.TypePtr, em)
	}

	// 检查 pet 参数是否齐全，缺失就补齐
	if len(pet.FieldsStr) == 0 {
		pet.FieldsStr, pet.FieldsIdx = em.ss.CTips()
	} else if len(pet.FieldsIdx) == 0 {
		pet.FieldsIdx = em.ss.CIndexes(strings.Split(pet.FieldsStr, ","))
	}

	se := subEncode{}
	se.em = em
	se.srcPtr = af.DataPtr

	se.newBytesBuf()
	se.encListGson(pet)
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	se.freeBytesBuf()

	return
}

func (se *subEncode) encListGson(pet gson.RowsEncPet) {
	tp := *se.bf
	tp = append(tp, '[')

	// struct slice
	sh := (*reflect.SliceHeader)(se.srcPtr)

	// 0. 当前记录数量
	tp = append(tp, strconv.FormatInt(int64(sh.Len), 10)...)
	tp = append(tp, ',')
	// 1. 总记录数量
	tp = append(tp, strconv.FormatInt(pet.Tt, 10)...)
	tp = append(tp, ",["...)

	// 2. 字段
	tp = append(tp, pet.FieldsStr...)
	tp = append(tp, "],["...)

	// 3. 记录值
	flsSize := len(pet.FieldsIdx)
	fls := se.em.ss.FieldsAttr
	// 循环记录
	for i := 0; i < sh.Len; i++ {
		se.srcPtr = unsafe.Pointer(sh.Data + uintptr(i*se.em.itemRawSize))

		tp = append(tp, '[')
		// 循环字段
		for j := 0; j < flsSize; j++ {
			idx := pet.FieldsIdx[j]

			ptr := unsafe.Pointer(uintptr(se.srcPtr) + fls[idx].Offset)
			ptrCt := fls[idx].PtrLevel
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
			se.em.fieldsEnc[idx](se.bf, ptr, fls[idx].Type)
			tp = *se.bf
		}
		if flsSize > 0 {
			tp = tp[:len(tp)-1]
		}
		tp = append(tp, "],"...)
	}

	if sh.Len > 0 {
		tp = tp[:len(tp)-1]
	}
	*se.bf = append(tp, "]]"...)
}
