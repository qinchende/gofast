package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"runtime/debug"
	"strconv"
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

		// 支持2种数据源：
		// A. struct B. cst.KV
		kd := itemType.Kind()
		if kd != reflect.Struct && itemType.String() != "cst.KV" {
			panic(errValueMustStruct)
		}

		if em = cacheGetEncMeta(itemType); em == nil {
			em = newEncodeMeta(itemType)
			cacheSetEncMeta(itemType, em)
		}
		cacheSetGsonEncMeta(af.TypePtr, em)
	}

	// 检查 pet 参数是否齐全，缺失就补齐
	if em.isStruct {
		if len(pet.Cls) == 0 {
			pet.Cls, pet.ClsIdx = em.ss.CTips()
		}
		if len(pet.ClsIdx) == 0 {
			pet.ClsIdx = em.ss.CIndexes(pet.Cls)
		}
	} else {

	}

	se := subEncode{}
	se.em = em
	se.srcPtr = af.DataPtr

	se.newBytesBuf()
	if em.isStruct {
		se.encStructListByPet(pet)
	} else {
		se.encMapListByPet(pet)
	}
	bs = make([]byte, len(*se.bf))
	copy(bs, *se.bf)
	se.freeBytesBuf()

	return
}

func (se *subEncode) encStructListByPet(pet gson.RowsEncPet) {
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
	for i := 0; i < len(pet.Cls); i++ {
		if i != 0 {
			tp = append(tp, ',')
		}
		tp = append(tp, '"')
		tp = append(tp, pet.Cls[i]...)
		tp = append(tp, '"')
	}
	tp = append(tp, "],["...)

	// 3. 记录值
	flsSize := len(pet.ClsIdx)
	fls := se.em.ss.FieldsAttr
	// 循环记录
	for i := 0; i < sh.Len; i++ {
		se.srcPtr = unsafe.Pointer(sh.Data + uintptr(i*se.em.itemRawSize))

		tp = append(tp, '[')
		// 循环字段
		for j := 0; j < flsSize; j++ {
			idx := pet.ClsIdx[j]

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

func (se *subEncode) encMapListByPet(pet gson.RowsEncPet) {
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
	for i := 0; i < len(pet.Cls); i++ {
		if i != 0 {
			tp = append(tp, ',')
		}
		tp = append(tp, '"')
		tp = append(tp, pet.Cls[i]...)
		tp = append(tp, '"')
	}
	tp = append(tp, "],["...)

	// 3. 记录值
	//flsStr := strings.ReplaceAll(pet.ClsStr, `"`, "")
	//keys := strings.Split(flsStr, ",")
	kvList, ok := pet.List.(*[]cst.KV)
	if !ok {
		panic(errValueType)
	}

	for i := 0; i < len(*kvList); i++ {
		kv := (*kvList)[i]

		tp = append(tp, '[')
		for idx := range pet.Cls {
			val := kv[pet.Cls[idx]]

			*se.bf = tp
			encAny(se.bf, unsafe.Pointer(&val), nil)
			tp = *se.bf
		}
		if len(pet.Cls) > 0 {
			tp = tp[:len(tp)-1]
		}
		tp = append(tp, "],"...)
	}

	if sh.Len > 0 {
		tp = tp[:len(tp)-1]
	}
	*se.bf = append(tp, "]]"...)
}
