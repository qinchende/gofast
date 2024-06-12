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

// NOTE：本编码将以小端方式存储数值
// --------------------------------
// 数值：0x11223344
// 地址：低 ----> 高
// 小端：44 33 22 11
// 大端：11 22 33 44
// --------------------------------

// 最大 MaxUint16
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func encU16By5(bf *[]byte, typ uint8, v uint64) {
//	*bf = encU16By5Ret(*bf, typ, v)
//}
//func encU16By5Ret(bs []byte, typ uint8, v uint64) []byte {
//	switch {
//	default:
//		panic(errOutOfRange)
//	case v <= 29:
//		bs = append(bs, typ|(uint8(v)))
//	case v <= math.MaxUint8:
//		bs = append(bs, typ|30, uint8(v))
//	case v <= math.MaxUint16:
//		bs = append(bs, typ|31, byte(v), byte(v>>8))
//	}
//	return bs
//}

func encU16By6(bf *[]byte, typ uint8, v uint64) {
	*bf = encU16By6Ret(*bf, typ, v)
}
func encU16By6Ret(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= 61:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|62, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|63, byte(v), byte(v>>8))
	}
	return bs
}

// 最大 MaxUint24
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encU24By5(bf *[]byte, typ uint8, v uint64) {
	*bf = encU24By5Ret(*bf, typ, v)
}
func encU24By5Ret(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= 28:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|29, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|30, byte(v), byte(v>>8))
	case v <= Max3BUint:
		bs = append(bs, typ|31, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

// 最大 MaxUint32
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encU32By6(bf *[]byte, typ uint8, v uint64) {
	*bf = encU32By6Ret(*bf, typ, v)
}

func encU32By6Ret(bs []byte, typ uint8, v uint64) []byte {
	if v <= Max3BUint {
		return encU32By6RetPart1(bs, typ, v)
	}
	return encU32By6RetPart2(bs, typ, v)
}

func encU32By6RetPart1(bs []byte, typ uint8, v uint64) []byte {
	switch {
	case v <= 59:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|60, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|61, byte(v), byte(v>>8))
	case v <= Max3BUint:
		bs = append(bs, typ|62, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

// Note: This func must mark as go:noinline
//
//go:noinline
func encU32By6RetPart2(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= math.MaxUint32:
		return append(bs, typ|63, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	}
}

// 最大 MaxUint64
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encU64By6(bf *[]byte, typ uint8, v uint64) {
	*bf = encU64By6Ret(*bf, typ, v)
}

func encU64By6Ret(bs []byte, typ uint8, v uint64) []byte {
	if v <= Max3BUint {
		return encU64By6RetPart1(bs, typ, v)
	}
	return encU64By6RetPart2(bs, typ, v)
}

func encU64By6RetPart1(bs []byte, typ uint8, v uint64) []byte {
	switch {
	case v <= 55:
		bs = append(bs, typ|(uint8(v)))
	case v <= math.MaxUint8:
		bs = append(bs, typ|56, uint8(v))
	case v <= math.MaxUint16:
		bs = append(bs, typ|57, byte(v), byte(v>>8))
	case v <= Max3BUint:
		bs = append(bs, typ|58, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

//go:noinline
func encU64By6RetPart2(bs []byte, typ uint8, v uint64) []byte {
	switch {
	default:
		panic(errOutOfRange)
	case v <= math.MaxUint32:
		return append(bs, typ|59, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	case v <= Max5BUint:
		return append(bs, typ|60, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32))
	case v <= Max6BUint:
		return append(bs, typ|61, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40))
	case v <= Max7BUint:
		return append(bs, typ|62, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48))
	case v <= math.MaxUint64:
		return append(bs, typ|63, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ---------- int && uint ----------
func encInt[T constraints.Integer](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *((*T)(ptr))
	if v >= 0 {
		encU64By6(bf, TypePosInt, uint64(v))
	} else {
		encU64By6(bf, TypeNegInt, uint64(-v))
	}
}

func encUint[T constraints.Unsigned](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	encU64By6(bf, TypePosInt, uint64(*((*T)(ptr))))
}

// ---------- float32 ----------
func encF32(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *(*uint32)(ptr)
	*bf = append(*bf, FixF32, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func encF32Ret(bs []byte, ptr unsafe.Pointer) []byte {
	v := *(*uint32)(ptr)
	return append(bs, FixF32, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func encF32ValRet(bs []byte, ptr unsafe.Pointer) []byte {
	v := *(*uint32)(ptr)
	return append(bs, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

// ---------- float64 ----------
func encF64(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	v := *(*uint64)(ptr)
	*bf = append(*bf, FixF64, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
}

func encF64Ret(bs []byte, ptr unsafe.Pointer) []byte {
	v := *(*uint64)(ptr)
	return append(bs, FixF64, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
}

func encF64ValRet(bs []byte, ptr unsafe.Pointer) []byte {
	v := *(*uint64)(ptr)
	return append(bs, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
}

// ---------- nil ----------
func encNil(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	*bf = append(*bf, FixNil)
}

// ---------- bool ----------
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

// ---------- string ----------
func encString(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	str := *((*string)(ptr))
	encU32By6(bf, TypeStr, uint64(len(str)))
	*bf = append(*bf, str...)
}

func encStringsDirect(bf *[]byte, strItems []string) {
	tp := *bf
	for idx, _ := range strItems {
		tp = encU32By6Ret(tp, TypeStr, uint64(len(strItems[idx])))
		tp = append(tp, strItems[idx]...)
	}
	*bf = tp
}

func encStringDirect(bf *[]byte, str string) {
	encU32By6(bf, TypeStr, uint64(len(str)))
	*bf = append(*bf, str...)
}

// ---------- bytes ----------
func encBytes(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	bs := *((*[]byte)(ptr))
	encU32By6(bf, TypeStr, uint64(len(bs)))
	*bf = append(*bf, bs...)
}

// ---------- time ----------
// 时间默认都是按 RFC3339 格式存储并解析
func encTime(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
	tp := *bf
	tp = append(tp, FixDateTime)
	*bf = append(tp, (*time.Time)(ptr).Format(cst.TimeFmtRFC3339)...)
}

// ---------- any value ----------
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
		encF32(bf, ptr, nil)
	case float64, *float64:
		encF64(bf, ptr, nil)

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
