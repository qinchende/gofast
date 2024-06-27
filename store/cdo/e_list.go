package cdo

import (
	"github.com/qinchende/gofast/core/rt"
	"golang.org/x/exp/constraints"
	"reflect"
	"unsafe"
)

// Note: add by sdx on 2024-06-06
// 这里将数组和切片的情况合并考虑，简化了代码；
// 但通常我们遇到的都是切片类型，如果分开处理，将能进一步提高约 10% 的性能。
func (se *subEncode) encList() {
	if se.em.isSlice {
		se.slice = *(*rt.SliceHeader)(se.srcPtr)
		se.srcPtr = se.slice.DataPtr
	} else {
		se.slice = rt.SliceHeader{DataPtr: se.srcPtr, Len: se.em.arrLen, Cap: se.em.arrLen}
	}
	se.em.listEnc(se)
}

// 这是通用方法，但不是最高效的
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// List type value
func encListAll(se *subEncode) {
	tLen := se.slice.Len
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(tLen)), ListAny)
	for i := 0; i < tLen; i++ {
		itemPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		// 一些本身就是引用类型的数据，需要找到他们指向值的地址
		// 比如 map | function | channel 等类型
		if se.em.itemKind == reflect.Map {
			itemPtr = *(*unsafe.Pointer)(itemPtr)
		}
		bs = se.em.itemEnc(bs, itemPtr, se.em.itemType)
	}
	*se.bf = bs
}

// List item is ptr
func encListAllPtr(se *subEncode) {
	tLen := se.slice.Len
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(tLen)), ListAny)

	ptrLevel := se.em.ptrLevel
	for i := 0; i < tLen; i++ {
		itemPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		ptrCt := ptrLevel
	peelPtr:
		itemPtr = *(*unsafe.Pointer)(itemPtr)
		if itemPtr == nil {
			bs = append(bs, FixNilMixed)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

		if se.em.itemKind == reflect.Map {
			itemPtr = *(*unsafe.Pointer)(itemPtr)
		}
		bs = se.em.itemEnc(bs, itemPtr, se.em.itemType)
	}
	*se.bf = bs
}

// int numbers
// Note：整形数组，用第一个字符的第一个bit位来代表正负符号
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListVarUint[T constraints.Unsigned](se *subEncode) {
	list := *(*[]T)(unsafe.Pointer(&se.slice))
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(list))), ListVarInt)
	for i := 0; i < len(list); i++ {
		v := uint64(list[i])
		if v <= MaxUint24 {
			bs = encListVarIntPart1(bs, ListVarIntPos, v)
		} else {
			bs = encListVarIntPart2(bs, ListVarIntPos, v)
		}
	}
	*se.bf = bs
}

func encListVarInt[T constraints.Integer](se *subEncode) {
	list := *(*[]T)(unsafe.Pointer(&se.slice))
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(list))), ListVarInt)
	for i := 0; i < len(list); i++ {
		v := list[i]
		if v >= 0 {
			if uint64(v) <= MaxUint24 {
				bs = encListVarIntPart1(bs, ListVarIntPos, uint64(v))
			} else {
				bs = encListVarIntPart2(bs, ListVarIntPos, uint64(v))
			}
		} else {
			if uint64(-v) <= MaxUint24 {
				bs = encListVarIntPart1(bs, ListVarIntNeg, uint64(-v))
			} else {
				bs = encListVarIntPart2(bs, ListVarIntNeg, uint64(-v))
			}
		}
	}
	*se.bf = bs
}

func encListVarIntPtr[T constraints.Integer](se *subEncode) {
	bs := *se.bf
	tLen := se.slice.Len
	bs = append(encU24By5Ret(bs, TypeList, uint64(tLen)), ListAny)
	ptrLevel := se.em.ptrLevel
	for i := 0; i < tLen; i++ {
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
			if uint64(v) <= MaxUint24 {
				bs = encListVarIntPart1(bs, ListVarIntPos, uint64(v))
			} else {
				bs = encListVarIntPart2(bs, ListVarIntPos, uint64(v))
			}
		} else {
			if uint64(-v) <= MaxUint24 {
				bs = encListVarIntPart1(bs, ListVarIntNeg, uint64(-v))
			} else {
				bs = encListVarIntPart2(bs, ListVarIntNeg, uint64(-v))
			}
		}
	}
	*se.bf = bs
}

