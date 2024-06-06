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

// 最长 MaxUint16 的整数编码
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encUint16L2(bf *[]byte, typ uint8, v uint64) {
	*bf = encUint16L2Ret(*bf, typ, v)
}
func encUint16L2Ret(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= 61:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|62, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|63, byte(v>>8), byte(v))
	}
	return bs
}

func encUint16(bf *[]byte, typ uint8, v uint64) {
	*bf = encUint16Ret(*bf, typ, v)
}
func encUint16Ret(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= 29:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|30, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|31, byte(v>>8), byte(v))
	}
	return bs
}

// 最长 MaxUint32 的整数编码
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encUint32(bf *[]byte, typ uint8, v uint64) {
	*bf = encUint32Ret(*bf, typ, v)
}

func encUint32Ret(bs []byte, typ uint8, v uint64) []byte {
	if v <= Max3BytesUint {
		return encUint32RetPart1(bs, typ, v)
	}
	return encUint32RetPart2(bs, typ, v)
}

func encUint32RetPart1(bs []byte, typ uint8, v uint64) []byte {
	switch {
	case v <= 27:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|28, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|29, byte(v>>8), byte(v))
	case v <= Max3BytesUint:
		bs = append(bs, typ|30, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	return bs
}

// Note: This func must mark as go:noinline
//
//go:noinline
func encUint32RetPart2(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= math.MaxUint32:
		return append(bs, typ|31, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

// 最长 MaxUint64 的整数编码
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encUint64(bf *[]byte, typ uint8, v uint64) {
	*bf = encUint64Ret(*bf, typ, v)
}

func encUint64Ret(bs []byte, typ uint8, v uint64) []byte {
	if v <= Max3BytesUint {
		return encUint64RetPart1(bs, typ, v)
	}
	return encUint64RetPart2(bs, typ, v)
}

func encUint64RetPart1(bs []byte, typ uint8, v uint64) []byte {
	switch {
	case v <= 23:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|24, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|25, byte(v>>8), byte(v))
	case v <= Max3BytesUint:
		bs = append(bs, typ|26, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	return bs
}

//go:noinline
func encUint64RetPart2(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= math.MaxUint32:
		return append(bs, typ|27, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	case v <= Max5BytesUint:
		return append(bs, typ|28, byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	case v <= Max6BytesUint:
		return append(bs, typ|29, byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	case v <= Max7BytesUint:
		return append(bs, typ|30, byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	case v <= math.MaxUint64:
		return append(bs, typ|31, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encInt[T constraints.Integer](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *((*T)(ptr))
	if v >= 0 {
		encUint64(bf, TypePosInt, uint64(v))
	} else {
		encUint64(bf, TypeNegInt, uint64(-v))
	}
}

func encUint[T constraints.Unsigned](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	encUint64(bf, TypePosInt, uint64(*((*T)(ptr))))
}

func encFloat32(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *(*uint32)(ptr)
	*bf = append(*bf, FixFloat32, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func encFloat32Ret(bs []byte, ptr unsafe.Pointer) []byte {
	v := *(*uint32)(ptr)
	return append(bs, FixFloat32, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func encFloat64(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *(*uint64)(ptr)
	*bf = append(*bf, FixFloat64, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func encFloat64Ret(bs []byte, ptr unsafe.Pointer) []byte {
	v := *(*uint64)(ptr)
	return append(bs, FixFloat64, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encBool(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	if *((*bool)(ptr)) {
		*bf = append(*bf, FixTrue)
	} else {
		*bf = append(*bf, FixFalse)
	}
}

func encBoolRet(bs []byte, ptr unsafe.Pointer) []byte {
	if *((*bool)(ptr)) {
		return append(bs, FixTrue)
	} else {
		return append(bs, FixFalse)
	}
}

func encNil(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, FixNil)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

func encStringDirect(bf *[]byte, str string) {
	encUint32(bf, TypeBytes, uint64(len(str)))
	*bf = append(*bf, str...)
}

func encBytes(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	bs := *((*[]byte)(ptr))
	encUint32(bf, TypeBytes, uint64(len(bs)))
	*bf = append(*bf, bs...)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 时间默认都是按 RFC3339 格式存储并解析
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
