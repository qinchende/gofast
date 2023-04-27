package jde

//
//func bindArrValue[T string | bool | float32 | float64](a *listPost, v T) {
//	*(*T)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = v
//	a.arrIdx++
//}
//
//var (
//	kindIntFunc = [27]arrIntFunc{
//		2: func(a *listPost, v int64) {
//			*(*int)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = int(v)
//			a.arrIdx++
//		},
//		3: func(a *listPost, v int64) {
//			*(*int8)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = int8(v)
//			a.arrIdx++
//		},
//		4: func(a *listPost, v int64) {
//			*(*int16)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = int16(v)
//			a.arrIdx++
//		},
//		5: func(a *listPost, v int64) {
//			*(*int32)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = int32(v)
//			a.arrIdx++
//		},
//		6: func(a *listPost, v int64) {
//			*(*int64)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = v
//			a.arrIdx++
//		},
//
//		7: func(a *listPost, v int64) {
//			*(*uint)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = uint(v)
//			a.arrIdx++
//		},
//		8: func(a *listPost, v int64) {
//			*(*uint8)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = uint8(v)
//			a.arrIdx++
//		},
//		9: func(a *listPost, v int64) {
//			*(*uint16)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = uint16(v)
//			a.arrIdx++
//		},
//		10: func(a *listPost, v int64) {
//			*(*uint32)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = uint32(v)
//			a.arrIdx++
//		},
//		11: func(a *listPost, v int64) {
//			*(*uint64)(unsafe.Pointer(a.arrPtr + uintptr(a.arrIdx*a.itemSize))) = uint64(v)
//			a.arrIdx++
//		},
//	}
//)
