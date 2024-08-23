package cdo

import (
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/rt"
	"unsafe"
)

func (d *decoder) decMap() {
	if d.dm.isMapStrStr {
		d.decMapStrStr()
	} else if d.dm.isMapSKV {
		d.decMapStrAny()
	} else {
		d.decMapOthers()
	}
}

func (d *decoder) decMapStrStr() {
	tLen, pos := d.kvLenPos()
	for i := 0; i < tLen; i++ {
		off1, k := scanString(d.str[pos:])
		pos += off1
		off2, v := scanString(d.str[pos:])
		pos += off2
		d.skv.SetString(k, v)
	}
	d.scan = pos
}

func (d *decoder) decMapStrAny() {
	tLen, pos := d.kvLenPos()
	d.scan = pos
	for i := 0; i < tLen; i++ {
		d.skv.Set(decStrVal(d), d.decAny())
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (d *decoder) decMapOthers() {
	tLen, pos := d.kvLenPos()
	d.scan = pos

	mapPtr := d.dstPtr
	mapMem := *(*unsafe.Pointer)(mapPtr)
	if mapMem == nil {
		mapMem = rt.MakeMap(d.dm.mapTypeAbi, 0)
	}
	for i := 0; i < tLen; i++ {
		// key
		k := rt.UnsafeNew(d.dm.mapKTypAbi)
		d.dstPtr = k
		d.dm.keyDec(d)

		// value
		v := rt.UnsafeNew(d.dm.mapVTypAbi)
		if !d.dm.isPtr {
			d.dstPtr = v
		} else {
			d.dstPtr = getPtrValAddr(v, d.dm.ptrLevel, d.dm.itemTypeAbi)
		}
		d.dm.itemDec(d)
		rt.MapAssign(d.dm.mapTypeAbi, mapMem, k, v)
	}
	*(*unsafe.Pointer)(mapPtr) = mapMem
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func makeKV(ptr unsafe.Pointer) cst.SuperKV {
	if *(*unsafe.Pointer)(ptr) == nil {
		newMap := make(cst.KV)
		*(*unsafe.Pointer)(ptr) = *((*unsafe.Pointer)(unsafe.Pointer(&newMap)))
	}
	return (*cst.KV)(ptr)
}

func makeMapStrAny(ptr unsafe.Pointer) cst.SuperKV {
	if *(*unsafe.Pointer)(ptr) == nil {
		newMap := make(map[string]any)
		*(*unsafe.Pointer)(ptr) = *((*unsafe.Pointer)(unsafe.Pointer(&newMap)))
	}
	return (*cst.KV)(ptr)
}

func makeWebKV(ptr unsafe.Pointer) cst.SuperKV {
	if *(*unsafe.Pointer)(ptr) == nil {
		newMap := make(cst.WebKV)
		*(*unsafe.Pointer)(ptr) = *((*unsafe.Pointer)(unsafe.Pointer(&newMap)))
	}
	return (*cst.WebKV)(ptr)
}

func makeMapStrStr(ptr unsafe.Pointer) cst.SuperKV {
	if *(*unsafe.Pointer)(ptr) == nil {
		newMap := make(map[string]string)
		*(*unsafe.Pointer)(ptr) = *((*unsafe.Pointer)(unsafe.Pointer(&newMap)))
	}
	return (*cst.WebKV)(ptr)
}
