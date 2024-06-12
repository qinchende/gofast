package cdo

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

// 这是通用方法，但不是最高效的
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// List type value
func encListAll(se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))
	*se.bf = append(*se.bf, ListAny)

	for i := 0; i < listSize; i++ {
		itemPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		// 一些本身就是引用类型的数据，需要找到他们指向值的地址
		// 比如 map | function | channel 等类型
		if se.em.itemKind == reflect.Map {
			itemPtr = *(*unsafe.Pointer)(itemPtr)
		}
		se.em.itemEnc(se.bf, itemPtr, se.em.itemType)
	}
}

// List item is ptr
func encListAllPtr(se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))
	*se.bf = append(*se.bf, ListAny)

	ptrLevel := se.em.ptrLevel
	for i := 0; i < listSize; i++ {
		itemPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		ptrCt := ptrLevel
	peelPtr:
		itemPtr = *(*unsafe.Pointer)(itemPtr)
		if itemPtr == nil {
			*se.bf = append(*se.bf, FixMixedNil)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

		if se.em.itemKind == reflect.Map {
			itemPtr = *(*unsafe.Pointer)(itemPtr)
		}
		se.em.itemEnc(se.bf, itemPtr, se.em.itemType)
	}
}

// int
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListUint[T constraints.Unsigned](se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))

	bs := *se.bf
	bs = append(bs, ListVarInt)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		v := uint64(*(*T)(iPtr))
		if v <= Max3BUint {
			bs = encU64By6RetPart1(bs, TypePosInt, v)
		} else {
			bs = encU64By6RetPart2(bs, TypePosInt, v)
		}
	}
	*se.bf = bs
}

func encListInt[T constraints.Integer](se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))

	bs := *se.bf
	bs = append(bs, ListVarInt)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		v := *((*T)(iPtr))
		if v >= 0 {
			if uint64(v) <= Max3BUint {
				bs = encU64By6RetPart1(bs, TypePosInt, uint64(v))
			} else {
				bs = encU64By6RetPart2(bs, TypePosInt, uint64(v))
			}
		} else {
			if uint64(-v) <= Max3BUint {
				bs = encU64By6RetPart1(bs, TypeNegInt, uint64(-v))
			} else {
				bs = encU64By6RetPart2(bs, TypeNegInt, uint64(-v))
			}
		}
	}
	*se.bf = bs
}

func encListIntPtr[T constraints.Integer](se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))
	ptrLevel := se.em.ptrLevel

	bs := *se.bf
	bs = append(bs, ListVarInt)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		// peel ptr ---------------------
		ptrCt := ptrLevel
	peelPtr:
		iPtr = *(*unsafe.Pointer)(iPtr)
		if iPtr == nil {
			bs = append(bs, FixNil)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}
		// END peel ---------------------

		v := *((*T)(iPtr))
		if v >= 0 {
			if uint64(v) <= Max3BUint {
				bs = encU64By6RetPart1(bs, TypePosInt, uint64(v))
			} else {
				bs = encU64By6RetPart2(bs, TypePosInt, uint64(v))
			}
		} else {
			if uint64(-v) <= Max3BUint {
				bs = encU64By6RetPart1(bs, TypeNegInt, uint64(-v))
			} else {
				bs = encU64By6RetPart2(bs, TypeNegInt, uint64(-v))
			}
		}
	}
	*se.bf = bs
}

// float
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListF32(se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))

	bs := *se.bf
	bs = append(bs, ListF32)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		bs = encF32ValRet(bs, iPtr)
	}
	*se.bf = bs
}

func encListF64(se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))

	bs := *se.bf
	bs = append(bs, ListF64)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		bs = encF64ValRet(bs, iPtr)
	}
	*se.bf = bs
}

// bool
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListBool(se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))

	bs := *se.bf
	bs = append(bs, ListBool)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		bs = encBoolRet(bs, iPtr)
	}
	*se.bf = bs
}

// string
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListString(se *subEncode, listSize int) {
	encU24By5(se.bf, TypeList, uint64(listSize))

	bs := *se.bf
	bs = append(bs, ListStr)
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		str := *((*string)(iPtr))

		// Note: 不要改变这里的任何逻辑
		// 这已经是测试过性能最好的写法，因为太长的函数将不会被内联优化
		v := uint64(len(str))
		if v <= Max3BUint {
			bs = encU32By6RetPart1(bs, TypeStr, v)
		} else {
			bs = encU32By6RetPart2(bs, TypeStr, v)
		}
		bs = append(bs, str...)
	}
	*se.bf = bs
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// list item is type of struct
func encListStruct(se *subEncode, listSize int) {
	// list size
	encU24By5(se.bf, TypeList, uint64(listSize))

	// 字段
	fls := se.em.ss.FieldsAttr
	encU16By6(se.bf, ListObjFields, uint64(len(fls)))
	encStringsDirect(se.bf, se.em.ss.Columns)

	// 循环记录
	for i := 0; i < listSize; i++ {
		rPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		if se.em.isPtr {
			ptrCt := se.em.ptrLevel
		peelItemPtr:
			rPtr = *(*unsafe.Pointer)(rPtr)
			if rPtr == nil {
				*se.bf = append(*se.bf, FixMixedNil)
				continue
			}
			ptrCt--
			if ptrCt > 0 {
				goto peelItemPtr
			}
		}

		// 循环字段
		for j := 0; j < len(fls); j++ {
			fPtr := fls[j].MyPtr(rPtr)

			ptrCt := fls[j].PtrLevel
			if ptrCt == 0 {
				goto encObjValue
			}

		peelFieldPtr:
			fPtr = *(*unsafe.Pointer)(fPtr)
			if fPtr == nil {
				*se.bf = append(*se.bf, FixNil)
				continue
			}
			ptrCt--
			if ptrCt > 0 {
				goto peelFieldPtr
			}

		encObjValue:
			se.em.fieldsEnc[j](se.bf, fPtr, fls[j].Type)
		}
	}
}
