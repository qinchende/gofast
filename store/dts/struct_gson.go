package dts

//
//func GetEncPet(list any) (pet gson.RowsEncPet) {
//	typ := reflect.TypeOf(list)
//	sliceTyp := typ.Elem()
//	stuTyp := sliceTyp.Elem()
//
//	ss := SchemaForDBByType(stuTyp)
//
//	af := (*rt.AFace)(unsafe.Pointer(&list))
//	sh := (*rt.SliceHeader)(af.DataPtr)
//
//	pet.Tt = int64(sh.Len)
//	pet.Fields = strings.Join(ss.cTips.items, ",")
//	pet.FlsIdxes = ss.cTips.idxes
//	pet.Target = list
//	return
//}
