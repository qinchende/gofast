package jde

import (
	"github.com/qinchende/gofast/cst"
	"golang.org/x/exp/constraints"
	"strconv"
	"unsafe"
)

//func encString(ptr unsafe.Pointer) string {
//	return *((*string)(ptr))
//}

func encString(bf []byte, ptr unsafe.Pointer) []byte {
	bf = append(bf, '"')
	bf = append(bf, *((*string)(ptr))...)
	bf = append(bf, "\","...)
	return bf
}

func encBool(bf []byte, ptr unsafe.Pointer) []byte {
	if *((*bool)(ptr)) {
		bf = append(bf, "true,"...)
	} else {
		bf = append(bf, "false,"...)
	}
	return bf
}

func encInt[T constraints.Signed](bf []byte, ptr unsafe.Pointer) []byte {
	bf = append(bf, strconv.FormatInt(int64(*((*T)(ptr))), 10)...)
	bf = append(bf, ',')
	return bf
}

func encUint[T constraints.Unsigned](bf []byte, ptr unsafe.Pointer) []byte {
	bf = append(bf, strconv.FormatUint(uint64(*((*T)(ptr))), 10)...)
	bf = append(bf, ',')
	return bf
}

func encFloat[T constraints.Float](bf []byte, ptr unsafe.Pointer) []byte {
	bf = append(bf, strconv.FormatFloat(float64(*((*T)(ptr))), 'g', -1, 64)...)
	bf = append(bf, ',')
	return bf
}

func encAny(bf []byte, ptr unsafe.Pointer) []byte {
	switch (*((*any)(ptr))).(type) {
	case int:
		return encInt[int](bf, ptr)
	case int8:
		return encInt[int8](bf, ptr)
	case int16:
		return encInt[int16](bf, ptr)
	case int32:
		return encInt[int32](bf, ptr)
	case int64:
		return encInt[int64](bf, ptr)

	case uint:
		return encUint[uint](bf, ptr)
	case uint8:
		return encUint[uint8](bf, ptr)
	case uint16:
		return encUint[uint16](bf, ptr)
	case uint32:
		return encUint[uint32](bf, ptr)
	case uint64:
		return encUint[uint64](bf, ptr)

	case float32:
		return encFloat[float32](bf, ptr)
	case float64:
		return encFloat[float64](bf, ptr)

	case bool:
		return encBool(bf, ptr)
	case string:
		return encString(bf, ptr)

	case cst.KV:

	}
	return bf
}
