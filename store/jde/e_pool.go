package jde

import (
	"sync"
)

var (
	jdeEncodeBufPool = sync.Pool{} // 字节序列缓存
	cachedEncMeta    sync.Map      // cached encode object meta info
)

func (se *subEncode) newBytesBuf() {
	if ret := jdeEncodeBufPool.Get(); ret != nil {
		se.bf = ret.(*[]byte)
		*se.bf = (*se.bf)[0:0]
	} else {
		bs := make([]byte, 0, defEncodeBufSize)
		se.bf = &bs
	}
}

func (se *subEncode) freeBytesBuf() {
	jdeEncodeBufPool.Put(se.bf)
	se.bf = nil
}

func cacheSetEncMeta(typAddr any, val *encMeta) {
	cachedEncMeta.Store(typAddr, val)
}

func cacheGetEncMeta(typAddr any) *encMeta {
	if ret, ok := cachedEncMeta.Load(typAddr); ok {
		return ret.(*encMeta)
	}
	return nil
}
