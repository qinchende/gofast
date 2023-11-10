package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"runtime/debug"
	"sync"
	"unsafe"
)

var (
	grsDecPool        = sync.Pool{New: func() any { return &gsonRowsDecode{} }}
	cachedGsonDecMeta sync.Map
)

func cacheSetGsonDecMeta(typAddr *rt.TypeAgent, val *decMeta) {
	cachedGsonDecMeta.Store(typAddr, val)
}

func cacheGetGsonDecMeta(typAddr *rt.TypeAgent) *decMeta {
	if ret, ok := cachedGsonDecMeta.Load(typAddr); ok {
		return ret.(*decMeta)
	}
	return nil
}

type gsonRowsDecode struct {
	sd subDecode // 共享的subDecode，用来解析子对象
	ct int64     // 记录条数
	tt int64     // 当前查询条件总记录数（用于分页）

	fc     int       // 字段数量
	flsIdx [128]int8 // 结构体不能超过128个字段（当然这里可以改大，不过建议不要定义那么多字段的结构体）
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Struct Json 解析含有 GsonRows字段，其它部分还需要按JSON格式继续解析
func scanObjGsonPet(sd *subDecode) {
	pet := (*gson.RowsDecPet)(fieldPtr(sd))
	ret := decGsonRows(pet.Target, sd.str[sd.scan:])
	if ret.Err != nil {
		panic(ret.Err)
	}

	pet.Ct = ret.Ct
	pet.Tt = ret.Tt
	sd.scan += ret.Scan
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decGsonRows(v any, str string) (ret gson.RowsDecRet) {
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				ret.Err = errors.New(fmt.Sprintf("error code: %d", code))
			} else {
				// 调试的时候打印错误信息
				fmt.Printf("%s\n%s", pic, debug.Stack())
				ret.Err = errors.New(fmt.Sprintf("other panic : %s", pic))
			}
		}
	}()

	af := (*rt.AFace)(unsafe.Pointer(&v))
	var dm *decMeta

	// check target object
	if dm = cacheGetGsonDecMeta(af.TypePtr); dm == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(v)
		if dstTyp.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		sliceType := dstTyp.Elem()
		if sliceType.Kind() != reflect.Slice {
			panic(errValueMustSlice)
		}
		itemType := sliceType.Elem()
		// TODO：只支持struct切片，而不是struct指针切片
		if itemType.Kind() != reflect.Struct {
			panic(errValueMustStruct)
		}

		if dm = cacheGetMeta(itemType); dm == nil {
			dm = newDecodeMeta(itemType)
			cacheSetMeta(itemType, dm)
		}
		cacheSetGsonDecMeta(af.TypePtr, dm)
	}

	// +++++++++++++++++++++++++++++++++++++++
	// var grs gsonRowsDecode
	grs := grsDecPool.Get().(*gsonRowsDecode)

	// init subDecode
	sd := &grs.sd
	sd.reset()
	sd.dm = dm
	sd.str = str
	sd.dstPtr = af.DataPtr
	grs.scanGsonRows()

	ret.Ct = grs.ct
	ret.Tt = grs.tt
	ret.Scan = grs.sd.scan

	grsDecPool.Put(grs)
	return
}

func (grs *gsonRowsDecode) scanGsonRows() {
	sd := &grs.sd
	sh := (*reflect.SliceHeader)(sd.dstPtr)
	flsCT := 0

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
	grs.ct = int64(sd.scanJsonRowUint()) // 0. Current count ++++++++++
	grs.tt = int64(sd.scanJsonRowUint()) // 1. Total count  +++++++++++
	pos = sd.scan

	// 2. Struct fields ++++++++++++++++++++++++++++++++++++++++
	if sd.str[pos] != '[' {
		goto errChar
	}
	pos++

	for {
		c = sd.str[pos]
		if c == ',' {
			pos++
		} else if c == ']' {
			pos++
			break
		} else if flsCT > 0 {
			goto errChar
		}

		if sd.str[pos] != '"' {
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
		grs.flsIdx[flsCT] = int8(sd.dm.ss.ColumnIndex(key)) // 比对 column 名称
		flsCT++
	}
	grs.fc = flsCT // 多少个有效字段

	// 3. values +++++++++++++++++++++++++++++++++++++++++++++++
	if sd.str[pos] != ',' {
		goto errChar
	}
	pos++
	if sd.str[pos] != '[' {
		goto errChar
	}
	pos++

	// 根据记录数量，初始化对象空间
	flsCT = int(grs.ct)
	if flsCT > sh.Cap {
		*(*[]byte)(sd.dstPtr) = make([]byte, sh.Len*sd.dm.itemRawSize, flsCT*sd.dm.itemRawSize)
		sh.Len, sh.Cap = flsCT, flsCT
	} else {
		sh.Len = flsCT
	}

	flsCT = 0
	for {
		c = sd.str[pos]
		if c == ',' {
			pos++
		} else if c == ']' {
			pos++
			goto finished
		}

		sd.scan = pos
		sd.dstPtr = unsafe.Pointer(sh.Data + uintptr(flsCT*sd.dm.itemRawSize))
		//// 如果是指针，需要分配空间
		//if sd.dm.isPtr {
		//	sd.dstPtr = getPtrValueAddr(sd.dstPtr, sd.dm.ptrLevel, sd.dm.itemKind, sd.dm.itemType)
		//}
		grs.scanJsonRowRecode()
		pos = sd.scan

		flsCT++
	}

errChar:
	sd.scan = pos
	panic(errChar)

finished:
	if sd.str[pos] != ']' {
		goto errChar
	}
	pos++

	// NOTE：结束的时候要讲后面的空字符排除掉
	if pos < len(sd.str) {
		for isBlankChar[sd.str[pos]] {
			pos++
		}
	}
	sd.scan = pos
}

func (grs *gsonRowsDecode) scanJsonRowRecode() {
	sd := &grs.sd

	if sd.str[sd.scan] != '[' {
		panic(errChar)
	}
	sd.scan++

	fc := 0
	for {
		if sd.str[sd.scan] == ',' {
			sd.scan++
		} else if sd.str[sd.scan] == ']' {
			sd.scan++
			return
		} else if fc > 0 {
			panic(errList)
		}

		sd.keyIdx = int(grs.flsIdx[fc])
		if sd.keyIdx >= 0 && fc < grs.fc {
			sd.dm.fieldsDec[sd.keyIdx](sd)
		} else {
			sd.skipOneValue()
		}

		fc++
	}
}

func (sd *subDecode) scanJsonRowUint() (ret uint64) {
	ret = lang.ParseUintFast(sd.str[sd.scanUintMust():sd.scan])
	if sd.str[sd.scan] != ',' {
		panic(errChar)
	}
	sd.scan++
	return
}
