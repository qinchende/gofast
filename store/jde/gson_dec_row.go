package jde

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"reflect"
	"runtime/debug"
	"unsafe"
)

// 解析GsonRow的值部分
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func decGsonRowOnlyValues(obj any, str string) (err error) {
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
	if dm = cacheGetDecMetaFast(af.TypePtr); dm == nil {
		// +++++++++++++ check type
		rfType := reflect.TypeOf(obj)
		if rfType.Kind() != reflect.Pointer {
			panic(errValueMustPtr)
		}
		rfType = rfType.Elem()
		if rfType.Kind() != reflect.Struct {
			panic(errValueMustStruct)
		}

		if dm = cacheGetDecMeta(rfType); dm == nil {
			dm = newDecodeMeta(rfType)
			cacheSetDecMeta(rfType, dm)
		}
		cacheSetDecMetaFast(af.TypePtr, dm)
	}

	// +++++++++++++++++++++++++++++++++++++++
	sd := jdeDecPool.Get().(*subDecode)
	sd.reset()
	sd.dm = dm
	sd.str = str
	sd.dstPtr = af.DataPtr
	sd.scanGsonRowJustValues()
	jdeDecPool.Put(sd)

	// NOTE：我们是特殊解析GsonRow的值部分，如果JSON字符串没有正常结束，需要报错。
	if sd.scan != len(str) {
		err = errJsonRowStr
	}
	return
}

// 这里待解析字符串的形式只能是 str: -> [v1,v2,v3,...]
func (sd *subDecode) scanGsonRowJustValues() {
	flsCount := len(sd.dm.ss.Fields)

	if sd.str[sd.scan] != '[' {
		panic(errChar)
	}
	sd.scan++

	fIndex := 0
	for {
		if sd.str[sd.scan] == ',' {
			sd.scan++
		} else if sd.str[sd.scan] == ']' {
			sd.scan++
			return
		} else if fIndex > 0 {
			panic(errList)
		}

		if fIndex < flsCount {
			sd.keyIdx = fIndex
			sd.dm.fieldsDec[fIndex](sd)
		} else {
			sd.skipOneValue()
		}

		fIndex++
	}
}
