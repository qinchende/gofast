package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"runtime/debug"
	"sync"
	"unsafe"
)

var (
	grDecPool     = sync.Pool{New: func() any { return &gsonRowDecode{} }}
	cachedGsonRow sync.Map
)

func cacheSetGsonRow(typAddr *rt.TypeAgent, val *decMeta) {
	cachedGsonRow.Store(typAddr, val)
}

func cacheGetGsonRow(typAddr *rt.TypeAgent) *decMeta {
	if ret, ok := cachedGsonRow.Load(typAddr); ok {
		return ret.(*decMeta)
	}
	return nil
}

type gsonRowDecode struct {
	sd subDecode // 共享的subDecode，用来解析子对象
	fc int       // 字段数量
	//flsIdx []uint8   // 结构体不能超过128个字段（当然这里可以改大，不过建议不要定义那么多字段的结构体）
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decGsonRow(obj any, str string) (err error) {
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

	af := (*rt.AFace)(unsafe.Pointer(&obj))
	var dm *decMeta

	// check target object
	if dm = cacheGetGsonRow(af.TypePtr); dm == nil {
		// +++++++++++++ check type
		dstTyp := reflect.TypeOf(obj)
		if dstTyp.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		objType := dstTyp.Elem()
		if objType.Kind() != reflect.Struct {
			panic(errValueMustStruct)
		}

		typAddr := (*rt.TypeAgent)((*rt.AFace)(unsafe.Pointer(&objType)).DataPtr)
		if dm = cacheGetMeta(typAddr); dm == nil {
			dm = newDecodeMeta(objType)
			cacheSetMeta(typAddr, dm)
		}
		cacheSetGsonRow(af.TypePtr, dm)
	}

	// +++++++++++++++++++++++++++++++++++++++
	// var gr gsonRowsDecode
	gr := grDecPool.Get().(*gsonRowDecode)
	// init subDecode
	sd := &gr.sd
	sd.reset()
	sd.dm = dm
	sd.str = str
	sd.dstPtr = af.DataPtr

	gr.fc = len(dm.ss.Columns)
	gr.scanJsonRowValueString()

	grDecPool.Put(gr)

	if gr.sd.scan != len(str) {
		err = errJsonRowStr
	}
	return
}

func (gr *gsonRowDecode) scanJsonRowValueString() {
	sd := &gr.sd

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

		//sd.keyIdx = int(gr.flsIdx[fc])
		if fc < gr.fc {
			sd.dm.fieldsDec[fc](sd)
		} else {
			sd.skipOneValue()
		}

		fc++
	}
}
