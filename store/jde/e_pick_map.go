package jde

import (
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

func (se *subEncode) encMap() {
	if se.em.isSuperKV {
		se.encMapKV()
		return
	} else if se.em.keyKind == reflect.String && !se.em.isPtr {
		switch se.em.itemRawSize {
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
	vs := se.em.itemRawSize

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
func (se *subEncode) encMapKV() {
	tp := *se.bf

	tp = append(tp, '{')
	theMap := *(*map[string]any)(se.srcPtr)
	for k, v := range theMap {
		tp = append(tp, '"')
		tp = append(tp, k...)
		tp = append(tp, "\":"...)

		*se.bf = tp
		encAny(se.bf, unsafe.Pointer(&v), nil)
		tp = *se.bf
	}
	if len(theMap) > 0 {
		tp = tp[:len(tp)-1]
	}
	tp = append(tp, '}')
	*se.bf = tp
}

func encMapStrAny[TV any](bf *[]byte, ptr unsafe.Pointer, valTyp reflect.Type, valEnc encValFunc) {
	tp := *bf

	tp = append(tp, '{')
	theMap := *(*map[string]TV)(ptr)
	for k, v := range theMap {
		// key
		tp = append(tp, '"')
		tp = append(tp, k...)
		tp = append(tp, "\":"...)

		// value
		ptr = unsafe.Pointer(&v)
		*bf = tp
		valEnc(bf, ptr, valTyp)
		tp = *bf
	}
	if len(theMap) > 0 {
		tp = tp[:len(tp)-1]
	}
	tp = append(tp, '}')
	*bf = tp
}

func encMapAnyAny[TK string | constraints.Integer, TV any](bf *[]byte, ptr unsafe.Pointer, ptrCt uint8,
	keyEnc encKeyFunc, valTyp reflect.Type, valEnc encValFunc) {

	*bf = append(*bf, '{')
	theMap := *(*map[TK]TV)(ptr)
	for k, v := range theMap {
		// key
		*bf = append(*bf, '"')
		keyEnc(bf, unsafe.Pointer(&k))
		*bf = append(*bf, "\":"...)

		// value
		ptrLevel := ptrCt
		ptr = unsafe.Pointer(&v)
		if ptrLevel == 0 {
			goto encMapValue
		}

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			*bf = append(*bf, "null,"...)
			continue
		}
		ptrLevel--
		if ptrLevel > 0 {
			goto peelPtr
		}

	encMapValue:
		valEnc(bf, ptr, valTyp)
	}
	if len(theMap) > 0 {
		*bf = (*bf)[:len(*bf)-1]
	}
	*bf = append(*bf, '}')
}
