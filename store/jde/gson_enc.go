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
	cachedGsonRowsEncMeta sync.Map
)

func cacheSetGsonRowsEncMeta(typAddr *rt.TypeAgent, val *encMeta) {
	cachedGsonRowsEncMeta.Store(typAddr, val)
}

func cacheGetGsonRowsEncMeta(typAddr *rt.TypeAgent) *encMeta {
	if ret, ok := cachedGsonRowsEncMeta.Load(typAddr); ok {
		return ret.(*encMeta)
	}
	return nil
}

// Encoder +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encToGsonRowsString(pet gson.RowsEncPet) (bs []byte, err error) {
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

	af := (*rt.AFace)(unsafe.Pointer(&pet.Target))
	var em *encMeta

	// check target object
	if em = cacheGetGsonRowsEncMeta(af.TypePtr); em == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(pet.Target)
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

		typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&itemType)).DataPtr)
		if em = cacheGetEncMeta(typAddr); em == nil {
			em = newEncodeMeta(itemType)
			cacheSetEncMeta(typAddr, em)
		}
		cacheSetEncMeta(af.TypePtr, em)
	}

	// 检查 pet 参数是否齐全
	if pet.Fields == "" {
		pet.Fields, pet.FlsIdxes = em.ss.CTips()
	} else if len(pet.FlsIdxes) == 0 {
		pet.FlsIdxes = em.ss.CIndexes(strings.Split(pet.Fields, ","))
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
	tp = append(tp, pet.Fields...)
	tp = append(tp, "],["...)
	*se.bf = tp

	// 3. 记录值
	flsSize := len(pet.FlsIdxes)
	fls := se.em.ss.FieldsAttr
	// 循环记录
	for i := 0; i < sh.Len; i++ {
		se.srcPtr = unsafe.Pointer(sh.Data + uintptr(i*se.em.itemRawSize))

		//tp := *se.bf
		*se.bf = append(*se.bf, '[')
		// 循环字段
		for j := 0; j < flsSize; j++ {
			idx := pet.FlsIdxes[j]

			ptr := unsafe.Pointer(uintptr(se.srcPtr) + fls[idx].Offset)
			ptrCt := fls[idx].PtrLevel
			if ptrCt == 0 {
				goto encObjValue
			}

		peelPtr:
			ptr = *(*unsafe.Pointer)(ptr)
			if ptr == nil {
				*se.bf = append(*se.bf, "null,"...)
				continue
			}
			ptrCt--
			if ptrCt > 0 {
				goto peelPtr
			}

		encObjValue:
			//*se.bf = tp
			se.em.fieldsEnc[idx](se.bf, ptr, fls[idx].Type)
			//tp = *se.bf
		}
		if flsSize > 0 {
			*se.bf = (*se.bf)[:len(*se.bf)-1]
		}
		*se.bf = append(*se.bf, "],"...)
	}
	if sh.Len > 0 {
		*se.bf = (*se.bf)[:len(*se.bf)-1]
	}
	*se.bf = append(*se.bf, ']')
}
