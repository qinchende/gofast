package pool

import (
	"bytes"
	"sync"
)

const (
	bytesSizeMini  = 128         // 128B
	bytesSizeDef   = 1024        // 1KB
	bytesSizeLarge = 1024 * 8    // 8KB
	bytesSizeMax   = 1024 * 1024 // 1MB 超过这个就直接丢
)

var (
	bytesPoolMini  = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeMini); return &bs }}  // 小字节序列 < 1K
	bytesPoolDef   = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeDef); return &bs }}   // 普通字节序列 < 8K
	bytesPoolLarge = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeLarge); return &bs }} // 超大字节序列 > 8K
)

func GetBytesMini() *[]byte {
	bf := bytesPoolMini.Get().(*[]byte)
	*bf = (*bf)[0:0]
	return bf
}

func GetBytes() *[]byte {
	bf := bytesPoolDef.Get().(*[]byte)
	*bf = (*bf)[0:0]
	return bf
}

func GetBytesLarge() *[]byte {
	bf := bytesPoolLarge.Get().(*[]byte)
	*bf = (*bf)[0:0]
	return bf
}

func FreeBytes(bs *[]byte) {
	if cap(*bs) > bytesSizeMax {
		return
	} else if cap(*bs) >= bytesSizeLarge {
		bytesPoolLarge.Put(bs)
	} else if cap(*bs) >= bytesSizeDef {
		bytesPoolDef.Put(bs)
	} else {
		bytesPoolMini.Put(bs)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自定义一个 BytesPool 对象，方便管理自定义大小以下的 BytesBuffer
type BytesPool struct {
	capability int
	pool       sync.Pool
}

func NewBytesPool(capability int) *BytesPool {
	return &BytesPool{
		capability: capability,
		pool: sync.Pool{
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
