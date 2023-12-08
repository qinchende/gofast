package jde

import (
	"unsafe"
)

func EncodeGsonRowOnlyValuePart(bf *[]byte, values []any) {
	*bf = append(*bf, '[')
	for _, val := range values {
		encAny(bf, unsafe.Pointer(&val), nil)
	}
	*bf = append(*bf, `],`...)
}
