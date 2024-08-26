package dts

import (
	"fmt"
	"github.com/qinchende/gofast/aid/lang"
	"strconv"
	"time"
	"unsafe"
)

type (
	SqlSkip int

	SqlInt      int
	SqlInt8     int8
	SqlInt16    int16
	SqlInt32    int32
	SqlInt64    int64
	SqlDuration int64

	SqlUint   uint
	SqlUint8  uint8
	SqlUint16 uint16
	SqlUint32 uint32
	SqlUint64 uint64

	SqlFloat32 float32
	SqlFloat64 float64

	SqlBool bool

	//sqlString string
	//sqlAny  any
	//sqlTime time.Time
)

// Note: 下面Scan方法中 src 参数是返回的字段数据，这个值的类型在 go-sql-driver/mysql 解析下只可能是：
// nil | int64 | float32 | float64 | []byte | time.Time 类型

func (val *SqlSkip) Scan(src any) error {
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (val *SqlInt) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlInt(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type int", src, src)
	}
	return nil
}

func (val *SqlInt8) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlInt8(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type int8", src, src)
	}
	return nil
}

func (val *SqlInt16) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlInt16(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type int16", src, src)
	}
	return nil
}

func (val *SqlInt32) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlInt32(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type int32", src, src)
	}
	return nil
}

func (val *SqlInt64) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlInt64(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type int64", src, src)
	}
	return nil
}

func (val *SqlDuration) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlDuration(s)
	case string:
		if d, err := time.ParseDuration(s); err != nil {
			return err
		} else {
			*val = SqlDuration(d)
		}
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type int64", src, src)
	}
	return nil
}

// ++++++++++++++++

func (val *SqlUint) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlUint(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type uint", src, src)
	}
	return nil
}

func (val *SqlUint8) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlUint8(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type uint8", src, src)
	}
	return nil
}

func (val *SqlUint16) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlUint16(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type uint16", src, src)
	}
	return nil
}

func (val *SqlUint32) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlUint32(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type uint32", src, src)
	}
	return nil
}

func (val *SqlUint64) Scan(src any) error {
	switch s := src.(type) {
	case int64:
		*val = SqlUint64(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type uint64", src, src)
	}
	return nil
}

// ++++++++++++++++

func (val *SqlFloat32) Scan(src any) error {
	switch s := src.(type) {
	case float32:
		*val = SqlFloat32(s)
	case float64:
		*val = SqlFloat32(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type float32", src, src)
	}
	return nil
}

func (val *SqlFloat64) Scan(src any) error {
	switch s := src.(type) {
	case float32:
		*val = SqlFloat64(s)
	case float64:
		*val = SqlFloat64(s)
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type float64", src, src)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Note: src共享了底层切片数组，需要copy到新内存
//func (val *sqlString) Scan(src any) error {
//	switch s := src.(type) {
//	case []byte:
//		*val = (sqlString)(lang.BTS(s))
//	}
//	return nil
//}

func (val *SqlBool) Scan(src any) error {
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
			*val = SqlBool(bv)
		}
		return err
	default:
		return fmt.Errorf("dts: couldn't convert %v (%T) into type bool", src, src)
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (fa *fieldAttr) intValue(oPtr unsafe.Pointer) any {
	return (*SqlInt)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) int8Value(oPtr unsafe.Pointer) any {
	return (*SqlInt8)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) int16Value(oPtr unsafe.Pointer) any {
	return (*SqlInt16)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) int32Value(oPtr unsafe.Pointer) any {
	return (*SqlInt32)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) int64Value(oPtr unsafe.Pointer) any {
	return (*SqlInt64)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) durationValue(oPtr unsafe.Pointer) any {
	return (*SqlDuration)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) uintValue(oPtr unsafe.Pointer) any {
	return (*SqlUint)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) uint8Value(oPtr unsafe.Pointer) any {
	return (*SqlUint8)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) uint16Value(oPtr unsafe.Pointer) any {
	return (*SqlUint16)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) uint32Value(oPtr unsafe.Pointer) any {
	return (*SqlUint32)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) uint64Value(oPtr unsafe.Pointer) any {
	return (*SqlUint64)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) float32Value(oPtr unsafe.Pointer) any {
	return (*SqlFloat32)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) float64Value(oPtr unsafe.Pointer) any {
	return (*SqlFloat64)(fa.MyPtr(oPtr))
}

// ++++++++++++++
func (fa *fieldAttr) boolValue(oPtr unsafe.Pointer) any {
	return (*SqlBool)(fa.MyPtr(oPtr))
}

// ++++++++++++++
// Note: 获取字符串切片，无法共享底层字节切片。因为 db.conn 读写数据用到的Buffer可能会被复用，值会被覆盖。
func (fa *fieldAttr) stringValue(oPtr unsafe.Pointer) any {
	return (*string)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) anyValue(oPtr unsafe.Pointer) any {
	return (*any)(fa.MyPtr(oPtr))
}

func (fa *fieldAttr) timeValue(oPtr unsafe.Pointer) any {
	return (*time.Time)(fa.MyPtr(oPtr))
}

// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

func IntValue(ptr unsafe.Pointer) any {
	return (*SqlInt)(ptr)
}

func Int8Value(ptr unsafe.Pointer) any {
	return (*SqlInt8)(ptr)
}

func Int16Value(ptr unsafe.Pointer) any {
	return (*SqlInt16)(ptr)
}

func Int32Value(ptr unsafe.Pointer) any {
	return (*SqlInt32)(ptr)
}

func Int64Value(ptr unsafe.Pointer) any {
	return (*SqlInt64)(ptr)
}

func UintValue(ptr unsafe.Pointer) any {
	return (*SqlUint)(ptr)
}

func Uint8Value(ptr unsafe.Pointer) any {
	return (*SqlUint8)(ptr)
}

func Uint16Value(ptr unsafe.Pointer) any {
	return (*SqlUint16)(ptr)
}

func Uint32Value(ptr unsafe.Pointer) any {
	return (*SqlUint32)(ptr)
}

func Uint64Value(ptr unsafe.Pointer) any {
	return (*SqlUint64)(ptr)
}

func Float32Value(ptr unsafe.Pointer) any {
	return (*SqlFloat32)(ptr)
}

func Float64Value(ptr unsafe.Pointer) any {
	return (*SqlFloat64)(ptr)
}

// ++++++++++++
func BoolValue(ptr unsafe.Pointer) any {
	return (*SqlBool)(ptr)
}

// ++++++++++++
// Note: 获取字符串切片，无法共享底层字节切片。因为 db.conn 读写数据用到的Buffer可能会被复用，值会被覆盖。
func StringValue(ptr unsafe.Pointer) any {
	return (*string)(ptr)
}

func AnyValue(ptr unsafe.Pointer) any {
	return (*any)(ptr)
}

func TimeValue(ptr unsafe.Pointer) any {
	return (*time.Time)(ptr)
}
