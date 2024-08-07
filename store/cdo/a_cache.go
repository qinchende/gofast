package cdo

import (
	"sync"
)

var (
	cdoEncPool    = sync.Pool{New: func() any { return &encoder{} }}
	cdoDecPool    = sync.Pool{New: func() any { return &decoder{} }}
	cachedEncMeta sync.Map
	cachedDecMeta sync.Map
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
