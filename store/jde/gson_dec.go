package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"runtime/debug"
	"unsafe"
)

type gsonDecode struct {
	sd subDecode // 共享的subDecode，用来解析子对象
	//dm *decMeta  // Struct | Slice,Array

	//dstPtr unsafe.Pointer // 目标值dst的地址
	//// 当前解析JSON的状态信息 ++++++
	//str  string // 本段字符串
	//scan int    // 自己的扫描进度，当解析错误时，这个就是定位
	//
	//keyIdx    int  // key index
	//arrIdx int // list解析的数量
	//skipValue bool // 跳过当前要解析的值

	ct uint64
	tt uint64

	flsIdx []int8
}

func decGsonRowsFromString(v any, str string) (err error) {
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = errors.New(fmt.Sprintf("error code: %d", code))
			} else {
				// 调试的时候打印错误信息
				fmt.Printf("%s\n%s", pic, debug.Stack())
				err = errors.New(fmt.Sprintf("other panic : %s", pic))
			}
		}
	}()

	dstTyp := reflect.TypeOf(v)
	if dstTyp.Kind() != reflect.Pointer {
		cst.PanicString("Target object must be pointer.")
	}
	sliceType := dstTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		cst.PanicString("Target object must be slice.")
	}

	var gd gsonDecode
	itemType := sliceType.Elem()
	typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&itemType)).DataPtr)
	if meta := cacheGetMeta(typAddr); meta != nil {
		gd.sd.dm = meta
	} else {
		gd.sd.dm = newDecodeMeta(itemType)
		cacheSetMeta(typAddr, gd.sd.dm)
	}
	gd.sd.str = str

	ptr := (*rt.AFace)(unsafe.Pointer(&v)).DataPtr
	gd.sd.dstPtr = ptr

	gd.scanGsonRows()
	return nil
}

func (gd *gsonDecode) scanGsonRows() {
	//pls := jdeBufPool.Get().(*listPool)
	gd.flsIdx = make([]int8, 0, 100)
	sd := &gd.sd
	sh := (*reflect.SliceHeader)(sd.dstPtr)
	recordIdx := 0
	newLen := 0

	pos := sd.scan

	for isBlankChar[sd.str[pos]] {
		pos++
	}

	c := sd.str[pos]
	if c != '[' {
		goto errChar
	}
	pos++

	sd.scan = pos

	// 0. Current count ++++++++++++++++++++++++++++++++++++++++
	gd.ct = sd.scanJsonRowUint()
	// 1. Total count  +++++++++++++++++++++++++++++++++++++++++
	gd.tt = sd.scanJsonRowUint()

	pos = sd.scan

	// 2. Struct fields ++++++++++++++++++++++++++++++++++++++++
	c = sd.str[pos]
	if c != '[' {
		goto errChar
	}
	pos++

	for {
		// 不用switch, 比较顺序相对比较明确
		if c == ',' {
			pos++
		} else if c == ']' {
			pos++
			break
		} else if len(gd.flsIdx) > 0 {
			goto errChar
		}

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}

		// scan field
		if c != '"' {
			goto errChar
		}

		start := pos
		sd.scan = pos
		slash := sd.scanQuoteStr()
		pos = sd.scan

		var key string
		if slash {
			key = sd.str[start+1 : sd.unescapeEnd()]
		} else {
			key = sd.str[start+1 : pos-1]
		}
		gd.flsIdx = append(gd.flsIdx, int8(sd.dm.ss.ColumnIndex(key)))

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}
	}

	// 3. values +++++++++++++++++++++++++++++++++++++++++++++++
	for isBlankChar[sd.str[pos]] {
		pos++
	}
	c = sd.str[pos]
	if c != ',' {
		goto errChar
	}
	pos++

	c = sd.str[pos]
	if c != '[' {
		goto errChar
	}
	pos++
	for isBlankChar[sd.str[pos]] {
		pos++
	}

	newLen = int(gd.ct)

	if newLen > sh.Cap {
		//var oldMem = []byte{}
		//old := (*reflect.SliceHeader)(unsafe.Pointer(&oldMem))
		//old.Data, old.Len, old.Cap = sh.Data, sh.Len*sd.dm.itemRawSize, sh.Cap*sd.dm.itemRawSize

		// TODO: 需要有更高效的扩容算法
		//newLen := int(float64(sh.Cap)*1.6) + 4 // 一种简易的动态扩容算法
		//fmt.Printf("growing len: %d, cap: %d \n\n", sh.Len*sd.dm.itemRawSize, newLen*sd.dm.itemRawSize)
		*(*[]byte)(sd.dstPtr) = make([]byte, sh.Len*sd.dm.itemRawSize, newLen*sd.dm.itemRawSize)

		//copy(*(*[]byte)(sd.dstPtr), oldMem)
		sh.Len, sh.Cap = newLen, newLen
	}

	for {
		// 不用switch, 比较顺序相对比较明确
		if recordIdx > 0 && c == ',' {
			pos++
		} else if c == ']' {
			sd.scan = pos + 1
			return
		}

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}

		// scan field
		sd.scan = pos
		sd.dstPtr = unsafe.Pointer(sh.Data + uintptr(recordIdx*sd.dm.itemRawSize))
		// TODO: slice index ptr
		gd.scanJsonRowRecode()
		recordIdx++
		pos = sd.scan

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}
	}

errChar:
	sd.scan = pos
	panic(errChar)
}

func (sd *subDecode) scanJsonRowUint() (ret uint64) {
	pos := sd.scan

	// 第一个整数
	for isBlankChar[sd.str[pos]] {
		pos++
	}

	sd.scan = pos
	ret = lang.ParseUint(sd.str[sd.scanUintMust():sd.scan])
	pos = sd.scan

	for isBlankChar[sd.str[pos]] {
		pos++
	}

	c := sd.str[pos]
	if c != ',' {
		goto errChar
	}
	pos++

	sd.scan = pos
	return

errChar:
	sd.scan = pos
	panic(errChar)
}

func (gd *gsonDecode) scanJsonRowRecode() {
	fc := 0
	sd := &gd.sd

	pos := sd.scan

	c := sd.str[pos]
	if c != '[' {
		goto errChar
	}
	pos++

	for {
		// 不用switch, 比较顺序相对比较明确
		if c == ',' {
			pos++
		} else if c == ']' {
			sd.scan = pos + 1
			return
		} else if fc > 0 {
			goto errChar
		}

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}

		// scan field
		sd.keyIdx = int(gd.flsIdx[fc])
		sd.scan = pos
		if sd.keyIdx >= 0 {
			sd.dm.fieldsDec[sd.keyIdx](sd)
		} else {
			sd.skipOneValue()
		}
		pos = sd.scan
		fc++

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}
	}

errChar:
	sd.scan = pos
	panic(errChar)
}
