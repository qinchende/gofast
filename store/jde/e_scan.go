package jde

import (
	"reflect"
	"unsafe"
)

func (se *subEncode) encStart() (err errType) {
	if se.em.isArray {
		if se.em.isPtr {
			se.encListPtr(se.em.itemLen)
		} else {
			se.encList(se.em.itemLen)
		}
	} else if se.em.isList {
		sh := (*reflect.SliceHeader)(unsafe.Pointer(se.srcPtr))
		se.srcPtr = unsafe.Pointer(sh.Data)

		if se.em.isPtr {
			se.encListPtr(sh.Len)
		} else {
			se.encList(sh.Len)
		}
	} else {
		if se.em.isPtr {
			se.encObjPtr()
		} else {
			se.encObj()
		}
	}
	return
}

// List
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) encList(size int) {
	bf := *se.bs
	//size := se.em.itemLen

	bf = append(bf, '[')
	for i := 0; i < size; i++ {
		if se.em.listItemEncMix == nil {
			bf = se.em.listItemEnc(bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemBytes)))
		} else {
			bf = se.em.listItemEncMix(bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemBytes)), *se)
		}
		//bf = append(bf, '"')
		////bf = append(bf, *((*string)(unsafe.Pointer(uintptr(se.srcPtr) + uintptr(i*se.em.itemBytes))))...)
		//bf = append(bf, se.em.listItemEnc(unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemBytes)))...)
		//bf = append(bf, "\","...)
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, ']')
	*se.bs = bf
}

func (se *subEncode) encListPtr(size int) {
	bf := *se.bs
	//size := se.em.itemLen
	ptrLevel := se.em.ptrLevel

	bf = append(bf, '[')
	for i := 0; i < size; i++ {
		ptrCt := ptrLevel
		ptr := unsafe.Pointer(uintptr(se.srcPtr) + uintptr(i*se.em.itemBytes))

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			bf = append(bf, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

		if se.em.listItemEncMix == nil {
			bf = se.em.listItemEnc(bf, ptr)
		} else {
			bf = se.em.listItemEncMix(bf, ptr, *se)
		}
		//bf = append(bf, '"')
		////bf = append(bf, *((*string)(ptr))...)
		//bf = append(bf, se.em.listItemEnc(ptr)...)
		//bf = append(bf, "\","...)
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, ']')
	*se.bs = bf
}

// Object
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) encObj() {
	bf := *se.bs

	size := len(se.em.fieldsEnc)

	bf = append(bf, '{')

	fAttrs := se.em.ss.FieldsAttr
	for i := 0; i < size; i++ {
		bf = append(bf, '"')
		bf = append(bf, se.em.ss.FieldName(i)...)
		bf = append(bf, "\":"...)

		if !fAttrs[i].IsMixType {
			bf = se.em.fieldsEnc[i](bf, unsafe.Pointer(uintptr(se.srcPtr)+fAttrs[i].Offset))
		} else {
			bf = se.em.fieldsEncMix[i](bf, unsafe.Pointer(uintptr(se.srcPtr)+fAttrs[i].Offset), *se, i)
		}
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, '}')

	*se.bs = bf
}

func (se *subEncode) encObjPtr() {
	bf := *se.bs
	size := len(se.em.fieldsEnc)

	bf = append(bf, '{')
	fAttrs := se.em.ss.FieldsAttr
	for i := 0; i < size; i++ {
		ptrCt := fAttrs[i].PtrLevel
		ptr := unsafe.Pointer(uintptr(se.srcPtr) + fAttrs[i].Offset)

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			bf = append(bf, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

		bf = append(bf, '"')
		bf = append(bf, se.em.ss.FieldName(i)...)
		bf = append(bf, "\":"...)

		if !fAttrs[i].IsMixType {
			bf = se.em.fieldsEnc[i](bf, ptr)
		} else {
			bf = se.em.fieldsEncMix[i](bf, ptr, *se, i)
		}

		//se.em.fieldsEnc[i](ptr)

		//bf = append(bf, '"')
		////bf = append(bf, *((*string)(ptr))...)
		//bf = append(bf, se.em.listItemEnc(ptr)...)
		//bf = append(bf, "\","...)
	}
	if size > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, '}')
	*se.bs = bf
}

// MixValue
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListMixItem(bf []byte, ptr unsafe.Pointer, se subEncode) []byte {
	se.initMeta(se.em.itemBaseType, ptr)
	*se.bs = bf
	se.encStart()
	return *se.bs
}

func encObjMixItem(bf []byte, ptr unsafe.Pointer, se subEncode, idx int) []byte {
	se.initMeta(se.em.ss.FieldsAttr[idx].Type, ptr)
	*se.bs = bf
	se.encStart()
	return *se.bs
}
