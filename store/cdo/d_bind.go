package cdo

import (
	"github.com/qinchende/gofast/core/cst"
	"time"
	"unsafe"
)

type intBinder func(p unsafe.Pointer, sym byte, v uint64)

// ++ int ++
func bindInt(p unsafe.Pointer, sym byte, v uint64) {
	*(*int)(p) = toInt(sym, v)
}

func toInt(sym byte, v uint64) int {
	if sym>>7 == 0 {
		if v >= OverInt {
			panic(errInfinity)
		}
		return int(v)
	} else {
		if v > OverInt {
			panic(errInfinity)
		}
		return int(-v)
	}
}

// ++ uint ++
func bindUint(p unsafe.Pointer, sym byte, v uint64) {
	*(*uint)(p) = toUint(sym, v)
}

func toUint(sym byte, v uint64) uint {
	if sym>>7 == 1 || v > MaxUint {
		panic(errInfinity)
	}
	return uint(v)
}

// ++ int8 ++
func bindInt8(p unsafe.Pointer, sym byte, v uint64) {
	*(*int8)(p) = toInt8(sym, v)
}

func toInt8(sym byte, v uint64) int8 {
	if sym>>7 == 0 {
		if v >= OverInt08 {
			panic(errInfinity)
		}
		return int8(v)
	} else {
		if v > OverInt08 {
			panic(errInfinity)
		}
		return int8(-v)
	}
}

// ++ uint8 ++
func bindUint8(p unsafe.Pointer, sym byte, v uint64) {
	*(*uint8)(p) = toUint8(sym, v)
}

func toUint8(sym byte, v uint64) uint8 {
	if sym>>7 == 1 || v > MaxUint08 {
		panic(errInfinity)
	}
	return uint8(v)
}

// ++ int16 ++
func bindInt16(p unsafe.Pointer, sym byte, v uint64) {
	*(*int16)(p) = toInt16(sym, v)
}

func toInt16(sym byte, v uint64) int16 {
	if sym>>7 == 0 {
		if v >= OverInt16 {
			panic(errInfinity)
		}
		return int16(v)
	} else {
		if v > OverInt16 {
			panic(errInfinity)
		}
		return int16(-v)
	}
}

// ++ uint16 ++
func bindUint16(p unsafe.Pointer, sym byte, v uint64) {
	*(*uint16)(p) = toUint16(sym, v)
}

func toUint16(sym byte, v uint64) uint16 {
	if sym>>7 == 1 || v > MaxUint16 {
		panic(errInfinity)
	}
	return uint16(v)
}

// ++ int32 ++
func bindInt32(p unsafe.Pointer, sym byte, v uint64) {
	*(*int32)(p) = toInt32(sym, v)
}

func toInt32(sym byte, v uint64) int32 {
	if sym>>7 == 0 {
		if v >= OverInt32 {
			panic(errInfinity)
		}
		return int32(v)
	} else {
		if v > OverInt32 {
			panic(errInfinity)
		}
		return int32(-v)
	}
}

// ++ uint32 ++
func bindUint32(p unsafe.Pointer, sym byte, v uint64) {
	*(*uint32)(p) = toUint32(sym, v)
}

func toUint32(sym byte, v uint64) uint32 {
	if sym>>7 == 1 || v > MaxUint32 {
		panic(errInfinity)
	}
	return uint32(v)
}

// ++ int64 ++
func bindInt64(p unsafe.Pointer, sym byte, v uint64) {
	*(*int64)(p) = toInt64(sym, v)
}

func toInt64(sym byte, v uint64) int64 {
	if sym>>7 == 0 {
		if v >= OverInt64 {
			panic(errInfinity)
		}
		return int64(v)
	} else {
		if v > OverInt64 {
			panic(errInfinity)
		}
		return int64(-v)
	}
}

// ++ uint64 ++
func bindUint64(p unsafe.Pointer, sym byte, v uint64) {
	*(*uint64)(p) = toUint64(sym, v)
}

func toUint64(sym byte, v uint64) uint64 {
	if sym>>7 == 1 {
		panic(errInfinity)
	}
	return v
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// float ++++++++++++++++++++++++++++++++++++++++++++
func bindF32(p unsafe.Pointer, v float32) {
	*(*float32)(p) = v
}

func bindF64(p unsafe.Pointer, v float64) {
	*(*float64)(p) = v
}

// []byte ++++++++++++++++++++++++++++++++++++++++++++
func bindBytes(p unsafe.Pointer, v []byte) {
	*(*[]byte)(p) = v
}

// string & bool & any +++++++++++++++++++++++++++++++
func bindString(p unsafe.Pointer, v string) {
	*(*string)(p) = v
}

func bindBool(p unsafe.Pointer, v bool) {
	*(*bool)(p) = v
}

func bindAny(p unsafe.Pointer, v any) {
	*(*any)(p) = v
}

// time ++++++++++++++++++++++++++++++++++++++++++++
// 时间默认都是按 RFC3339 格式存储并解析
func bindTime(p unsafe.Pointer, v string) {
	if tm, err := time.Parse(cst.TimeFmtRFC3339, v); err != nil {
		panic(err)
	} else {
		*(*time.Time)(p) = tm
	}
}
