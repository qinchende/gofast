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
func encU16By6Ret(bs []byte, sym byte, v uint64) []byte {
	switch {
	default:
		panic(errOutRange)
	case v <= 61:
		bs = append(bs, sym|(byte(v)))
	case v <= MaxUint08:
		bs = append(bs, sym|62, byte(v))
	case v <= MaxUint16:
		bs = append(bs, sym|63, byte(v), byte(v>>8))
	}
	return bs
}

// 最大 MaxUint24
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encU24By5Ret(bs []byte, sym byte, v uint64) []byte {
	switch {
	default:
		panic(errOutRange)
	case v <= 28:
		bs = append(bs, sym|(byte(v)))
	case v <= MaxUint08:
		bs = append(bs, sym|29, byte(v))
	case v <= MaxUint16:
		bs = append(bs, sym|30, byte(v), byte(v>>8))
	case v <= MaxUint24:
		bs = append(bs, sym|31, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

// 最大 MaxUint32
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encU32By6Ret(bs []byte, sym byte, v uint64) []byte {
	if v <= MaxUint24 {
		return encU32By6RetPart1(bs, sym, v)
	}
	return encU32By6RetPart2(bs, sym, v)
}

func encU32By6RetPart1(bs []byte, sym byte, v uint64) []byte {
	switch {
	case v <= 59:
		bs = append(bs, sym|(byte(v)))
	case v <= MaxUint08:
		bs = append(bs, sym|60, byte(v))
	case v <= MaxUint16:
		bs = append(bs, sym|61, byte(v), byte(v>>8))
	case v <= MaxUint24:
		bs = append(bs, sym|62, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

// Note: This func must mark as go:noinline
//
//go:noinline
func encU32By6RetPart2(bs []byte, sym byte, v uint64) []byte {
	switch {
	default:
		panic(errOutRange)
	case v <= MaxUint32:
		return append(bs, sym|63, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	}
}

// 最大 MaxUint64
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encU64By6Ret(bs []byte, sym byte, v uint64) []byte {
	if v <= MaxUint24 {
		return encU64By6RetPart1(bs, sym, v)
	}
	return encU64By6RetPart2(bs, sym, v)
}

func encU64By6RetPart1(bs []byte, sym byte, v uint64) []byte {
	switch {
	case v <= 55:
		bs = append(bs, sym|(byte(v)))
	case v <= MaxUint08:
		bs = append(bs, sym|56, byte(v))
	case v <= MaxUint16:
		bs = append(bs, sym|57, byte(v), byte(v>>8))
	case v <= MaxUint24:
		bs = append(bs, sym|58, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

//go:noinline
func encU64By6RetPart2(bs []byte, sym byte, v uint64) []byte {
	switch {
	case v <= MaxUint32:
		bs = append(bs, sym|59, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	case v <= MaxUint40:
		bs = append(bs, sym|60, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32))
	case v <= MaxUint48:
		bs = append(bs, sym|61, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40))
	case v <= MaxUint56:
		bs = append(bs, sym|62, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48))
	case v <= MaxUint64:
		bs = append(bs, sym|63, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
	}
	return bs
}

// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// Note：特例，此时是按照大端法存储的数据
func encListVarIntPart1(bs []byte, sym byte, v uint64) []byte {
	switch {
	case v <= MaxUint05:
		bs = append(bs, sym|0x00|(byte(v)))
	case v <= MaxUint13:
		bs = append(bs, sym|0x20|byte(v>>8), byte(v))
	case v <= MaxUint21:
		bs = append(bs, sym|0x40|byte(v>>16), byte(v>>8), byte(v))
	case v <= MaxUint24:
		bs = append(bs, sym|0x63, byte(v), byte(v>>8), byte(v>>16))
	}
	return bs
}

//go:noinline
func encListVarIntPart2(bs []byte, sym byte, v uint64) []byte {
	switch {
	case v <= MaxUint32:
		bs = append(bs, sym|0x64, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	case v <= MaxUint40:
		bs = append(bs, sym|0x65, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32))
	case v <= MaxUint48:
		bs = append(bs, sym|0x66, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40))
	case v <= MaxUint56:
		bs = append(bs, sym|0x67, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48))
	case v <= MaxUint64:
		bs = append(bs, sym|0x68, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
	}
	return bs
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ---------- int && uint ----------
func encIntRet[T constraints.Integer](bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	v := *((*T)(ptr))
	if v >= 0 {
		return encU64By6Ret(bf, TypeVarIntPos, uint64(v))
	} else {
		return encU64By6Ret(bf, TypeVarIntNeg, uint64(-v))
	}
}

func encUintRet[T constraints.Unsigned](bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	return encU64By6Ret(bf, TypeVarIntPos, uint64(*((*T)(ptr))))
}

// ---------- float32 ----------
func encF32Ret(bs []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	v := *(*uint32)(ptr)
	return append(bs, FixF32, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func encF32ValRet(bs []byte, f float32) []byte {
	v := math.Float32bits(f)
	return append(bs, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

// ---------- float64 ----------
func encF64Ret(bs []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	v := *(*uint64)(ptr)
	return append(bs, FixF64, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
}

func encF64ValRet(bs []byte, f float64) []byte {
	v := math.Float64bits(f)
	return append(bs, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
}

// ---------- nil ----------

// ---------- bool ----------
func encBoolRet(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	if *((*bool)(ptr)) {
		return append(bf, FixTrue)
	} else {
		return append(bf, FixFalse)
	}
}

// ---------- string ----------
func encStringRet(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	str := *((*string)(ptr))
	bf = encU32By6Ret(bf, TypeStr, uint64(len(str)))
	return append(bf, str...)
}

func encStringsDirectRet(bs []byte, strItems []string) []byte {
	for idx, _ := range strItems {
		v := uint64(len(strItems[idx]))
		if v <= MaxUint24 {
			bs = encU32By6RetPart1(bs, TypeStr, v)
		} else {
			bs = encU32By6RetPart2(bs, TypeStr, v)
		}
		bs = append(bs, strItems[idx]...)
	}
	return bs
}

func encStringDirectRet(bs []byte, str string) []byte {
	bs = encU32By6Ret(bs, TypeStr, uint64(len(str)))
	return append(bs, str...)
}

// ---------- bytes ----------

func encBytesRet(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bs := *((*[]byte)(ptr))
	bf = encU32By6Ret(bf, TypeStr, uint64(len(bs)))
	return append(bf, bs...)
}

// ---------- time ----------
// 时间默认都是按 RFC3339 格式存储并解析
func encTimeRet(bf []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	bf = append(bf, FixTime)
	return append(bf, (*time.Time)(ptr).Format(cst.TimeFmtRFC3339)...)
}

// ---------- any value ----------
func encAnyRet(bs []byte, ptr unsafe.Pointer, typ reflect.Type) []byte {
	oldPtr := ptr

	// ei := (*rt.AFace)(ptr)
	ptr = (*rt.AFace)(ptr).DataPtr
	if ptr == nil {
		return append(bs, FixNil)
	}

	switch (*((*any)(oldPtr))).(type) {

	case int, *int:
		return encIntRet[int](bs, ptr, nil)
	case int8, *int8:
		return encIntRet[int8](bs, ptr, nil)
	case int16, *int16:
		return encIntRet[int16](bs, ptr, nil)
	case int32, *int32:
		return encIntRet[int32](bs, ptr, nil)
	case int64, *int64:
		return encIntRet[int64](bs, ptr, nil)

	case uint, *uint:
		return encUintRet[uint](bs, ptr, nil)
	case uint8, *uint8:
		return encUintRet[uint8](bs, ptr, nil)
	case uint16, *uint16:
		return encUintRet[uint16](bs, ptr, nil)
	case uint32, *uint32:
		return encUintRet[uint32](bs, ptr, nil)
	case uint64, *uint64:
		return encUintRet[uint64](bs, ptr, nil)

	case float32, *float32:
		return encF32Ret(bs, ptr, nil)
	case float64, *float64:
		return encF64Ret(bs, ptr, nil)

	case bool, *bool:
		return encBoolRet(bs, ptr, nil)
	case string, *string:
		return encStringRet(bs, ptr, nil)

	case []byte, *[]byte:
		return encBytesRet(bs, ptr, nil)

	case time.Time, *time.Time:
		return encTimeRet(bs, ptr, nil)

	default:
		return encMixedItemRet(bs, ptr, reflect.TypeOf(*((*any)(oldPtr)))) // ei.TypePtr
	}
}