// float
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListF32(se *subEncode) {
	list := *(*[]float32)(unsafe.Pointer(&se.slice))
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(list))), ListF32)
	for i := 0; i < len(list); i++ {
		bs = encF32ValRet(bs, list[i])
	}
	*se.bf = bs
}

func encListF64(se *subEncode) {
	list := *(*[]float64)(unsafe.Pointer(&se.slice))
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(list))), ListF64)
	for i := 0; i < len(list); i++ {
		bs = encF64ValRet(bs, list[i])
	}
	*se.bf = bs
}

// bool
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListBool(se *subEncode) {
	list := *(*[]bool)(unsafe.Pointer(&se.slice))
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(list))), ListBool)
	for i := 0; i < len(list); i++ {
		if list[i] {
			bs = append(bs, FixTrue)
		} else {
			bs = append(bs, FixFalse)
		}
	}
	*se.bf = bs
}

// string
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func encListString(se *subEncode) {
	list := *(*[]string)(unsafe.Pointer(&se.slice))
	bs := *se.bf
	bs = append(encU24By5Ret(bs, TypeList, uint64(len(list))), ListStr)
	for i := 0; i < len(list); i++ {
		// Note: 不要改变这里的任何逻辑
		// 这已经是测试过性能最好的写法，因为太长的函数将不会被内联优化
		v := uint64(len(list[i]))
		if v <= MaxUint24 {
			bs = encU32By6RetPart1(bs, TypeStr, v)
		} else {
			bs = encU32By6RetPart2(bs, TypeStr, v)
		}
		bs = append(bs, list[i]...)
	}
	*se.bf = bs
}

// []struct & []*struct
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// list item is type of struct
func encListStruct(se *subEncode) {
	// list size
	tLen := se.slice.Len
	bs := *se.bf
	bs = encU24By5Ret(bs, TypeList, uint64(tLen))

	// 字段名称
	fls := se.em.ss.FieldsAttr
	bs = encU16By6Ret(bs, ListObjFields, uint64(len(fls)))
	bs = encStringsDirectRet(bs, se.em.ss.Columns)

	// 循环记录 + 循环字段
	for i := 0; i < tLen; i++ {
		rPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)
		for j := 0; j < len(fls); j++ {
			bs = se.em.fieldsEnc[j](bs, fls[j].MyPtr(rPtr), fls[j].Type)
		}
	}
	*se.bf = bs
}

func encListStructPtr(se *subEncode) {
	// list size
	tLen := se.slice.Len
	bs := *se.bf
	bs = encU24By5Ret(bs, TypeList, uint64(tLen))

	// 字段名称
	fls := se.em.ss.FieldsAttr
	bs = encU16By6Ret(bs, ListObjFields, uint64(len(fls)))
	bs = encStringsDirectRet(bs, se.em.ss.Columns)

	// 循环记录
	for i := 0; i < tLen; i++ {
		rPtr := unsafe.Add(se.srcPtr, i*se.em.itemMemSize)

		// []*Struct的时候，需要判断值是否为 nil
		if se.em.isPtr {
			ptrCt := se.em.ptrLevel
		peelItemPtr:
			rPtr = *(*unsafe.Pointer)(rPtr)
			if rPtr == nil {
				bs = append(bs, FixNilMixed)
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

			// --- struct field is ptr ---
			ptrCt := fls[j].PtrLevel
			if ptrCt == 0 {
				goto encObjValue
			}

		peelFieldPtr:
			fPtr = *(*unsafe.Pointer)(fPtr)
			if fPtr == nil {
				bs = append(bs, FixNil)
				continue
			}
			ptrCt--
			if ptrCt > 0 {
				goto peelFieldPtr
			}
			// ----------------------------

		encObjValue:
			bs = se.em.fieldsEnc[j](bs, fPtr, fls[j].Type)
		}
	}
	*se.bf = bs
}
