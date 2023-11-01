package dts

import (
	"fmt"
	"github.com/qinchende/gofast/skill/lang"
	"strconv"
	"time"
	"unsafe"
)

type (
	SqlSkip int

	sqlInt   int
	sqlInt8  int8
	sqlInt16 int16
	sqlInt32 int32
	sqlInt64 int64

	sqlUint   uint
	sqlUint8  uint8
	sqlUint16 uint16
	sqlUint32 uint32
	sqlUint64 uint64

	sqlFloat32 float32
	sqlFloat64 float64

	//sqlString string
	sqlBool   bool

	//sqlAny  any
	//sqlTime time.Time
)

func (val *SqlSkip) Scan(src any) error {
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (val *sqlInt) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlInt(s)
	}
	return nil
}

func (val *sqlInt8) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlInt8(s)
	}
	return nil
}

func (val *sqlInt16) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlInt16(s)
	}
	return nil
}

func (val *sqlInt32) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlInt32(s)
	}
	return nil
}

func (val *sqlInt64) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlInt64(s)
	}
	return nil
}

// ++++++++++++++++

func (val *sqlUint) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlUint(s)
	}
	return nil
}

func (val *sqlUint8) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlUint8(s)
	}
	return nil
}

func (val *sqlUint16) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlUint16(s)
	}
	return nil
}

func (val *sqlUint32) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlUint32(s)
	}
	return nil
}

func (val *sqlUint64) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = sqlUint64(s)
	}
	return nil
}

// ++++++++++++++++

func (val *sqlFloat32) Scan(src any) error {
	switch s := src.(type) {
	case float64:
		*val = sqlFloat32(s)
	}
	return nil
}

func (val *sqlFloat64) Scan(src any) error {
	switch s := src.(type) {
	case float64:
		*val = sqlFloat64(s)
	}
	return nil
}

// ++++++++++++++++

//func (val *sqlString) Scan(src any) error {
//	switch s := src.(type) {
//	case []byte:
//		*val = (sqlString)(lang.BTS(s))
//	}
//	return nil
//}

func (val *sqlBool) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		if s == 1 || s == 0 {
			*val = (s == 1)
		} else {
			return fmt.Errorf("dts: couldn't convert %d into type bool", s)
		}
	case []byte:
		bv, err := strconv.ParseBool(lang.BTS(s))
		if err == nil {
			*val = sqlBool(bv)
		}
		return err
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type bool", src, src)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (fa *fieldAttr) intValue(oPtr unsafe.Pointer) any {
	return (*sqlInt)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int8Value(oPtr unsafe.Pointer) any {
	return (*sqlInt8)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int16Value(oPtr unsafe.Pointer) any {
	return (*sqlInt16)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int32Value(oPtr unsafe.Pointer) any {
	return (*sqlInt32)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int64Value(oPtr unsafe.Pointer) any {
	return (*sqlInt64)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) uintValue(oPtr unsafe.Pointer) any {
	return (*sqlUint)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) uint8Value(oPtr unsafe.Pointer) any {
	return (*sqlUint8)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) uint16Value(oPtr unsafe.Pointer) any {
	return (*sqlUint16)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) uint32Value(oPtr unsafe.Pointer) any {
	return (*sqlUint32)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) uint64Value(oPtr unsafe.Pointer) any {
	return (*sqlUint64)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) float32Value(oPtr unsafe.Pointer) any {
	return (*sqlFloat32)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) float64Value(oPtr unsafe.Pointer) any {
	return (*sqlFloat64)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

// ++++++++++++++
func (fa *fieldAttr) stringValue(oPtr unsafe.Pointer) any {
	return (*string)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) boolValue(oPtr unsafe.Pointer) any {
	return (*sqlBool)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (fa *fieldAttr) anyValue(oPtr unsafe.Pointer) any {
	return (*any)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) timeValue(oPtr unsafe.Pointer) any {
	return (*time.Time)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}
