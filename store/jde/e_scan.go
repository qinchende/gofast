package jde

import (
	"fmt"
	"reflect"
	"runtime/debug"
)

func (se *subEncode) encStart() (err errType) {
	defer func() {
		if pic := recover(); pic != nil {
			if code, ok := pic.(errType); ok {
				err = code
			} else {
				fmt.Printf("%s\n%s", pic, debug.Stack())
				err = errJson
			}
		}
	}()

	if se.dm.isArray {
		se.encArray()
	} else if se.dm.isList {
		se.encList()
	} else {
		se.encObject()
	}
	return
}

func (se *subEncode) encArray() {
	se.bs = append(se.bs, "["...)

	sh := (*reflect.SliceHeader)(se.dstPtr)
	if sh.Len >= 1 {
		se.dm.listItemEnc(se, 0)
	}
	for i := 1; i < sh.Len; i++ {
		se.bs = append(se.bs, ", "...)
		se.dm.listItemEnc(se, i)
	}

	se.bs = append(se.bs, ']')
}

func (se *subEncode) encList() {
	se.bs = append(se.bs, "["...)

	sh := (*reflect.SliceHeader)(se.dstPtr)
	if sh.Len >= 1 {
		se.dm.listItemEnc(se, 0)
	}
	for i := 1; i < sh.Len; i++ {
		se.bs = append(se.bs, ", "...)
		se.dm.listItemEnc(se, i)
	}

	se.bs = append(se.bs, ']')
}

func (se *subEncode) encObject() {

}
