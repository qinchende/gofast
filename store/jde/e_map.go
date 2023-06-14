package jde

import (
	"github.com/qinchende/gofast/core/rt"
	"unsafe"
)

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

// TODO: 非map[string]any需要encode
func (se *subEncode) encMapGeneral() {
	bf := *se.bf

	//keyIsStr := se.em.keyType.Kind() == reflect.String
	ct := 0

	bf = append(bf, '{')
	//mpIter2 := reflect.ValueOf(reflect.ValueOf(reflect.ValueOf(reflect.ValueOf(se.srcPtr).Elem().Interface()).Interface())).MapRange()

	for mpIter.Next() {
		ct++
		bf = append(bf, '"')
		key := mpIter.Key().Interface()
		bf = se.em.keyEnc(bf, (*rt.AFace)(unsafe.Pointer(&key)).DataPtr)
		bf = append(bf, "\":"...)

		val := mpIter.Value().Interface()
		ptr := (*rt.AFace)(unsafe.Pointer(&val)).DataPtr
		ptrCt := se.em.ptrLevel
		if ptrCt == 0 {
			goto encMapValue
		}

	peelPtr:
		ptr = *(*unsafe.Pointer)(ptr)
		if ptr == nil {
			bf = append(bf, "null,"...)
			continue
		}
		ptrCt--
		if ptrCt > 0 {
			goto peelPtr
		}

	encMapValue:
		bf = se.em.itemEnc(bf, ptr, se.em.itemType)
	}

	if ct > 0 {
		bf = bf[:len(bf)-1]
	}
	bf = append(bf, '}')
	*se.bf = bf
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// map type value
//func (se *subEncode) encMapKV() {
//	bf := *se.bf
//
//	//itra := reflect.MakeMap(se.em.itemType).MapRange()
//	//itra.Key()
//
//	theMap := *(*map[string]any)(se.srcPtr)
//	//keyIsStr := se.em.keyType.Kind() == reflect.String
//
//	bf = append(bf, '{')
//	for k, v := range theMap {
//		bf = append(bf, '"')
//		//if keyIsStr {
//		bf = append(bf, k...)
//		//} else {
//		//	bf = se.em.keyEnc(bf, unsafe.Pointer(&k), se.em.keyType)
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
//		bf = se.em.itemEnc(bf, ptr, se.em.itemType)
//	}
//	if len(theMap) > 0 {
//		bf = bf[:len(bf)-1]
//	}
//	bf = append(bf, '}')
//	*se.bf = bf
//}
