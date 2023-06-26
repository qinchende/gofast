package jde

import (
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
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
	//arrIdx    int  // list解析的数量
	//skipValue bool // 跳过当前要解析的值

	flsIdx []int8
}

func decGsonRowsFromString(v any, str string) {
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

func (gd *gsonDecode) scanGsonRows() {
	//pls := jdeBufPool.Get().(*listPool)
	gd.flsIdx = make([]int8, 0, 100)

	sd := gd.sd
	pos := sd.scan

	for isBlankChar[sd.str[pos]] {
		pos++
	}

	c := sd.str[pos]
	if c != '[' {
		goto errChar
	}
	pos++

	// 0. Current count ++++++++++++++++++++++++++++++++++++++++
	_ = sd.scanJsonRowUint()
	// 1. Total count  +++++++++++++++++++++++++++++++++++++++++
	_ = sd.scanJsonRowUint()

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
			break
		} else if sd.arrIdx > 0 {
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

	for {
		// 不用switch, 比较顺序相对比较明确
		if c == ',' {
			pos++
		} else if c == ']' {
			break
		} else if sd.arrIdx > 0 {
			goto errChar
		}

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}

		// scan field
		sd.scan = pos
		gd.scanJsonRowRecode()
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

func (gd *gsonDecode) scanJsonRowRecode() {
	ct := 0

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
			break
		} else if sd.arrIdx > 0 {
			goto errChar
		}

		c = sd.str[pos]
		for isBlankChar[c] {
			pos++
			c = sd.str[pos]
		}

		// scan field
		idx := gd.flsIdx[ct]
		if idx >= 0 {
			sd.dm.fieldsDec[idx](sd)
		} else {
			sd.skipOneValue()
		}
		ct++

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
