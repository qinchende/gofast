package dts

//
//func ExtractListType(dst any) (reflect.Type, reflect.Type, bool, bool) {
//	dstTyp := reflect.TypeOf(dst)
//	if dstTyp.Kind() != reflect.Pointer {
//		cst.PanicString("Target object must be pointer.")
//	}
//	sliceType := dstTyp.Elem()
//	if sliceType.Kind() != reflect.Slice {
//		cst.PanicString("Target object must be slice.")
//	}
//
//	isPtr := false
//	isKV := false
//	recordType := sliceType.Elem()
//	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
//	if recordType.Kind() == reflect.Pointer {
//		isPtr = true
//		recordType = recordType.Elem()
//	} else {
//		typStr := recordType.String()
//		if typStr == cst.StrTypeOfKV {
//			isKV = true
//		}
//	}
//
//	return sliceType, recordType, isPtr, isKV
//}
