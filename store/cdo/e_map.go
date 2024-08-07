package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

// Note: 这里是投机取巧了。不是所有的Map都通用，只是尽量解决常见 Map 的编码。比如：cst.KV | map[string]any
// 这样避免了复杂的map反射操作，提高性能。

func (se *encoder) encMap() {
	if se.em.isSuperKV {
		se.encMapKV()
		return
	} else if se.em.keyKind == reflect.String && !se.em.isPtr {
		switch se.em.itemMemSize {
		case 1:
			encMapStrAny[rt.TByte1](se.bf, se.srcPtr, se.em.itemType, se.em.itemEnc)
		case 2:
			encMapStrAny[rt.TByte2](se.bf, se.srcPtr, se.em.itemType, se.em.itemEnc)
		case 4:
			encMapStrAny[rt.TByte4](se.bf, se.srcPtr, se.em.itemType, se.em.itemEnc)
		case 8:
			encMapStrAny[rt.TByte8](se.bf, se.srcPtr, se.em.itemType, se.em.itemEnc)
		case 16:
			encMapStrAny[rt.TByte16](se.bf, se.srcPtr, se.em.itemType, se.em.itemEnc)
		case 24:
			encMapStrAny[rt.TByte24](se.bf, se.srcPtr, se.em.itemType, se.em.itemEnc)
		default:
			panic(errMapType)
		}
		return
	}

	ks := se.em.keySize
	vs := se.em.itemMemSize

	switch ks {
	case 1:
		switch vs {
		case 1:
			encMapAnyAny[int8, rt.TByte1](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 2:
			encMapAnyAny[int8, rt.TByte2](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 4:
			encMapAnyAny[int8, rt.TByte4](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 8:
			encMapAnyAny[int8, rt.TByte8](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 16:
			encMapAnyAny[int8, rt.TByte16](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 24:
			encMapAnyAny[int8, rt.TByte24](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		}

	case 2:
		switch vs {
		case 1:
			encMapAnyAny[int16, rt.TByte1](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 2:
			encMapAnyAny[int16, rt.TByte2](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 4:
			encMapAnyAny[int16, rt.TByte4](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 8:
			encMapAnyAny[int16, rt.TByte8](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 16:
			encMapAnyAny[int16, rt.TByte16](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 24:
			encMapAnyAny[int16, rt.TByte24](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		}

	case 4:
		switch vs {
		case 1:
			encMapAnyAny[int32, rt.TByte1](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 2:
			encMapAnyAny[int32, rt.TByte2](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 4:
			encMapAnyAny[int32, rt.TByte4](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 8:
			encMapAnyAny[int32, rt.TByte8](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 16:
			encMapAnyAny[int32, rt.TByte16](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 24:
			encMapAnyAny[int32, rt.TByte24](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		}

	case 8:
		switch vs {
		case 1:
			encMapAnyAny[int64, rt.TByte1](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 2:
			encMapAnyAny[int64, rt.TByte2](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 4:
			encMapAnyAny[int64, rt.TByte4](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 8:
			encMapAnyAny[int64, rt.TByte8](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 16:
			encMapAnyAny[int64, rt.TByte16](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 24:
			encMapAnyAny[int64, rt.TByte24](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		}

	case 16:
		switch vs {
		case 1:
			encMapAnyAny[string, rt.TByte1](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 2:
			encMapAnyAny[string, rt.TByte2](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 4:
			encMapAnyAny[string, rt.TByte4](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 8:
			encMapAnyAny[string, rt.TByte8](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 16:
			encMapAnyAny[string, rt.TByte16](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		case 24:
			encMapAnyAny[string, rt.TByte24](se.bf, se.srcPtr, se.em.ptrLevel, se.em.keyEnc, se.em.itemType, se.em.itemEnc)
		}

	default:
		panic(errMapType)
	}
}

// cst.KV is map[string]any
// 这是最常见的场景，单独拿出来快速处理
func (se *encoder) encMapKV() {
	// planA
	// theMap := *(*map[string]any)(se.srcPtr)
	// planB
	var theMap map[string]any
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = se.srcPtr

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		bs = encStringDirectRet(bs, k)
		bs = encAnyRet(bs, unsafe.Pointer(&v), nil)
	}
	*se.bf = bs
}

func encMapStrAny[TV any](bf *[]byte, ptr unsafe.Pointer, valTyp reflect.Type, valEnc encValFunc) {
	var theMap map[string]TV
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = ptr

	bs := *bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		bs = encStringDirectRet(bs, k)
		bs = valEnc(bs, unsafe.Pointer(&v), valTyp)
	}
	*bf = bs
}

func encMapAnyAny[TK string | constraints.Integer, TV any](bf *[]byte, ptr unsafe.Pointer, ptrCt uint8,
	keyEnc encValFunc, valTyp reflect.Type, valEnc encValFunc) {
	var theMap map[TK]TV
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = ptr

	bs := *bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		// key
		bs = keyEnc(bs, unsafe.Pointer(&k), nil)

		// value
		// --- ptr ---
		ptrLevel := ptrCt
		ptr = unsafe.Pointer(&v)
		if ptrLevel == 0 {
			goto encMapValue
		}

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			bs = append(bs, FixNil)
			continue
		}
		ptrLevel--
		if ptrLevel > 0 {
			goto peelPtr
		}
		// -----------

	encMapValue:
		bs = valEnc(bs, ptr, valTyp)
	}
}
