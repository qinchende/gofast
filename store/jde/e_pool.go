package jde

import (
	"sync"
)

var (
	//jdeEncPool    = sync.Pool{New: func() any { return &subEncode{} }}
	jdeBytesPool  = sync.Pool{}
	cachedEncMeta sync.Map // cached dest value meta info
)

func newBytes() *[]byte {
	if ret := jdeBytesPool.Get(); ret != nil {
		bs := ret.(*[]byte)
		*bs = (*bs)[:0]
		return bs
	} else {
		bs := make([]byte, 0, 8*1024)
		return &bs
	}
}

func cacheSetEncMeta(typ *dataType, val *encMeta) {
	cachedEncMeta.Store(typ, val)
}

func cacheGetEncMeta(typ *dataType) *encMeta {
	if ret, ok := cachedEncMeta.Load(typ); ok {
		return ret.(*encMeta)
	}
	return nil
}

//// TODO: buffer pool 需要有个机制，释放那些某次偶发申请太大的buffer，而导致长时间不释放的问题
//type bytesPool struct {
//	buf []byte
//}
