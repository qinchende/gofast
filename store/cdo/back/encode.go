package back

//
//func encU16By6(bf *[]byte, sym byte, v uint64) {
//	*bf = encU16By6Ret(*bf, sym, v)
//}
//func encU32By6(bf *[]byte, sym byte, v uint64) {
//	*bf = encU32By6Ret(*bf, sym, v)
//}
//func encU64By6(bf *[]byte, sym byte, v uint64) {
//	*bf = encU64By6Ret(*bf, sym, v)
//}
//func encU24By5(bf *[]byte, sym byte, v uint64) {
//	*bf = encU24By5Ret(*bf, sym, v)
//}
//func encBytes(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	bs := *((*[]byte)(ptr))
//	encU32By6(bf, TypeStr, uint64(len(bs)))
//	*bf = append(*bf, bs...)
//}
//
//func encInt[T constraints.Integer](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	v := *((*T)(ptr))
//	if v >= 0 {
//		encU64By6(bf, TypeVarIntPos, uint64(v))
//	} else {
//		encU64By6(bf, TypeVarIntNeg, uint64(-v))
//	}
//}
//
//func encUint[T constraints.Unsigned](bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	encU64By6(bf, TypeVarIntPos, uint64(*((*T)(ptr))))
//}
//
//func encF32(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	v := *(*uint32)(ptr)
//	*bf = append(*bf, FixF32, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
//}
//
//func encF64(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	v := *(*uint64)(ptr)
//	*bf = append(*bf, FixF64, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
//}
//
//func encNil(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	*bf = append(*bf, FixNil)
//}
//
//func encBool(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	if *((*bool)(ptr)) {
//		*bf = append(*bf, FixTrue)
//	} else {
//		*bf = append(*bf, FixFalse)
//	}
//}
//
//func encString(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	str := *((*string)(ptr))
//	encU32By6(bf, TypeStr, uint64(len(str)))
//	*bf = append(*bf, str...)
//}
//
//func encTime(bf *[]byte, ptr unsafe.Pointer, typ reflect.Type) {
//	tp := *bf
//	tp = append(tp, FixTime)
//	*bf = append(tp, (*time.Time)(ptr).Format(cst.TimeFmtRFC3339)...)
//}
