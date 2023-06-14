package jde

import (
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"unsafe"
)

func (se *subEncode) encStart() {
	if se.em.isList {
		if se.em.isArray {
			if se.em.isPtr {
				se.encListPtr(se.em.arrLen)
			} else {
				se.encList(se.em.arrLen)
			}
		} else {
			sh := (*reflect.SliceHeader)(se.srcPtr)
			se.srcPtr = unsafe.Pointer(sh.Data)

			if se.em.isPtr {
				se.encListPtr(sh.Len)
			} else {
				se.encList(sh.Len)
			}
		}
	} else if se.em.isStruct {
		se.encStruct()
	} else if se.em.isMap {
		if se.em.isSuperKV {
			se.encMapKV()
		} else {
			se.encMapGeneral()
		}
	} else if se.em.isPtr {
		se.encPointer()
	} else {
		se.encBasic()
	}
}

// +++++++++++++++++++++++++++++++++++++++++++
// Basic type value
func (se *subEncode) encBasic() {
	bf := *se.bf
	bf = se.em.itemEnc(bf, se.srcPtr, se.em.itemType)
	bf = bf[:len(bf)-1]
	*se.bf = bf
}

// Pointer type value
func (se *subEncode) encPointer() {
	bf := *se.bf

	ptrCt := se.em.ptrLevel
	ptr := se.srcPtr
peelPtr:
	ptr = *(*unsafe.Pointer)(ptr)
	if ptr == nil {
		bf = append(bf, "null,"...)
		goto finished
	}
	ptrCt--
	if ptrCt > 0 {
		goto peelPtr
	}

	bf = encMixItem(bf, ptr, se.em.itemType)

finished:
	bf = bf[:len(bf)-1]
	*se.bf = bf
}

// List type value
func (se *subEncode) encList(size int) {
	bf := *se.bf
	bf = append(bf, '[')
	for i := 0; i < size; i++ {
		bf = se.em.itemEnc(bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemRawSize)), se.em.itemType)
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, ']')
	*se.bf = bf
}

// List item is ptr
func (se *subEncode) encListPtr(size int) {
	bf := *se.bf
	ptrLevel := se.em.ptrLevel

	bf = append(bf, '[')
	for i := 0; i < size; i++ {
		ptrCt := ptrLevel
		ptr := unsafe.Pointer(uintptr(se.srcPtr) + uintptr(i*se.em.itemRawSize))

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			bf = append(bf, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

		bf = se.em.itemEnc(bf, ptr, se.em.itemType)
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, ']')
	*se.bf = bf
}

// Struct type value
func (se *subEncode) encStruct() {
	bf := *se.bf
	fls := se.em.ss.FieldsAttr
	size := len(fls)

	bf = append(bf, '{')
	for i := 0; i < size; i++ {
		bf = append(bf, '"')
		bf = append(bf, se.em.ss.ColumnName(i)...)
		bf = append(bf, "\":"...)

		ptr := unsafe.Pointer(uintptr(se.srcPtr) + fls[i].Offset)
		ptrCt := fls[i].PtrLevel
		if ptrCt == 0 {
			goto encObjValue
		}

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			bf = append(bf, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

	encObjValue:
		bf = se.em.fieldsEnc[i](bf, ptr, fls[i].Type)
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, '}')
	*se.bf = bf
}

// Use SubEncode to encode Mix Item Value
// +++++++++++++++++++++++++++++++++++++++++++
func encMixItem(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	se := subEncode{}
	se.getEncMeta(typ, ptr)
	se.bf = &bf // Note: 此处产生了切片变量逃逸

	se.encStart()
	*se.bf = append(*se.bf, ',')

	return *se.bf
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encInt[T constraints.Signed](bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bf = append(bf, strconv.FormatInt(int64(*((*T)(ptr))), 10)...)
	return append(bf, ',')
}

func encIntOnly[T constraints.Signed](bf []byte, ptr unsafe.Pointer) []byte {
	return append(bf, strconv.FormatInt(int64(*((*T)(ptr))), 10)...)
}

func encUint[T constraints.Unsigned](bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bf = append(bf, strconv.FormatUint(uint64(*((*T)(ptr))), 10)...)
	return append(bf, ',')
}

func encUintOnly[T constraints.Unsigned](bf []byte, ptr unsafe.Pointer) []byte {
	return append(bf, strconv.FormatUint(uint64(*((*T)(ptr))), 10)...)
}

func encFloat64(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bf = append(bf, strconv.FormatFloat(*((*float64)(ptr)), 'g', -1, 64)...)
	return append(bf, ',')
}

func encFloat32(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bf = append(bf, strconv.FormatFloat(float64(*((*float32)(ptr))), 'g', -1, 32)...)
	return append(bf, ',')
}

func encString(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bf = append(bf, '"')
	bf = append(bf, *((*string)(ptr))...)
	return append(bf, "\","...)
}

func encStringOnly(bf []byte, ptr unsafe.Pointer) []byte {
	return append(bf, *((*string)(ptr))...)
}

func encBool(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	if *((*bool)(ptr)) {
		bf = append(bf, "true,"...)
	} else {
		bf = append(bf, "false,"...)
	}
	return bf
}

func encAny(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	oldPtr := ptr

	//ei := (*rt.AFace)(ptr)
	ptr = (*rt.AFace)(ptr).DataPtr
	if ptr == nil {
		bf = append(bf, "null,"...)
		return bf
	}

	switch (*((*any)(oldPtr))).(type) {

	case int:
		return encInt[int](bf, ptr, nil)
	case int8:
		return encInt[int8](bf, ptr, nil)
	case int16:
		return encInt[int16](bf, ptr, nil)
	case int32:
		return encInt[int32](bf, ptr, nil)
	case int64:
		return encInt[int64](bf, ptr, nil)

	case uint:
		return encUint[uint](bf, ptr, nil)
	case uint8:
		return encUint[uint8](bf, ptr, nil)
	case uint16:
		return encUint[uint16](bf, ptr, nil)
	case uint32:
		return encUint[uint32](bf, ptr, nil)
	case uint64:
		return encUint[uint64](bf, ptr, nil)

	case float32:
		return encFloat32(bf, ptr, nil)
	case float64:
		return encFloat64(bf, ptr, nil)

	case bool:
		return encBool(bf, ptr, nil)
	case string:
		return encString(bf, ptr, nil)

	default:
		return encMixItem(bf, ptr, reflect.TypeOf(*((*any)(oldPtr))))
		//return encMixItem(bf, ptr, ei.TypePtr)
	}
	//return bf
}
