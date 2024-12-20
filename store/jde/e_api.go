package jde

import (
	"golang.org/x/exp/constraints"
	"strconv"
	"time"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// NOTE: 一些方法可以供外部使用 Build JSON String
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ++ string
func AppendStrNoQuotes(bs []byte, v string) []byte {
	return addStrNoQuotes(bs, v)
}

func AppendStr(bs []byte, v string) []byte {
	return append(addStrNoQuotes(append(bs, '"'), v), '"')
}

func AppendKey(bs []byte, k string) []byte {
	return addStrQuotes(bs, k, ':')
}

func AppendStrField(bs []byte, k, v string) []byte {
	return addStrQuotes(addStrQuotes(bs, k, ':'), v, ',')
}

func AppendStrListField(bs []byte, k string, list []string) []byte {
	bs = AppendKey(bs, k)

	if len(list) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range list {
		bs = addStrQuotes(bs, list[idx], ',')
	}
	return append(bs[:len(bs)-1], "],"...)
}

func AppendJsonField(bs []byte, k string, v []byte) []byte {
	return append(append(addStrQuotes(bs, k, ':'), v...), ',')
}

// ++ int
func AppendIntField[T constraints.Signed](bs []byte, k string, v T) []byte {
	return append(strconv.AppendInt(AppendKey(bs, k), int64(v), 10), ',')
}

func AppendIntListField[T constraints.Signed](bs []byte, k string, v []T) []byte {
	bs = AppendKey(bs, k)
	if len(v) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range v {
		bs = append(strconv.AppendInt(bs, int64(v[idx]), 10), ',')
	}
	return append(bs[:len(bs)-1], "],"...)
}

// ++ uint
func AppendUintField[T constraints.Unsigned](bs []byte, k string, v T) []byte {
	return append(strconv.AppendUint(AppendKey(bs, k), uint64(v), 10), ',')
}

func AppendUintListField[T constraints.Unsigned](bs []byte, k string, v []T) []byte {
	bs = AppendKey(bs, k)
	if len(v) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range v {
		bs = append(strconv.AppendUint(bs, uint64(v[idx]), 10), ',')
	}
	return append(bs[:len(bs)-1], "],"...)
}

// ++ float
func AppendF32Field(bs []byte, k string, v float32) []byte {
	return append(strconv.AppendFloat(AppendKey(bs, k), float64(v), 'g', -1, 32), ',')
}

func AppendF32sField(bs []byte, k string, v []float32) []byte {
	bs = AppendKey(bs, k)
	if len(v) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range v {
		bs = append(strconv.AppendFloat(bs, float64(v[idx]), 'g', -1, 32), ',')
	}
	return append(bs[:len(bs)-1], "],"...)
}

func AppendF64Field(bs []byte, k string, v float64) []byte {
	return append(strconv.AppendFloat(AppendKey(bs, k), v, 'g', -1, 64), ',')
}

func AppendF64sField(bs []byte, k string, v []float64) []byte {
	bs = AppendKey(bs, k)
	if len(v) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range v {
		bs = append(strconv.AppendFloat(bs, float64(v[idx]), 'g', -1, 64), ',')
	}
	return append(bs[:len(bs)-1], "],"...)
}

// ++ bool
func AppendBoolField(bs []byte, k string, v bool) []byte {
	bs = AppendKey(bs, k)
	if v {
		return append(bs, "true,"...)
	} else {
		return append(bs, "false,"...)
	}
}

func AppendBoolListField(bs []byte, k string, v []bool) []byte {
	bs = AppendKey(bs, k)
	if len(v) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range v {
		if v[idx] {
			bs = append(bs, "true,"...)
		} else {
			bs = append(bs, "false,"...)
		}
	}
	return append(bs[:len(bs)-1], "],"...)
}

// ++ time.Time
func AppendTimeField(bs []byte, k string, v time.Time, fmt string) []byte {
	bs = append(addStrNoQuotes(append(bs, '"'), k), "\":\""...)
	return append(v.AppendFormat(bs, fmt), "\","...)
}

func AppendTimeListField(bs []byte, k string, list []time.Time, fmt string) []byte {
	bs = AppendKey(bs, k)
	if len(list) == 0 {
		return append(bs, "[],"...)
	}

	bs = append(bs, '[')
	for idx := range list {
		bs = append(list[idx].AppendFormat(append(bs, '"'), fmt), "\","...)
	}
	return append(bs[:len(bs)-1], "],"...)
}
