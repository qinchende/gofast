package jde

import (
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func (se *subEncode) encStart() {
	switch {
	case se.em.isList:
		if se.em.isArray {
			if se.em.isPtr {
				se.encListPtr(se.em.arrLen)
			} else {
				se.encList(se.em.arrLen)
			}
			break
		}

		// slice
		sh := (*rt.SliceHeader)(se.srcPtr)
		se.srcPtr = sh.DataPtr
		if se.em.isPtr {
			se.encListPtr(sh.Len)
		} else {
			se.encList(sh.Len)
		}
	case se.em.isStruct:
		se.encStruct()
	case se.em.isMap:
		se.encMap()
	case se.em.isPtr:
		se.encPointer()
	default:
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
		itPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		// TODO: 一些本身就是引用类型的数据，需要找到他们指向值的地址
		// 比如 map | function | channel 等类型
		if se.em.itemKind == reflect.Map {
			itPtr = *(*unsafe.Pointer)(itPtr)
			se.em.itemEnc(se.bf, itPtr, se.em.itemType)
		} else {
			se.em.itemEnc(se.bf, itPtr, se.em.itemType)
		}
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
		itPtr := unsafe.Pointer(uintptr(se.srcPtr) + uintptr(i*se.em.itemMemSize))

	peelPtr:
		itPtr = *(*unsafe.Pointer)(itPtr)
		if itPtr == nil {
			tp = append(tp, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

		*se.bf = tp
		// add by cd.net on 2023-12-01
		// []*map 类似这种数据源，不要再剥离一次指针，指向map值的内存
		if se.em.itemKind == reflect.Map {
			itPtr = *(*unsafe.Pointer)(itPtr)
		}
		se.em.itemEnc(se.bf, itPtr, se.em.itemType)
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

func encBytes(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	tp := *bf
	tp = append(tp, '"')
	tp = append(tp, *((*[]byte)(ptr))...)
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

// 时间默认都是按 RFC3339 格式存储并解析
func encTime(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	tp := *bf
	tp = append(tp, '"')
	tp = append(tp, (*time.Time)(ptr).Format(cst.TimeFmtRFC3339)...)
	*bf = append(tp, "\","...)
}

func encAny(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	oldPtr := ptr

	// ei := (*rt.AFace)(ptr)
	ptr = (*rt.AFace)(ptr).DataPtr
	if ptr == nil {
		*bf = append(*bf, "null,"...)
		return
	}

	switch (*((*any)(oldPtr))).(type) {

	case int, *int:
		encInt[int](bf, ptr, nil)
	case int8, *int8:
		encInt[int8](bf, ptr, nil)
	case int16, *int16:
		encInt[int16](bf, ptr, nil)
	case int32, *int32:
		encInt[int32](bf, ptr, nil)
	case int64, *int64:
		encInt[int64](bf, ptr, nil)

	case uint, *uint:
		encUint[uint](bf, ptr, nil)
	case uint8, *uint8:
		encUint[uint8](bf, ptr, nil)
	case uint16, *uint16:
		encUint[uint16](bf, ptr, nil)
	case uint32, *uint32:
		encUint[uint32](bf, ptr, nil)
	case uint64, *uint64:
		encUint[uint64](bf, ptr, nil)

	case float32, *float32:
		encFloat32(bf, ptr, nil)
	case float64, *float64:
		encFloat64(bf, ptr, nil)

	case bool, *bool:
		encBool(bf, ptr, nil)
	case string, *string:
		encString(bf, ptr, nil)

	case []byte, *[]byte:
		encBytes(bf, ptr, nil)

	case time.Time, *time.Time:
		encTime(bf, ptr, nil)

	default:
		encMixItem(bf, ptr, reflect.TypeOf(*((*any)(oldPtr))))
		//return encMixItem(bf, ptr, ei.TypePtr)
	}
	//return bf
}
