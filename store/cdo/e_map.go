package cdo

import (
	"encoding"
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"unsafe"
)

// Note: 这里投机取巧了。 尽量将常见Map单独处理，提高性能。不常见的用反射

func (se *encoder) encMap() {
	// 特殊处理：map[string]string & map[string]all-type
	if se.em.isMapStrStr {
		se.encMapStrStr()
	} else if se.em.isMapStrAll && !se.em.isPtr {
		switch se.em.itemMemSize {
		default:
			encMapStrAllReflect(se)
		case 1:
			encMapStrAll[rt.TByte1](se)
		case 2:
			encMapStrAll[rt.TByte2](se)
		case 4:
			encMapStrAll[rt.TByte4](se)
		case 8:
			encMapStrAll[rt.TByte8](se)
		case 16:
			encMapStrAll[rt.TByte16](se)
		case 24:
			encMapStrAll[rt.TByte24](se)
		}
	} else {
		se.encMapOthers()
	}
}

func (se *encoder) encMapOthers() {
	kSize := se.em.keySize
	vSize := se.em.itemMemSize

	switch kSize {
	default:
		encMapAllReflect(se)

	case 1:
		switch vSize {
		default:
			encMapAllReflect(se)
		case 1:
			encMapAll[uint8, rt.TByte1](se)
		case 2:
			encMapAll[uint8, rt.TByte2](se)
		case 4:
			encMapAll[uint8, rt.TByte4](se)
		case 8:
			encMapAllPtr[uint8, rt.TByte8](se)
		case 16:
			encMapAll[uint8, rt.TByte16](se)
		case 24:
			encMapAll[uint8, rt.TByte24](se)
		}

	case 2:
		switch vSize {
		default:
			encMapAllReflect(se)
		case 1:
			encMapAll[uint16, rt.TByte1](se)
		case 2:
			encMapAll[uint16, rt.TByte2](se)
		case 4:
			encMapAll[uint16, rt.TByte4](se)
		case 8:
			encMapAllPtr[uint16, rt.TByte8](se)
		case 16:
			encMapAll[uint16, rt.TByte16](se)
		case 24:
			encMapAll[uint16, rt.TByte24](se)
		}

	case 4:
		switch vSize {
		default:
			encMapAllReflect(se)
		case 1:
			encMapAll[uint32, rt.TByte1](se)
		case 2:
			encMapAll[uint32, rt.TByte2](se)
		case 4:
			encMapAll[uint32, rt.TByte4](se)
		case 8:
			encMapAllPtr[uint32, rt.TByte8](se)
		case 16:
			encMapAll[uint32, rt.TByte16](se)
		case 24:
			encMapAll[uint32, rt.TByte24](se)
		}

	case 8:
		switch vSize {
		default:
			encMapAllReflect(se)
		case 1:
			encMapAll[uint64, rt.TByte1](se)
		case 2:
			encMapAll[uint64, rt.TByte2](se)
		case 4:
			encMapAll[uint64, rt.TByte4](se)
		case 8:
			encMapAllPtr[uint64, rt.TByte8](se)
		case 16:
			encMapAll[uint64, rt.TByte16](se)
		case 24:
			encMapAll[uint64, rt.TByte24](se)
		}

	case 16:
		switch vSize {
		default:
			encMapAllReflect(se)
		case 1:
			encMapAll[string, rt.TByte1](se)
		case 2:
			encMapAll[string, rt.TByte2](se)
		case 4:
			encMapAll[string, rt.TByte4](se)
		case 8:
			encMapAllPtr[string, rt.TByte8](se)
		case 16:
			encMapAll[string, rt.TByte16](se)
		case 24:
			encMapAll[string, rt.TByte24](se)
		}
	}
}

// map[string]string
func (se *encoder) encMapStrStr() {
	var theMap map[string]string
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = se.srcPtr

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		bs = encStringDirectRet(bs, k)
		bs = encStringDirectRet(bs, v)
	}
	*se.bf = bs
}

