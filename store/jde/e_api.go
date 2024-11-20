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
	bs = append(bs, '"')
	bs = addStrNoQuotes(bs, v)
	return append(bs, '"')
}

func AppendKey(bs []byte, k string) []byte {
	return addStrQuotes(bs, k, ':')
}

func AppendStrField(bs []byte, k, v string) []byte {
	bs = addStrQuotes(bs, k, ':')
	return addStrQuotes(bs, v, ',')
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
	bs = bs[:len(bs)-1]
	bs = append(bs, "],"...)
	return bs
}

// ++ int
func AppendIntField[T constraints.Signed](bs []byte, k string, v T) []byte {
	bs = AppendKey(bs, k)
	bs = strconv.AppendInt(bs, int64(v), 10)
	bs = append(bs, ',')
	return bs
}

func AppendIntListField[T constraints.Signed](bs []byte, k string, v []T) []byte {
	bs = AppendKey(bs, k)
	bs = appendIntListField[T](bs, v)
	bs = append(bs, ',')
	return bs
}

func appendIntListField[T constraints.Signed](bs []byte, list []T) []byte {
	if len(list) == 0 {
		return append(bs, '[', ']')
	}
	bs = append(bs, '[')
	for idx := range list {
		bs = strconv.AppendInt(bs, int64(list[idx]), 10)
		bs = append(bs, ',')
	}
	bs = bs[:len(bs)-1]
	bs = append(bs, ']')
	return bs
}

// ++ uint
func AppendUintField(bs []byte, k string, v uint64) []byte {
	bs = AppendKey(bs, k)
	bs = strconv.AppendUint(bs, v, 10)
	bs = append(bs, ',')
	return bs
}

func AppendUintListField[T constraints.Unsigned](bs []byte, k string, v []T) []byte {
	bs = AppendKey(bs, k)
	bs = appendUintListField[T](bs, v)
	bs = append(bs, ',')
	return bs
}

func appendUintListField[T constraints.Unsigned](bs []byte, list []T) []byte {
	if len(list) == 0 {
		return append(bs, '[', ']')
	}
	bs = append(bs, '[')
	for idx := range list {
		bs = strconv.AppendUint(bs, uint64(list[idx]), 10)
		bs = append(bs, ',')
	}
	bs = bs[:len(bs)-1]
	bs = append(bs, ']')
	return bs
}

// ++ float
func AppendF32Field(bs []byte, k string, v float32) []byte {
	bs = AppendKey(bs, k)
	bs = strconv.AppendFloat(bs, float64(v), 'g', -1, 32)
	bs = append(bs, ',')
	return bs
}

func AppendF64Field(bs []byte, k string, v float64) []byte {
	bs = AppendKey(bs, k)
	bs = strconv.AppendFloat(bs, v, 'g', -1, 64)
	bs = append(bs, ',')
	return bs
}

func AppendBoolField(bs []byte, k string, v bool) []byte {
	bs = AppendKey(bs, k)
	if v {
		bs = append(bs, "true,"...)
	} else {
		bs = append(bs, "false,"...)
	}
	return bs
}

// ++ time.Time
func AppendTimeField(bs []byte, k string, v time.Time, fmt string) []byte {
	bs = append(bs, '"')
	bs = addStrNoQuotes(bs, k)
	bs = append(bs, "\":\""...)
	bs = v.AppendFormat(bs, fmt)
	bs = append(bs, "\","...)
	return bs
}
