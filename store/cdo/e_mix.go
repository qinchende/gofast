package cdo

import (
	"unsafe"
)

// Basic type value
func (se *subEncode) encBasic() {
	se.em.itemEnc(se.bf, se.srcPtr, se.em.itemType)
}

// Pointer type value
func (se *subEncode) encPointer() {
	ptr := se.srcPtr

	ptrCt := se.em.ptrLevel
peelPtr:
	ptr = *(*unsafe.Pointer)(ptr)
	if ptr == nil {
		*se.bf = append(*se.bf, FixNil)
		return
	}
	ptrCt--
	if ptrCt > 0 {
		goto peelPtr
	}

	encMixedItem(se.bf, ptr, se.em.itemType)
}

// A struct object encode same as map[string]any
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *subEncode) encStruct() {
	fls := se.em.ss.FieldsAttr
	size := len(fls)

	encU24By5(se.bf, TypeList, uint64(size))
	*se.bf = append(*se.bf, ListKV)

	for i := 0; i < size; i++ {
		// key
		encStringDirect(se.bf, se.em.ss.ColumnName(i))

		fPtr := fls[i].MyPtr(se.srcPtr)
		ptrCt := fls[i].PtrLevel
		if ptrCt == 0 {
			goto encObjValue
		}

	peelPtr:
		fPtr = *(*unsafe.Pointer)(fPtr)
		if fPtr == nil {
			*se.bf = append(*se.bf, FixNil)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

	encObjValue:
		// value
		se.em.fieldsEnc[i](se.bf, fPtr, fls[i].Type)
	}
}
