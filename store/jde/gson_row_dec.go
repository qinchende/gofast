package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/store/gson"
	"reflect"
	"runtime/debug"
	"unsafe"
)

//var (
//	grsDecPool     = sync.Pool{New: func() any { return &gsonRowsDecode{} }}
//	cachedGsonRows sync.Map
//)
//
//func cacheSetGsonRows(typAddr *rt.TypeAgent, val *decMeta) {
//	cachedGsonRows.Store(typAddr, val)
//}
//
//func cacheGetGsonRows(typAddr *rt.TypeAgent) *decMeta {
//	if ret, ok := cachedGsonRows.Load(typAddr); ok {
//		return ret.(*decMeta)
//	}
//	return nil
//}

type gsonRowDecode struct {
	sd     subDecode // 共享的subDecode，用来解析子对象
	fc     int       // 字段数量
	flsIdx [128]int8 // 结构体不能超过128个字段（当然这里可以改大，不过建议不要定义那么多字段的结构体）
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decGsonRow(obj any, str string) (ret gson.RowsDecRet) {
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

	af := (*rt.AFace)(unsafe.Pointer(&obj))
	var dm *decMeta

	// check target object
	if dm = cacheGetGsonRows(af.TypePtr); dm == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(obj)
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

		typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&itemType)).DataPtr)
		if dm = cacheGetMeta(typAddr); dm == nil {
			dm = newDecodeMeta(itemType)
			cacheSetMeta(typAddr, dm)
		}
		cacheSetGsonRows(af.TypePtr, dm)
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

func (grs *gsonRowDecode) scanJsonRowRecode() {
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
