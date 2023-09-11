package jde

import (
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"time"
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
		se.encMap()
	} else if se.em.isPtr {
		se.encPointer()
	} else {
		se.encBasic()
	}
}

// +++++++++++++++++++++++++++++++++++++++++++
// Basic type value
func (se *subEncode) encBasic() {
	se.em.itemEnc(se.bf, se.srcPtr, se.em.itemType)
	*se.bf = (*se.bf)[:len(*se.bf)-1]
}

// Pointer type value
func (se *subEncode) encPointer() {
	ptrCt := se.em.ptrLevel
	ptr := se.srcPtr

peelPtr:
	ptr = *(*unsafe.Pointer)(ptr)
	if ptr == nil {
		*se.bf = append(*se.bf, nullBytes...)
		return
	}
	ptrCt--
	if ptrCt > 0 {
		goto peelPtr
	}

	encMixItem(se.bf, ptr, se.em.itemType)
	*se.bf = (*se.bf)[:len(*se.bf)-1]
}

// List type value
func (se *subEncode) encList(size int) {
	*se.bf = append(*se.bf, '[')
	for i := 0; i < size; i++ {
		se.em.itemEnc(se.bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemRawSize)), se.em.itemType)
	}
	if size > 0 {
		*se.bf = (*se.bf)[:len(*se.bf)-1]
	}
	*se.bf = append(*se.bf, ']')
}

// List item is ptr
func (se *subEncode) encListPtr(size int) {
	ptrLevel := se.em.ptrLevel

	tp := *se.bf
	tp = append(tp, '[')
	for i := 0; i < size; i++ {
		ptrCt := ptrLevel
		ptr := unsafe.Pointer(uintptr(se.srcPtr) + uintptr(i*se.em.itemRawSize))

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

		*se.bf = tp
		se.em.itemEnc(se.bf, ptr, se.em.itemType)
		tp = *se.bf
	}
	if size > 0 {
		tp = tp[:len(tp)-1]
	}
	*se.bf = append(tp, ']')
}

// Struct type value
func (se *subEncode) encStruct() {
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

// Use SubEncode to encode Mix Item Value
// +++++++++++++++++++++++++++++++++++++++++++
func encMixItem(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	se := subEncode{}
	se.getEncMeta(typ, ptr)
	se.bf = bf
	se.encStart()
	*bf = append(*se.bf, ',')
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encInt[T constraints.Signed](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, strconv.FormatInt(int64(*((*T)(ptr))), 10)...)
	*bf = append(*bf, ',')
}

func encIntOnly[T constraints.Signed](bf *[]byte, ptr unsafe.Pointer) {
	*bf = append(*bf, strconv.FormatInt(int64(*((*T)(ptr))), 10)...)
}

func encUint[T constraints.Unsigned](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, strconv.FormatUint(uint64(*((*T)(ptr))), 10)...)
	*bf = append(*bf, ',')
}

func encUintOnly[T constraints.Unsigned](bf *[]byte, ptr unsafe.Pointer) {
	*bf = append(*bf, strconv.FormatUint(uint64(*((*T)(ptr))), 10)...)
}

func encFloat64(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, strconv.FormatFloat(*((*float64)(ptr)), 'g', -1, 64)...)
	*bf = append(*bf, ',')
}

func encFloat32(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, strconv.FormatFloat(float64(*((*float32)(ptr))), 'g', -1, 32)...)
	*bf = append(*bf, ',')
}

func encString(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	tp := *bf
	tp = append(tp, '"')
	tp = append(tp, *((*string)(ptr))...)
	*bf = append(tp, "\","...)
}

func encStringOnly(bf *[]byte, ptr unsafe.Pointer) {
	*bf = append(*bf, *((*string)(ptr))...)
}

func encBool(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	if *((*bool)(ptr)) {
		*bf = append(*bf, "true,"...)
	} else {
		*bf = append(*bf, "false,"...)
	}
}

func encTime(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	tp := *bf
	tp = append(tp, '"')
	tp = append(tp, (*time.Time)(ptr).Format(cst.TimeFmtSaveRFC3339)...)
	*bf = append(tp, "\","...)
}

func encAny(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	oldPtr := ptr

	//ei := (*rt.AFace)(ptr)
	ptr = (*rt.AFace)(ptr).DataPtr
	if ptr == nil {
		*bf = append(*bf, "null,"...)
		return
	}

	switch (*((*any)(oldPtr))).(type) {

	case int:
		encInt[int](bf, ptr, nil)
	case int8:
		encInt[int8](bf, ptr, nil)
	case int16:
		encInt[int16](bf, ptr, nil)
	case int32:
		encInt[int32](bf, ptr, nil)
	case int64:
		encInt[int64](bf, ptr, nil)

	case uint:
		encUint[uint](bf, ptr, nil)
	case uint8:
		encUint[uint8](bf, ptr, nil)
	case uint16:
		encUint[uint16](bf, ptr, nil)
	case uint32:
		encUint[uint32](bf, ptr, nil)
	case uint64:
		encUint[uint64](bf, ptr, nil)

	case float32:
		encFloat32(bf, ptr, nil)
	case float64:
		encFloat64(bf, ptr, nil)

	case bool:
		encBool(bf, ptr, nil)
	case string:
		encString(bf, ptr, nil)

	default:
		encMixItem(bf, ptr, reflect.TypeOf(*((*any)(oldPtr))))
		//return encMixItem(bf, ptr, ei.TypePtr)
	}
	//return bf
}
