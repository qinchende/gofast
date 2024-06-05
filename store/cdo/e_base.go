package cdo

import (
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"math"
	"reflect"
	"time"
	"unsafe"
)

func encNumMax2BWith6(bf *[]byte, typ uint8, v uint64) {
	if v <= math.MaxUint8 {
		if v <= 61 {
			*bf = append(*bf, typ|uint8(v))
		} else {
			*bf = append(*bf, typ|62, uint8(v))
		}
	} else if v <= math.MaxUint16 {
		*bf = append(*bf, typ|63, byte(v>>8), byte(v))
	} else {
		panic(errOutOfRange)
	}
}

//go:inline
func encUint16(bf *[]byte, typ uint8, v uint64) {
	if v <= math.MaxUint8 {
		if v <= 29 {
			*bf = append(*bf, typ|(uint8(v)))
		} else {
			*bf = append(*bf, typ|30, uint8(v))
		}
	} else if v <= math.MaxUint16 {
		*bf = append(*bf, typ|31, byte(v>>8), byte(v))
	} else {
		panic(errOutOfRange)
	}
}

//go:inline
func encUint32(bf *[]byte, typ uint8, v uint64) {
	*bf = encUint32Ret(*bf, typ, v)
}

func encUint32Ret(bs []byte, typ uint8, v uint64) []byte {
	if v <= math.MaxUint8 {
		if v <= 27 {
			bs = append(bs, typ|(uint8(v)))
		} else {
			bs = append(bs, typ|28, uint8(v))
		}
	} else if v <= math.MaxUint16 {
		bs = append(bs, typ|29, byte(v>>8), byte(v))
	} else if v <= Max3BytesUint {
		bs = append(bs, typ|30, byte(v>>16), byte(v>>8), byte(v))
	} else if v <= math.MaxUint32 {
		bs = append(bs, typ|31, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	} else {
		panic(errOutOfRange)
	}
	return bs
}

//go:inline
func encUint64(bf *[]byte, typ uint8, v uint64) {
	if v <= math.MaxUint8 {
		if v <= 23 {
			*bf = append(*bf, typ|uint8(v))
		} else {
			*bf = append(*bf, typ|24, uint8(v))
		}
	} else if v <= math.MaxUint16 {
		*bf = append(*bf, typ|25, byte(v>>8), byte(v))
	} else if v <= Max3BytesUint {
		*bf = append(*bf, typ|26, byte(v>>16), byte(v>>8), byte(v))
	} else if v <= math.MaxUint32 {
		*bf = append(*bf, typ|27, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	} else if v <= Max5BytesUint {
		*bf = append(*bf, typ|28, byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	} else if v <= Max6BytesUint {
		*bf = append(*bf, typ|29, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	} else if v <= Max7BytesUint {
		*bf = append(*bf, typ|30, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	} else {
		*bf = append(*bf, typ|31, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//go:inline
func encInt[T constraints.Integer](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := int64(*((*T)(ptr)))
	if v >= 0 {
		encUint64(bf, TypePosInt, uint64(v))
	} else {
		encUint64(bf, TypeNegInt, uint64(-v))
	}
}

//go:inline
func encUint[T constraints.Unsigned](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	encUint64(bf, TypePosInt, uint64(*((*T)(ptr))))
}

//go:inline
func encFloat32(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *(*uint32)(ptr)
	*bf = append(*bf, FixFloat32, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

//go:inline
func encFloat64(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *(*uint64)(ptr)
	*bf = append(*bf, FixFloat64, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

//go:inline
func encBool(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	if *((*bool)(ptr)) {
		*bf = append(*bf, FixTrue)
	} else {
		*bf = append(*bf, FixFalse)
	}
}

//go:inline
func encNil(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, FixNil)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
//go:inline
func encString(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	str := *((*string)(ptr))
	encUint32(bf, TypeBytes, uint64(len(str)))
	*bf = append(*bf, str...)
}

func encStringsDirect(bf *[]byte, strItems []string) {
	tp := *bf
	for idx, _ := range strItems {
		tp = encUint32Ret(tp, TypeBytes, uint64(len(strItems[idx])))
		tp = append(tp, strItems[idx]...)
	}
	*bf = tp
}

//go:inline
func encStringDirect(bf *[]byte, str string) {
	encUint32(bf, TypeBytes, uint64(len(str)))
	*bf = append(*bf, str...)
}

//go:inline
func encBytes(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	bs := *((*[]byte)(ptr))
	encUint32(bf, TypeBytes, uint64(len(bs)))
	*bf = append(*bf, bs...)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 时间默认都是按 RFC3339 格式存储并解析
//
//go:inline
func encTime(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	tp := *bf
	tp = append(tp, FixDateTime)
	*bf = append(tp, (*time.Time)(ptr).Format(cst.TimeFmtRFC3339)...)
}

func encAny(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	oldPtr := ptr

	// ei := (*rt.AFace)(ptr)
	ptr = (*rt.AFace)(ptr).DataPtr
	if ptr == nil {
		*bf = append(*bf, FixNil)
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
		encMixedItem(bf, ptr, reflect.TypeOf(*((*any)(oldPtr))))
		//return encMixedItem(bf, ptr, ei.TypePtr)
	}
	//return bf
}
