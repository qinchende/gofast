package dts

import (
	"reflect"
	"unsafe"
)

func GetInt(kd reflect.Kind, p unsafe.Pointer) int64 {
	switch kd {
	default:
		panic(errNumOutOfRange)
	case reflect.Int:
		return int64(*(*int)(p))
	case reflect.Int8:
		return int64(*(*int8)(p))
	case reflect.Int16:
		return int64(*(*int16)(p))
	case reflect.Int32:
		return int64(*(*int32)(p))
	case reflect.Int64:
		return *(*int64)(p)
	}
}

func GetUint(kd reflect.Kind, p unsafe.Pointer) uint64 {
	switch kd {
	default:
		panic(errNumOutOfRange)
	case reflect.Uint:
		return uint64(*(*uint)(p))
	case reflect.Uint8:
		return uint64(*(*uint8)(p))
	case reflect.Uint16:
		return uint64(*(*uint16)(p))
	case reflect.Uint32:
		return uint64(*(*uint32)(p))
	case reflect.Uint64:
		return *(*uint64)(p)
	case reflect.Uintptr:
		return uint64(*(*uintptr)(p))
	}
}
