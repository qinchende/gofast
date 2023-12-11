package pool

import (
	"bytes"
	"sync"
)

const (
	bytesSizeMini   = 128
	bytesSizeNormal = 1024
	bytesSizeLarge  = 1024 * 8
	bytesSizeMax    = 1024 * 1024 // 最大1M缓存
)

var (
	bytesPoolMini   = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeMini); return &bs }}   // 小字节序列 < 1K
	bytesPoolNormal = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeNormal); return &bs }} // 普通字节序列 < 8K
	bytesPoolLarge  = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeLarge); return &bs }}  // 超大字节序列 > 8K
)

func GetBytesMini() *[]byte {
	bf := bytesPoolMini.Get().(*[]byte)
	*bf = (*bf)[0:0]
	return bf
}

func GetBytesNormal() *[]byte {
	bf := bytesPoolNormal.Get().(*[]byte)
	*bf = (*bf)[0:0]
	return bf
}

func GetBytesLarge() *[]byte {
	bf := bytesPoolLarge.Get().(*[]byte)
	*bf = (*bf)[0:0]
	return bf
}

func FreeBytes(bs *[]byte) {
	if len(*bs) > bytesSizeMax {
		return
	} else if len(*bs) > bytesSizeLarge {
		bytesPoolLarge.Put(bs)
	} else if len(*bs) > bytesSizeNormal {
		bytesPoolNormal.Put(bs)
	} else {
		bytesPoolMini.Put(bs)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自定义一个 BytesPool 对象，方便管理自定义大小以下的 BytesBuffer
type BytesPool struct {
	capability int
	pool       *sync.Pool
}

func NewBytesPool(capability int) *BytesPool {
	return &BytesPool{
		capability: capability,
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (bp *BytesPool) Get() *bytes.Buffer {
	buf := bp.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (bp *BytesPool) Put(bf *bytes.Buffer) {
	if bf.Cap() < bp.capability {
		bp.pool.Put(bf)
	}
}