// map[string]all-type
func encMapStrAll[TV any](se *encoder) {
	var theMap map[string]TV
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = se.srcPtr

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		bs = encStringDirectRet(bs, k)
		bs = se.em.itemEnc(bs, unsafe.Pointer(&v), se.em.itemType)
	}
	*se.bf = bs
}

//// map[string]any
//func (se *encoder) encMapStrAny() {
//	var theMap map[string]any
//	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = se.srcPtr
//
//	bs := *se.bf
//	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
//	for k, v := range theMap {
//		bs = encStringDirectRet(bs, k)
//		bs = encAnyRet(bs, unsafe.Pointer(&v), nil)
//	}
//	*se.bf = bs
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// theMap := *(*map[string]any)(se.srcPtr)
// var theMap map[TK]TV
// *(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = ptr
func encMapAll[TK string | constraints.Unsigned, TV any](se *encoder) {
	var theMap map[TK]TV
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = se.srcPtr

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		bs = se.em.keyEnc(bs, unsafe.Pointer(&k), nil)
		bs = se.em.itemEnc(bs, unsafe.Pointer(&v), se.em.itemType)
	}
	*se.bf = bs
}

func encMapAllPtr[TK string | constraints.Unsigned, TV any](se *encoder) {
	var theMap map[TK]TV
	*(*unsafe.Pointer)(unsafe.Pointer(&theMap)) = se.srcPtr

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(theMap))), ListKV)
	for k, v := range theMap {
		// key
		bs = se.em.keyEnc(bs, unsafe.Pointer(&k), nil)

		// value
		// --- if ptr ---
		ptrLevel := se.em.ptrLevel
		ptr := unsafe.Pointer(&v)
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
		bs = se.em.itemEnc(bs, ptr, se.em.itemType)
	}
	*se.bf = bs
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// encode map using reflect
func encMapStrAllReflect(se *encoder) {
	theMap := reflect.MakeMapWithSize(reflect.MapOf(se.em.keyType, se.em.itemType), 0)
	refVal := (*rt.ReflectValue)(unsafe.Pointer(&theMap))
	refVal.DataPtr = se.srcPtr

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(theMap.Len())), ListKV)

	iter := theMap.MapRange()
	for iter.Next() {
		k := resolveKeyName(iter.Key())
		bs = encStringDirectRet(bs, k)

		val := iter.Value()
		vRef := (*rt.ReflectValue)(unsafe.Pointer(&val))
		bs = se.em.itemEnc(bs, vRef.DataPtr, se.em.itemType)
	}
	*se.bf = bs
}

func encMapAllReflect(se *encoder) {
	theMap := reflect.MakeMap(reflect.MapOf(se.em.keyType, se.em.itemType))
	theMap.SetPointer(se.srcPtr)

	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(theMap.Len())), ListKV)
	iter := theMap.MapRange()
	for iter.Next() {
		k := iter.Key().Addr().UnsafePointer()
		v := iter.Value().Addr().UnsafePointer()

		// key
		bs = se.em.keyEnc(bs, k, nil)

		// value
		// --- ptr ---
		ptrLevel := se.em.ptrLevel
		ptr := v
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
		bs = se.em.itemEnc(bs, ptr, se.em.itemType)
	}
	*se.bf = bs
}

// utils +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func resolveKeyName(k reflect.Value) string {
	if k.Kind() == reflect.String {
		return k.String()
	}
	if tm, ok := k.Interface().(encoding.TextMarshaler); ok {
		if k.Kind() == reflect.Pointer && k.IsNil() {
			return ""
		}
		if buf, err := tm.MarshalText(); err == nil {
			return string(buf)
		} else {
			panic(errKey)
		}
	}

	switch k.Kind() {
	default:
		panic(errKey)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(k.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(k.Uint(), 10)
	}
}
