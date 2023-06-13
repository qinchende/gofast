package jde

import "unsafe"

// cst.KV is map[string]any
// 这是最常见的场景，单独拿出来快速处理
func (se *subEncode) encMapKV() {
	bf := *se.bf

	bf = append(bf, '{')
	theMap := *(*map[string]any)(se.srcPtr)
	for k, v := range theMap {
		bf = append(bf, '"')
		bf = append(bf, k...)
		bf = append(bf, "\":"...)
		bf = encAny(bf, unsafe.Pointer(&v), nil)
	}
	if len(theMap) > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, '}')
	*se.bf = bf
}

//// map type value
//func (se *subEncode) encMapKV() {
//	bf := *se.bf
//
//	//itra := reflect.MakeMap(se.em.itemBaseType).MapRange()
//	//itra.Key()
//
//	theMap := *(*map[string]any)(se.srcPtr)
//	//keyIsStr := se.em.keyBaseType.Kind() == reflect.String
//
//	bf = append(bf, '{')
//	for k, v := range theMap {
//		bf = append(bf, '"')
//		//if keyIsStr {
//		bf = append(bf, k...)
//		//} else {
//		//	bf = se.em.keyPick(bf, unsafe.Pointer(&k), se.em.keyBaseType)
//		//}
//		bf = append(bf, "\":"...)
//
//		//v := theMap[k]
//
//		//ptr := (*rt.AFace)(unsafe.Pointer(&v)).DataPtr
//		ptr := unsafe.Pointer(&v)
//		ptrCt := se.em.ptrLevel
//		if ptrCt == 0 {
//			goto encMapValue
//		}
//
//	peelPtr:
//		ptr = *(*unsafe.Pointer)(ptr)
//		if ptr == nil {
//			bf = append(bf, "null,"...)
//			continue
//		}
//		ptrCt--
//		if ptrCt > 0 {
//			goto peelPtr
//		}
//
//	encMapValue:
//		bf = se.em.itemPick(bf, ptr, se.em.itemBaseType)
//	}
//	if len(theMap) > 0 {
//		bf = bf[:len(bf)-1]
//	}
//	bf = append(bf, '}')
//	*se.bf = bf
//}
