package pool

import (
	"sync"
)

const (
	bsSizeMini  = 128        // 0. 128B
	bsSizeDef   = 1024       // 1. 1KB 默认 1KB -> 8KB 随机大小的内存
	bsSizeLarge = 1024 * 8   // 2. 8KB
	bsSizeHuge  = 1024 * 64  // 3. 64KB
	bsSizeMax   = 1024 * 512 // 4. 512KB 超过这个就直接丢
)

//var steps = [4]int{bytesSizeMini, bytesSizeDef, bytesSizeLarge, bytesSizeMax}

var (
	bsPoolMini  = sync.Pool{New: func() any { bs := make([]byte, 0, bsSizeMini); return &bs }}  // 小字节序列 	< 1K
	bsPoolDef   = sync.Pool{New: func() any { bs := make([]byte, 0, bsSizeDef); return &bs }}   // 普通字节序列 < 8K
	bsPoolLarge = sync.Pool{New: func() any { bs := make([]byte, 0, bsSizeLarge); return &bs }} // 大字节序列 	< 64K
	bsPoolHuge  = sync.Pool{New: func() any { bs := make([]byte, 0, bsSizeHuge); return &bs }}  // 超大字节序列 < 512KB
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
	if needSize <= bsSizeMini {
		bf = bsPoolMini.Get().(*[]byte)
	} else if needSize <= bsSizeDef {
		bf = bsPoolDef.Get().(*[]byte)
	} else if needSize <= bsSizeLarge {
		bf = bsPoolLarge.Get().(*[]byte)
	} else {
		bf = bsPoolHuge.Get().(*[]byte)
	}
	*bf = (*bf)[0:0]
	return bf
}

func FreeBytes(bs *[]byte) {
	lastSize = cap(*bs)
	if lastSize > bsSizeMax {
		return
	} else if lastSize >= bsSizeHuge {
		bsPoolHuge.Put(bs)
	} else if lastSize >= bsSizeLarge {
		bsPoolLarge.Put(bs)
	} else if lastSize >= bsSizeDef {
		bsPoolDef.Put(bs)
	} else {
		bsPoolMini.Put(bs)
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
