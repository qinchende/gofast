package dts

import (
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"unsafe"
)

func BindBaseValueAsConfig(kd reflect.Kind, itPtr unsafe.Pointer, itVal any) {
	switch kd {
	case reflect.Int:
		setInt(itPtr, itVal)
	case reflect.Int8:
		setInt8(itPtr, itVal)
	case reflect.Int16:
		setInt16(itPtr, itVal)
	case reflect.Int32:
		setInt32(itPtr, itVal)
	case reflect.Int64:
		setInt64(itPtr, itVal)

	case reflect.Uint:
		setUint(itPtr, itVal)
	case reflect.Uint8:
		setUint8(itPtr, itVal)
	case reflect.Uint16:
		setUint16(itPtr, itVal)
	case reflect.Uint32:
		setUint32(itPtr, itVal)
	case reflect.Uint64:
		setUint64(itPtr, itVal)

	case reflect.Float32:
		setFloat32(itPtr, itVal)
	case reflect.Float64:
		setFloat64(itPtr, itVal)

	case reflect.String:
		setString(itPtr, itVal)
	case reflect.Bool:
		setBool(itPtr, itVal)
	case reflect.Interface:
		setAny(itPtr, itVal)

	}
}

// NOTE：不通用
// 下面的绑定函数，只针对类似Web请求提交的数据。只支持string|number|bool等基础知识，或者 KV|Array
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// int
func setInt(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt(p, v)
	case string:
		BindInt(p, lang.ParseInt(v))
	case *string:
		BindInt(p, lang.ParseInt(*v))
	}
}

func setInt8(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt8(p, v)
	case string:
		BindInt8(p, lang.ParseInt(v))
	case *string:
		BindInt8(p, lang.ParseInt(*v))
	}
}

func setInt16(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt16(p, v)
	case string:
		BindInt16(p, lang.ParseInt(v))
	case *string:
		BindInt16(p, lang.ParseInt(*v))
	}
}

func setInt32(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt32(p, v)
	case string:
		BindInt32(p, lang.ParseInt(v))
	case *string:
		BindInt32(p, lang.ParseInt(*v))
	}
}

func setInt64(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case int64:
		BindInt64(p, v)
	case string:
		BindInt64(p, lang.ParseInt(v))
	case *string:
		BindInt32(p, lang.ParseInt(*v))
	}
}

// uint
func setUint(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint(p, v)
	case string:
		BindUint(p, lang.ParseUint(v))
	case *string:
		BindUint(p, lang.ParseUint(*v))
	}
}

func setUint8(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint8(p, v)
	case string:
		BindUint8(p, lang.ParseUint(v))
	case *string:
		BindUint8(p, lang.ParseUint(*v))
	}
}

func setUint16(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint16(p, v)
	case string:
		BindUint16(p, lang.ParseUint(v))
	case *string:
		BindUint16(p, lang.ParseUint(*v))
	}
}

func setUint32(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint32(p, v)
	case string:
		BindUint32(p, lang.ParseUint(v))
	case *string:
		BindUint32(p, lang.ParseUint(*v))
	}
}

func setUint64(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case uint64:
		BindUint64(p, v)
	case string:
		BindUint64(p, lang.ParseUint(v))
	case *string:
		BindUint64(p, lang.ParseUint(*v))
	}
}

// float
func setFloat32(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case float64:
		BindFloat32(p, v)
	case string:
		BindFloat32(p, lang.ParseFloat(v))
	case *string:
		BindFloat32(p, lang.ParseFloat(*v))
	}
}

func setFloat64(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case float64:
		BindFloat64(p, v)
	case string:
		BindFloat64(p, lang.ParseFloat(v))
	case *string:
		BindFloat64(p, lang.ParseFloat(*v))
	}
}

// string
func setString(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case string:
		BindString(p, v)
	case *string:
		BindString(p, *v)
	default:
		BindString(p, lang.ToString(v))
	}
}

func setBool(p unsafe.Pointer, val any) {
	switch v := val.(type) {
	case bool:
		BindBool(p, v)
	case string:
		BindBool(p, lang.ParseBool(v))
	case *string:
		BindBool(p, lang.ParseBool(*v))
	}
}

func setAny(p unsafe.Pointer, val any) {
	BindAny(p, val)
}
