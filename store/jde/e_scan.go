package jde

import (
	"reflect"
	"unsafe"
)

func (se *subEncode) encStart() (err errType) {
	if se.em.isArray {
		if se.em.isPtr {
			se.encListPtr(se.em.arrLen)
		} else {
			se.encList(se.em.arrLen)
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
		se.encObject()
	}
	return
}

// List
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) encList(size int) {
	bf := *se.bs
	//size := se.em.arrLen

	bf = append(bf, '[')
	for i := 0; i < size; i++ {
		if se.em.listItemEncMix == nil {
			bf = se.em.listItemEnc(bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemBytes)))
		} else {
			bf = se.em.listItemEncMix(bf, unsafe.Pointer(uintptr(se.srcPtr)+uintptr(i*se.em.itemBytes)), se.em.itemBaseType)
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
	//size := se.em.arrLen
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
			bf = se.em.listItemEncMix(bf, ptr, se.em.itemBaseType)
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
func (se *subEncode) encObject() {
	bf := *se.bs
	size := len(se.em.fieldsEnc)

	bf = append(bf, '{')
	fAttrs := se.em.ss.FieldsAttr
	for i := 0; i < size; i++ {
		bf = append(bf, '"')
		bf = append(bf, se.em.ss.FieldName(i)...)
		bf = append(bf, "\":"...)

		ptr := unsafe.Pointer(uintptr(se.srcPtr) + fAttrs[i].Offset)
		ptrCt := fAttrs[i].PtrLevel
		if ptrCt == 0 {
			goto encObjValue
		}

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

	encObjValue:
		if !fAttrs[i].IsMixType {
			bf = se.em.fieldsEnc[i](bf, ptr)
		} else {
			bf = se.em.fieldsEncMix[i](bf, ptr, fAttrs[i].Type)
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

// Use SubEncode to encode Mix Item Value
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encMixItem(bf []byte, ptr unsafe.Pointer, rfType reflect.Type) []byte {
	se := subEncode{bs: new([]byte)}
	se.initMeta(rfType, ptr)
	*se.bs = bf
	se.encStart()
	*se.bs = append(*se.bs, ',')
	return *se.bs
}
