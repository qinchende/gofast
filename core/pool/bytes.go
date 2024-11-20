package pool

import (
	"sync"
)

const (
	bytesSizeMini  = 128        // 0. 128B
	bytesSizeDef   = 1024       // 1. 1KB 默认 1KB -> 8KB 随机大小的内存
	bytesSizeLarge = 1024 * 8   // 2. 8KB
	bytesSizeMax   = 1024 * 512 // 3. 512KB 超过这个就直接丢
)

//var steps = [4]int{bytesSizeMini, bytesSizeDef, bytesSizeLarge, bytesSizeMax}

var (
	bytesPoolMini  = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeMini); return &bs }}  // 小字节序列 < 1K
	bytesPoolDef   = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeDef); return &bs }}   // 普通字节序列 < 8K
	bytesPoolLarge = sync.Pool{New: func() any { bs := make([]byte, 0, bytesSizeLarge); return &bs }} // 超大字节序列 > 8K
)

// TODO：这里需要完善
// 比如用一个滑动窗口算法估算最近的使用大小，优化内存分配方案
var (
	lastSize  int
	totalSize int
	loopTimes int
)

//func GetBytesMini() *[]byte {
//	bf := bytesPoolMini.Get().(*[]byte)
//	*bf = (*bf)[0:0]
//	return bf
//}

//func GetBytesLarge() *[]byte {
//	bf := bytesPoolLarge.Get().(*[]byte)
//	*bf = (*bf)[0:0]
//	return bf
//}

func GetBytes() *[]byte {
	return getBySize(lastSize)
}

func GetBytesMin(minSize int) *[]byte {
	return getBySize(minSize)
}

func getBySize(needSize int) *[]byte {
	var bf *[]byte
	if needSize <= bytesSizeMini {
		bf = bytesPoolMini.Get().(*[]byte)
	} else if needSize <= bytesSizeDef {
		bf = bytesPoolDef.Get().(*[]byte)
	} else {
		bf = bytesPoolLarge.Get().(*[]byte)
	}
	*bf = (*bf)[0:0]
	return bf
}

func FreeBytes(bs *[]byte) {
	lastSize = cap(*bs)
	if lastSize > bytesSizeMax {
		return
	} else if lastSize >= bytesSizeLarge {
		bytesPoolLarge.Put(bs)
	} else if lastSize >= bytesSizeDef {
		bytesPoolDef.Put(bs)
	} else {
		bytesPoolMini.Put(bs)
	}
}

//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 自定义一个 BytesPool 对象，方便管理自定义大小以下的 BytesBuffer
//type BytesPool struct {
//	capability int
//	pool       sync.Pool
//}
//
//func NewBytesPool(capability int) *BytesPool {
//	return &BytesPool{
//		capability: capability,
//		pool: sync.Pool{
//			New: func() any {
//				return new(bytes.Buffer)
//			},
//		},
//	}
//}
//
//func (bp *BytesPool) Get() *bytes.Buffer {
//	buf := bp.pool.Get().(*bytes.Buffer)
//	buf.Reset()
//	return buf
//}
//
//func (bp *BytesPool) Put(bf *bytes.Buffer) {
//	if bf.Cap() < bp.capability {
//		bp.pool.Put(bf)
//	}
//}
