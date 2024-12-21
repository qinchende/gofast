package wit

import (
	"golang.org/x/exp/constraints"
	"math"
	"unsafe"
)

// 6个字表示一组KV值，6 * 8 = 48Bytes
// 主要是为了避免string,int,float等基础数据类型值用any字段保存时，暗藏的堆分配
type KVItem struct {
	_   [0]func() // disallow ==
	Key string
	Val any
	ptr unsafe.Pointer
	num uint64
}

func Str(k string, v string) (f KVItem) {
	f.Key = k
	pStr := (*string)(unsafe.Pointer(&f.ptr))
	*pStr = v
	f.Val = pStr
	return
}

func Num[T constraints.Integer](k string, v T) *KVItem {
	f := &KVItem{Key: k, num: uint64(v)}
	f.Val = (*T)(unsafe.Pointer(&f.num))
	return f
}

func Bool(k string, v bool) *KVItem {
	val := uint64(0)
	if v {
		val = 1
	}
	f := &KVItem{Key: k, num: val << 63}
	f.Val = (*bool)(unsafe.Pointer(&f.num))
	return f
}

// Note: 只适用于小端机
func F32(k string, v float32) *KVItem {
	f := &KVItem{Key: k, num: uint64(math.Float32bits(v)) << 32}
	f.Val = (*float32)(unsafe.Pointer(&f.num))
	return f
}

func F64(k string, v float64) *KVItem {
	f := &KVItem{Key: k, num: math.Float64bits(v)}
	f.Val = (*float64)(unsafe.Pointer(&f.num))
	return f
}

func Time(k string, v float64) *KVItem {
	f := &KVItem{Key: k, num: math.Float64bits(v)}
	f.Val = (*float64)(unsafe.Pointer(&f.num))
	return f
}

func Any(k string, v any) *KVItem {
	return &KVItem{Key: k, Val: v}
}
