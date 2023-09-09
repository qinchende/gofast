package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"runtime/debug"
	"sync"
	"unsafe"
)

var (
	cachedGsonRowEncMeta sync.Map
)

func cacheSetGsonRowEncMeta(typAddr *rt.TypeAgent, val *encMeta) {
	cachedGsonRowEncMeta.Store(typAddr, val)
}

func cacheGetGsonRowEncMeta(typAddr *rt.TypeAgent) *encMeta {
	if ret, ok := cachedGsonRowEncMeta.Load(typAddr); ok {
		return ret.(*encMeta)
	}
	return nil
}

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
	if em = cacheGetGsonRowsEncMeta(af.TypePtr); em == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(obj)
		if dstTyp.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		objType := dstTyp.Elem()
		if objType.Kind() != reflect.Struct {
			panic(errValueMustStruct)
		}

		typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&objType)).DataPtr)
		if em = cacheGetEncMeta(typAddr); em == nil {
			em = newEncodeMeta(objType)
			cacheSetEncMeta(typAddr, em)
		}
		cacheSetGsonRowsEncMeta(af.TypePtr, em)
	}

	//// 检查 pet 参数是否齐全，缺失就补齐
	//if len(pet.FlsStr) == 0 {
	//	pet.FlsStr, pet.FlsIdxes = em.ss.CTips()
	//} else if len(pet.FlsIdxes) == 0 {
	//	pet.FlsIdxes = em.ss.CIndexes(strings.Split(pet.FlsStr, ","))
	//}

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

	//// struct slice
	//sh := (*reflect.SliceHeader)(se.srcPtr)
	//
	//// 0. 当前记录数量
	//tp = append(tp, strconv.FormatInt(int64(sh.Len), 10)...)
	//tp = append(tp, ',')
	//// 1. 总记录数量
	//tp = append(tp, strconv.FormatInt(pet.Tt, 10)...)
	//tp = append(tp, ",["...)

	//// 2. 字段
	//tp = append(tp, pet.FlsStr...)
	//tp = append(tp, "],["...)

	//// 3. 记录值
	fls := se.em.ss.FieldsAttr
	flsSize := len(fls)
	//// 循环记录
	//for i := 0; i < sh.Len; i++ {
	//	se.srcPtr = unsafe.Pointer(sh.Data + uintptr(i*se.em.itemRawSize))
	//
	//	tp = append(tp, '[')
	// 循环字段
	for fIndex := 0; fIndex < flsSize; fIndex++ {
		//idx := pet.FlsIdxes[fIndex]

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
	//tp = append(tp, "],"...)
	//}

	//if sh.Len > 0 {
	//	tp = tp[:len(tp)-1]
	//}
	*se.bf = append(tp, ']')
}
