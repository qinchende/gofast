package cdo

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

// 这是通用方法，但不是最高效的
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// List type value
func encAllList(se *subEncode, listSize int) {
	encUint32(se.bf, TypeArray, uint64(listSize))
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
func encAllListPtr(se *subEncode, listSize int) {
	encUint32(se.bf, TypeArray, uint64(listSize))
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

		// add by cd.net on 2023-12-01
		// []*map 类似这种数据源，不要再剥离一次指针，指向map值的内存
		if se.em.itemKind == reflect.Map {
			itemPtr = *(*unsafe.Pointer)(itemPtr)
		}
		se.em.itemEnc(se.bf, itemPtr, se.em.itemType)
	}
}

// number
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encUintList[T constraints.Unsigned](se *subEncode, listSize int) {
	encUint32(se.bf, TypeArray, uint64(listSize))
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		encUint64(se.bf, TypePosInt, uint64(*(*T)(iPtr)))
	}
}

func encIntList[T constraints.Integer](se *subEncode, listSize int) {
	encUint32(se.bf, TypeArray, uint64(listSize))
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		v := *((*T)(iPtr))
		if v >= 0 {
			encUint64(se.bf, TypePosInt, uint64(v))
		} else {
			encUint64(se.bf, TypeNegInt, uint64(-v))
		}
	}
}

func encIntListPtr[T constraints.Integer](se *subEncode, listSize int) {
	encUint32(se.bf, TypeArray, uint64(listSize))
	ptrLevel := se.em.ptrLevel

	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		// peel ptr ---------------------
		ptrCt := ptrLevel
	peelPtr:
		iPtr = *(*unsafe.Pointer)(iPtr)
		if iPtr == nil {
			*se.bf = append(*se.bf, FixNil)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}
		// END peel ---------------------

		v := *((*T)(iPtr))
		if v >= 0 {
			encUint64(se.bf, TypePosInt, uint64(v))
		} else {
			encUint64(se.bf, TypeNegInt, uint64(-v))
		}
	}
}

// string
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encStringList(se *subEncode, listSize int) {
	encUint32(se.bf, TypeArray, uint64(listSize))
	for i := 0; i < listSize; i++ {
		iPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		encString(se.bf, iPtr, nil)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// list item is type of struct
func encStructList(se *subEncode, listSize int) {
	// list size
	encUint32(se.bf, TypeArrSame, uint64(listSize))

	// 字段
	fls := se.em.ss.FieldsAttr
	encNumMax2BWith6(se.bf, ArrSameObjFields, uint64(len(fls)))
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
