package mapx

import "reflect"

func checkDestType(dest any) (reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		panic("dest must be pointer.")
	}
	dSliceTyp := dTyp.Elem()
	if dSliceTyp.Kind() != reflect.Slice {
		panic("dest must be slice.")
	}

	isPtr := false
	isKV := false
	dItemType := dSliceTyp.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if dItemType.Kind() == reflect.Ptr {
		isPtr = true
		dItemType = dItemType.Elem()
	} else if dItemType.Name() == "KV" {
		isKV = true
	}

	return dSliceTyp, dItemType, isPtr, isKV
}
