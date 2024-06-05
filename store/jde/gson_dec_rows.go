package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"runtime/debug"
	"unsafe"
)

type gsonRowsDecode struct {
	sd subDecode // 共享的subDecode，用来解析子对象
	ct int64     // 记录条数
	tt int64     // 当前查询条件总记录数（用于分页）

	clsCt   int       // 字段数量
	clsIdx  [128]int8 // 结构体不能超过128个字段（当然这里可以改大，不过建议不要定义那么多字段的结构体）
	columns []string  // 字段名称
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Struct Json 解析含有 GsonRows字段，其它部分还需要按JSON格式继续解析
func scanObjGsonPet(sd *subDecode) {
	pet := (*gson.RowsDecPet)(fieldPtr(sd))
	ret := decGsonRows(pet.List, sd.str[sd.scan:])
	if ret.Err != nil {
		panic(ret.Err)
	}

	pet.Ct = ret.Ct
	pet.Tt = ret.Tt
	sd.scan += ret.Scan
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decGsonRows(v any, source string) (ret gson.RowsDecRet) {
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

	var dm *decMeta
	af := (*rt.AFace)(unsafe.Pointer(&v))

	// check target object
	if dm = cacheGetDecMetaFast(af.TypePtr); dm == nil {
		// +++++++++++++ check type
		rfType := reflect.TypeOf(v)
		if rfType.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		rfType = rfType.Elem()
		if rfType.Kind() != reflect.Slice {
			panic(errValueMustSlice)
		}
		rfType = rfType.Elem()

		// 支持2种数据源：
		// A. struct B. cst.KV
		kd := rfType.Kind()
		if kd != reflect.Struct && rfType != cst.TypeCstKV {
			panic(errValueMustStruct)
		}

		if dm = cacheGetDecMeta(rfType); dm == nil {
			dm = newDecodeMeta(rfType)
			cacheSetDecMeta(rfType, dm)
		}
		cacheSetDecMetaFast(af.TypePtr, dm)
	}

	// +++++++++++++++++++++++++++++++++++++++
	grs := grsDecPool.Get().(*gsonRowsDecode)
	defer grsDecPool.Put(grs)

	grs.initDecode(dm, af.DataPtr, source)
	grs.scanGsonRows()

	ret.Ct = grs.ct
	ret.Tt = grs.tt
	ret.Scan = grs.sd.scan
	return
}

func (grs *gsonRowsDecode) initDecode(dm *decMeta, ptr unsafe.Pointer, source string) {
	grs.sd.reset()
	grs.sd.dm = dm
	grs.sd.str = source
	grs.sd.dstPtr = ptr
	grs.columns = grs.columns[0:0]
}

func (grs *gsonRowsDecode) scanGsonRows() {
	sd := &grs.sd
	dm := sd.dm
	//sh := (*rt.SliceHeader)(sd.dstPtr)
	tmpCT := 0

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
		} else if tmpCT > 0 {
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
		if dm.isStruct {
			grs.clsIdx[tmpCT] = int8(dm.ss.ColumnIndex(key)) // 比对 column 名称
		} else {
			grs.columns = append(grs.columns, key)
		}
		tmpCT++
	}
	grs.clsCt = tmpCT // 多少个有效字段

	// 3. values +++++++++++++++++++++++++++++++++++++++++++++++
	if sd.str[pos] != ',' {
		goto errChar
	}
	pos++
	if sd.str[pos] != '[' {
		goto errChar
	}
	pos++

	if dm.isStruct {
		// 根据记录数量，初始化对象空间 +++
		tmpCT = int(grs.ct)
		//if tmpCT > sh.Cap {
		//	*(*[]byte)(sd.dstPtr) = make([]byte, sh.Len*dm.itemMemSize, tmpCT*dm.itemMemSize)
		//	sh.Len, sh.Cap = tmpCT, tmpCT
		//} else {
		//	sh.Len = tmpCT
		//}
		ptr := rt.SliceToArray(sd.dstPtr, dm.itemMemSize, tmpCT)
		// END分配内存空间 ++++++++++++++++

		tmpCT = 0
		for {
			c = sd.str[pos]
			if c == ',' {
				pos++
			} else if c == ']' {
				pos++
				goto finished
			}

			sd.scan = pos
			sd.dstPtr = unsafe.Add(ptr, tmpCT*dm.itemMemSize)
			//// 如果是指针，需要分配空间
			//if sd.dm.isPtr {
			//	sd.dstPtr = getPtrValueAddr(sd.dstPtr, sd.dm.ptrLevel, sd.dm.itemKind, sd.dm.itemType)
			//}
			grs.scanStructRecord()
			pos = sd.scan

			tmpCT++
		}
	} else {
		// 根据记录数量，初始化对象空间 +++
		tmpCT = int(grs.ct)
		//if tmpCT > sh.Cap {
		//	*(*[]cst.KV)(sd.dstPtr) = make([]cst.KV, tmpCT)
		//	sh.Len, sh.Cap = tmpCT, tmpCT
		//} else {
		//	sh.Len = tmpCT
		//}
		ptr := rt.SliceToArray(sd.dstPtr, ptrMemSize, tmpCT)
		// END分配内存空间 ++++++++++++++++

		tmpCT = 0
		for {
			c = sd.str[pos]
			if c == ',' {
				pos++
			} else if c == ']' {
				pos++
				goto finished
			}

			sd.scan = pos
			sd.dstPtr = unsafe.Add(ptr, tmpCT*ptrMemSize)

			// 给 cst.KV类型指针 初始化变量
			theMap := make(cst.KV, grs.clsCt)
			sd.mp = &theMap
			*(*unsafe.Pointer)(sd.dstPtr) = *(*unsafe.Pointer)((unsafe.Pointer)(sd.mp))

			grs.scanKVRecord()
			pos = sd.scan

			tmpCT++
		}
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

func (grs *gsonRowsDecode) scanStructRecord() {
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

		sd.keyIdx = int(grs.clsIdx[fc])
		if sd.keyIdx >= 0 && fc < grs.clsCt {
			sd.dm.fieldsDec[sd.keyIdx](sd)
		} else {
			sd.skipOneValue()
		}

		fc++
	}
}

func (grs *gsonRowsDecode) scanKVRecord() {
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
		sd.dm.kvPairDec(sd, grs.columns[fc])
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
