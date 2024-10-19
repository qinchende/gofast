package lang

import (
	"encoding/json"
	"math"
	"strconv"
	"unsafe"
)

// 任意类型的值，转换成 Int64，只要能转，不丢失精度都转，否则给出错误
func ToInt64(v any) (i64 int64, err error) {
	if v == nil {
		return 0, errNilValue
	}

	switch vt := v.(type) {
	case string:
		i64, err = strconv.ParseInt(vt, 10, 64)
	case json.Number:
		i64, err = strconv.ParseInt(string(vt), 10, 64)
	case int:
		i64 = int64(vt)
	case uint:
		if vt <= math.MaxInt64 {
			i64 = int64(vt)
		} else {
			err = errNumOutOfRange
		}
	case int64:
		i64 = vt
	case uint64:
		if vt <= math.MaxInt64 {
			i64 = int64(vt)
		} else {
			err = errNumOutOfRange
		}
	case int32:
		i64 = int64(vt)
	case uint32:
		i64 = int64(vt)
	case int16:
		i64 = int64(vt)
	case uint16:
		i64 = int64(vt)
	case int8:
		i64 = int64(vt)
	case uint8:
		i64 = int64(vt)
	case []byte:
		i64, err = strconv.ParseInt(string(vt), 10, 64)
	default:
		err = errConvertValue
	}
	return
}

// 任意类型的值，转换成 Int64，只要能转，不丢失精度都转，否则给出错误
func ToUint64(v any) (ui64 uint64, err error) {
	if v == nil {
		return 0, errNilValue
	}

	switch vt := v.(type) {
	case string:
		ui64, err = strconv.ParseUint(vt, 10, 64)
	case json.Number:
		ui64, err = strconv.ParseUint(string(vt), 10, 64)
	case int:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errNumOutOfRange
		}
	case uint:
		ui64 = uint64(vt)
	case int64:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errNumOutOfRange
		}
	case uint64:
		ui64 = vt
	case int32:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errNumOutOfRange
		}
	case uint32:
		ui64 = uint64(vt)
	case int16:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errNumOutOfRange
		}
	case uint16:
		ui64 = uint64(vt)
	case int8:
		if vt >= 0 {
			ui64 = uint64(vt)
		} else {
			err = errNumOutOfRange
		}
	case uint8:
		ui64 = uint64(vt)
	case []byte:
		ui64, err = strconv.ParseUint(string(vt), 10, 64)
	default:
		err = errConvertValue
	}
	return
}

func ToFloat64(v any) (f64 float64, err error) {
	if v == nil {
		return 0.0, errNilValue
	}

	switch vt := v.(type) {
	case string:
		f64, err = strconv.ParseFloat(vt, 64)
	case json.Number:
		f64, err = strconv.ParseFloat(string(vt), 64)
	case float32:
		f64 = float64(vt)
	case float64:
		f64 = vt
	case int:
		f64, err = strconv.ParseFloat(strconv.FormatInt(int64(vt), 10), 64)
	case uint:
		f64, err = strconv.ParseFloat(strconv.FormatUint(uint64(vt), 10), 64)
	case int64:
		f64, err = strconv.ParseFloat(strconv.FormatInt(vt, 10), 64)
	case uint64:
		f64, err = strconv.ParseFloat(strconv.FormatUint(vt, 10), 64)
	case int32:
		f64 = float64(vt)
	case uint32:
		f64 = float64(vt)
	case int16:
		f64 = float64(vt)
	case uint16:
		f64 = float64(vt)
	case int8:
		f64 = float64(vt)
	case uint8:
		f64 = float64(vt)
	case []byte:
		f64, err = strconv.ParseFloat(string(vt), 64)
	default:
		err = errConvertValue
	}
	return
}

func ToFloat32(v any) (f32 float32, err error) {
	if v2, err2 := ToFloat64(v); err2 != nil {
		return 0.0, err2
	} else if v2 <= math.MaxFloat32 && v2 >= math.SmallestNonzeroFloat32 {
		return float32(v2), nil
	} else {
		return 0.0, errNumOutOfRange
	}
}

