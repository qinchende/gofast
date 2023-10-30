package dts

import (
	"github.com/qinchende/gofast/skill/lang"
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

	sqlString string
	sqlBool   bool

	//sqlTime time.Time
)

func (val *SqlSkip) Scan(src any) error {
	return nil
}

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

func (val *sqlString) Scan(src any) error {
	switch s := src.(type) {
	case []byte:
		*val = (sqlString)(lang.BTS(s))
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (fa *fieldAttr) intScanner(oPtr unsafe.Pointer) any {
	return (*sqlInt)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int8Scanner(oPtr unsafe.Pointer) any {
	return (*sqlInt8)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int16Scanner(oPtr unsafe.Pointer) any {
	return (*sqlInt16)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int32Scanner(oPtr unsafe.Pointer) any {
	return (*sqlInt32)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) int64Scanner(oPtr unsafe.Pointer) any {
	return (*sqlInt64)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) stringScanner(oPtr unsafe.Pointer) any {
	return (*sqlString)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}

func (fa *fieldAttr) timeScanner(oPtr unsafe.Pointer) any {
	return (*time.Time)(unsafe.Pointer(uintptr(oPtr) + fa.Offset))
}
