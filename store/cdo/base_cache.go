package cdo

import (
	"encoding/binary"
	"github.com/qinchende/gofast/core/rt"
	"math"
	"sync"
)

var (
	jdeDecPool = sync.Pool{New: func() any { return &subDecode{} }}
	//grsDecPool = sync.Pool{New: func() any { return &gsonRowsDecode{} }}
	jdeBufPool = sync.Pool{New: func() any { return &listPool{} }}

	cachedDecMeta     sync.Map
	cachedDecMetaFast sync.Map
	cachedEncMeta     sync.Map
	cachedEncMetaFast sync.Map
)

// Decode
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func cacheSetDecMeta(typAddr any, val *decMeta) {
	cachedDecMeta.Store(typAddr, val)
}

func cacheGetDecMeta(typAddr any) *decMeta {
	if ret, ok := cachedDecMeta.Load(typAddr); ok {
		return ret.(*decMeta)
	}
	return nil
}

func cacheSetDecMetaFast(typAddr *rt.TypeAgent, val *decMeta) {
	cachedDecMetaFast.Store(typAddr, val)
}

func cacheGetDecMetaFast(typAddr *rt.TypeAgent) *decMeta {
	if ret, ok := cachedDecMetaFast.Load(typAddr); ok {
		return ret.(*decMeta)
	}
	return nil
}

// Encode
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func cacheSetEncMeta(typAddr any, val *encMeta) {
	cachedEncMeta.Store(typAddr, val)
}

func cacheGetEncMeta(typAddr any) *encMeta {
	if ret, ok := cachedEncMeta.Load(typAddr); ok {
		return ret.(*encMeta)
	}
	return nil
}

func cacheSetEncMetaFast(typAddr *rt.TypeAgent, val *encMeta) {
	cachedEncMetaFast.Store(typAddr, val)
}

func cacheGetEncMetaFast(typAddr *rt.TypeAgent) *encMeta {
	if ret, ok := cachedEncMetaFast.Load(typAddr); ok {
		return ret.(*encMeta)
	}
	return nil
}

func toBytes(f float64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(f))
	return bytes
}

// 将字节切片转换回浮点数
func fromBytes(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
}