func ToInt(v any) (i int, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt && v2 >= math.MinInt {
		return int(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToUint(v any) (ui uint, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint {
		return uint(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToInt32(v any) (i32 int32, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt32 && v2 >= math.MinInt32 {
		return int32(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToUint32(v any) (ui32 uint32, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint32 {
		return uint32(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToInt16(v any) (i16 int16, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt16 && v2 >= math.MinInt16 {
		return int16(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToUint16(v any) (ui16 uint16, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint16 {
		return uint16(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToInt8(v any) (i8 int8, err error) {
	if v2, err2 := ToInt64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxInt8 && v2 >= math.MinInt8 {
		return int8(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

func ToUint8(v any) (ui8 uint8, err error) {
	if v2, err2 := ToUint64(v); err2 != nil {
		return 0, err2
	} else if v2 <= math.MaxUint8 {
		return uint8(v2), nil
	} else {
		return 0, errNumOutOfRange
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//

func ParseInt(s string) int64 {
	if i64, err := strconv.ParseInt(s, 10, 64); err != nil {
		panic(errNumberFmt)
	} else {
		return i64
	}
}

func ParseUint(s string) uint64 {
	if ui64, err := strconv.ParseUint(s, 10, 64); err != nil {
		panic(errNumberFmt)
	} else {
		return ui64
	}
}

func ParseFloat(s string) float64 {
	if f64, err := strconv.ParseFloat(s, 64); err != nil {
		panic(errNumberFmt)
	} else {
		return f64
	}
}

func ParseBool(s string) bool {
	switch s {
	case "1", "t", "T", "True", "true", "TRUE":
		return true
	case "0", "f", "F", "False", "false", "FALSE":
		return false
	default:
		panic(errBoolFmt)
	}
}

// Note: 特殊转换函数，提高性能
// Fast number value parser
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var (
	pow10u64 = [...]uint64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	}
	pow10u64Len = len(pow10u64)
)

// 参数s 必须是严格的uint类型
func ParseUintFast(s string) uint64 {
	maxDigit := len(s)
	if maxDigit > pow10u64Len {
		panic(errNumberFmt)
	}
	sum := uint64(0)
	for i := 0; i < maxDigit; i++ {
		c := uint64(s[i]) - 48
		digitValue := pow10u64[maxDigit-i-1]
		sum += c * digitValue
	}
	return sum
}

var (
	pow10i64 = [...]int64{
		1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
	}
	pow10i64Len = len(pow10i64)
)

// 参数s 必须是严格的int类型
func ParseIntFast(s string) int64 {
	isNegative := false
	if s[0] == '-' {
		s = s[1:]
		isNegative = true
	}
	maxDigit := len(s)
	if maxDigit > pow10i64Len {
		panic(errNumberFmt)
	}
	sum := int64(0)
	for i := 0; i < maxDigit; i++ {
		c := int64(s[i]) - 48
		digitValue := pow10i64[maxDigit-i-1]
		sum += c * digitValue
	}
	if isNegative {
		return -1 * sum
	}
	return sum
}

// fast Int to string
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

var endianness int

func init() {
	var b [2]byte
	*(*uint16)(unsafe.Pointer(&b)) = uint16(0xABCD)

	switch b[0] {
	case 0xCD:
		endianness = 0 // LE
	case 0xAB:
		endianness = 1 // BE
	default:
		panic("could not determine endianness")
	}
}

// "00010203...96979899" cast to []uint16
var intLELookup = [100]uint16{
	0x3030, 0x3130, 0x3230, 0x3330, 0x3430, 0x3530, 0x3630, 0x3730, 0x3830, 0x3930,
	0x3031, 0x3131, 0x3231, 0x3331, 0x3431, 0x3531, 0x3631, 0x3731, 0x3831, 0x3931,
	0x3032, 0x3132, 0x3232, 0x3332, 0x3432, 0x3532, 0x3632, 0x3732, 0x3832, 0x3932,
	0x3033, 0x3133, 0x3233, 0x3333, 0x3433, 0x3533, 0x3633, 0x3733, 0x3833, 0x3933,
	0x3034, 0x3134, 0x3234, 0x3334, 0x3434, 0x3534, 0x3634, 0x3734, 0x3834, 0x3934,
	0x3035, 0x3135, 0x3235, 0x3335, 0x3435, 0x3535, 0x3635, 0x3735, 0x3835, 0x3935,
	0x3036, 0x3136, 0x3236, 0x3336, 0x3436, 0x3536, 0x3636, 0x3736, 0x3836, 0x3936,
	0x3037, 0x3137, 0x3237, 0x3337, 0x3437, 0x3537, 0x3637, 0x3737, 0x3837, 0x3937,
	0x3038, 0x3138, 0x3238, 0x3338, 0x3438, 0x3538, 0x3638, 0x3738, 0x3838, 0x3938,
	0x3039, 0x3139, 0x3239, 0x3339, 0x3439, 0x3539, 0x3639, 0x3739, 0x3839, 0x3939,
}

var intBELookup = [100]uint16{
	0x3030, 0x3031, 0x3032, 0x3033, 0x3034, 0x3035, 0x3036, 0x3037, 0x3038, 0x3039,
	0x3130, 0x3131, 0x3132, 0x3133, 0x3134, 0x3135, 0x3136, 0x3137, 0x3138, 0x3139,
	0x3230, 0x3231, 0x3232, 0x3233, 0x3234, 0x3235, 0x3236, 0x3237, 0x3238, 0x3239,
	0x3330, 0x3331, 0x3332, 0x3333, 0x3334, 0x3335, 0x3336, 0x3337, 0x3338, 0x3339,
	0x3430, 0x3431, 0x3432, 0x3433, 0x3434, 0x3435, 0x3436, 0x3437, 0x3438, 0x3439,
	0x3530, 0x3531, 0x3532, 0x3533, 0x3534, 0x3535, 0x3536, 0x3537, 0x3538, 0x3539,
	0x3630, 0x3631, 0x3632, 0x3633, 0x3634, 0x3635, 0x3636, 0x3637, 0x3638, 0x3639,
	0x3730, 0x3731, 0x3732, 0x3733, 0x3734, 0x3735, 0x3736, 0x3737, 0x3738, 0x3739,
	0x3830, 0x3831, 0x3832, 0x3833, 0x3834, 0x3835, 0x3836, 0x3837, 0x3838, 0x3839,
	0x3930, 0x3931, 0x3932, 0x3933, 0x3934, 0x3935, 0x3936, 0x3937, 0x3938, 0x3939,
}

var intLookup = [2]*[100]uint16{&intLELookup, &intBELookup}

func numMask(numBitSize uint8) uint64 {
	return 1<<numBitSize - 1
}

func AppendInt(out []byte, ptr unsafe.Pointer, numBits uint8) []byte {
	var u64 uint64
	switch numBits {
	case 8:
		u64 = (uint64)(*(*uint8)(ptr))
	case 16:
		u64 = (uint64)(*(*uint16)(ptr))
	case 32:
		u64 = (uint64)(*(*uint32)(ptr))
	case 64:
		u64 = *(*uint64)(ptr)
	}
	mask := numMask(numBits)
	n := u64 & mask
	negative := (u64>>(numBits-1))&1 == 1
	if !negative {
		if n < 10 {
			return append(out, byte(n+'0'))
		} else if n < 100 {
			u := intLELookup[n]
			return append(out, byte(u), byte(u>>8))
		}
	} else {
		n = -n & mask
	}

	lookup := intLookup[endianness]

	var b [22]byte
	u := (*[11]uint16)(unsafe.Pointer(&b))
	i := 11

	for n >= 100 {
		j := n % 100
		n /= 100
		i--
		u[i] = lookup[j]
	}

	i--
	u[i] = lookup[n]

	i *= 2 // convert to byte index
	if n < 10 {
		i++ // remove leading zero
	}
	if negative {
		i--
		b[i] = '-'
	}

	return append(out, b[i:]...)
}

func AppendUint(out []byte, ptr unsafe.Pointer, numBits uint8) []byte {
	var u64 uint64
	switch numBits {
	case 8:
		u64 = (uint64)(*(*uint8)(ptr))
	case 16:
		u64 = (uint64)(*(*uint16)(ptr))
	case 32:
		u64 = (uint64)(*(*uint32)(ptr))
	case 64:
		u64 = *(*uint64)(ptr)
	}
	mask := numMask(numBits)
	n := u64 & mask
	if n < 10 {
		return append(out, byte(n+'0'))
	} else if n < 100 {
		u := intLELookup[n]
		return append(out, byte(u), byte(u>>8))
	}

	lookup := intLookup[endianness]

	var b [22]byte
	u := (*[11]uint16)(unsafe.Pointer(&b))
	i := 11

	for n >= 100 {
		j := n % 100
		n /= 100
		i--
		u[i] = lookup[j]
	}

	i--
	u[i] = lookup[n]

	i *= 2 // convert to byte index
	if n < 10 {
		i++ // remove leading zero
	}
	return append(out, b[i:]...)
}
