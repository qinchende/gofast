package cdo

import (
	"unsafe"
)

// Basic type value
func (se *encoder) encBasic() {
	*se.bf = se.em.itemEnc(*se.bf, se.srcPtr, se.em.itemType)
}

// Pointer type value
func (se *encoder) encPointer() {
	ptr := se.srcPtr
	bs := *se.bf

	ptrCt := se.em.ptrLevel
peelPtr:
	ptr = *(*unsafe.Pointer)(ptr)
	if ptr == nil {
		bs = append(bs, FixNil)
		return
	}
	ptrCt--
	if ptrCt > 0 {
		goto peelPtr
	}

	bs = se.em.itemEnc(bs, ptr, se.em.itemType)
	*se.bf = bs
}

// A struct object encode same as map[string]any
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (se *encoder) encStruct() {
	fls := se.em.ss.FieldsAttr
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(fls))), ListKV)

	for i := 0; i < len(fls); i++ {
		str := se.em.ss.ColumnName(i)
		v := uint64(len(str))
		if v <= MaxUint24 {
			bs = encU32By6RetPart1(bs, TypeStr, v)
		} else {
			bs = encU32By6RetPart2(bs, TypeStr, v)
		}
		bs = append(bs, str...)

		// --- ptr ---
		fPtr := fls[i].MyPtr(se.srcPtr)
		ptrCt := fls[i].PtrLevel
		if ptrCt == 0 {
			goto encFieldVal
		}

	peelPtr:
		fPtr = *(*unsafe.Pointer)(fPtr)
		if fPtr == nil {
			bs = append(bs, FixNil)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}
		// -----------

	encFieldVal:
		bs = se.em.fieldsEnc[i](bs, fPtr, fls[i].Type)
	}
	*se.bf = bs
}
