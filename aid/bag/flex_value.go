// Copyright 2024 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// Note: this is copy from stand-lib ~/log/slog
package bag

import (
	"fmt"
	"math"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// A TFlex can represent any Go value, but unlike type any,
// it can represent most small values without an allocation.
// The zero TFlex corresponds to nil.
type TFlex struct {
	_    [0]func() // disallow ==
	num  uint64    // 多出这个64字节，存放简单类型的值
	mate any       // 综合字段
}

type (
	stringptr *byte   // used in TFlex.mate when the TFlex is a string
	groupptr  *KVItem // used in TFlex.mate when the TFlex is a []KVItem
)

// Kind is the kind of a [TFlex].
type Kind int

// The following list is sorted alphabetically, but it's also important that
// KindAny is 0 so that a zero TFlex represents nil.

const (
	KindAny Kind = iota
	KindBool
	KindDuration
	KindF64
	KindI64
	KindStr
	KindTime
	KindU64
	KindGroup
	KindLogValuer
)

var kindStrings = []string{
	"Any",
	"Bool",
	"Duration",
	"Float64",
	"Int64",
	"String",
	"Time",
	"Uint64",
	"Group",
	"LogValuer",
}

func (k Kind) String() string {
	if k >= 0 && int(k) < len(kindStrings) {
		return kindStrings[k]
	}
	return "<unknown slog.Kind>"
}

// Unexported version of Kind, just so we can store Kinds in Values.
// (No user-provided value has this type.)
type kind Kind

// Kind returns v's Kind.
func (v TFlex) Kind() Kind {
	switch x := v.mate.(type) {
	case Kind:
		return x
	case stringptr:
		return KindStr
	case timeLocation, timeTime:
		return KindTime
	case groupptr:
		return KindGroup
	case LogValuer:
		return KindLogValuer
	case kind: // a kind is just a wrapper for a Kind
		return KindAny
	default:
		return KindAny
	}
}

//////////////// Constructors

// StringValue returns a new [TFlex] for a string.
func StringValue(value string) TFlex {
	return TFlex{num: uint64(len(value)), mate: stringptr(unsafe.StringData(value))}
}

// IntValue returns a [TFlex] for an int.
func IntValue(v int) TFlex {
	return Int64Value(int64(v))
}

// Int64Value returns a [TFlex] for an int64.
func Int64Value(v int64) TFlex {
	return TFlex{num: uint64(v), mate: KindI64}
}

// Uint64Value returns a [TFlex] for a uint64.
func Uint64Value(v uint64) TFlex {
	return TFlex{num: v, mate: KindU64}
}

// Float64Value returns a [TFlex] for a floating-point number.
func Float64Value(v float64) TFlex {
	return TFlex{num: math.Float64bits(v), mate: KindF64}
}

// BoolValue returns a [TFlex] for a bool.
func BoolValue(v bool) TFlex {
	u := uint64(0)
	if v {
		u = 1
	}
	return TFlex{num: u, mate: KindBool}
}

type (
	// Unexported version of *time.Location, just so we can store *time.Locations in
	// Values. (No user-provided value has this type.)
	timeLocation *time.Location

	// timeTime is for times where UnixNano is undefined.
	timeTime time.Time
)

// TimeValue returns a [TFlex] for a [time.Time].
// It discards the monotonic portion.
func TimeValue(v time.Time) TFlex {
	if v.IsZero() {
		// UnixNano on the zero time is undefined, so represent the zero time
		// with a nil *time.Location instead. time.Time.Location method never
		// returns nil, so a TFlex with any == timeLocation(nil) cannot be
		// mistaken for any other TFlex, time.Time or otherwise.
		return TFlex{mate: timeLocation(nil)}
	}
	nsec := v.UnixNano()
	t := time.Unix(0, nsec)
	if v.Equal(t) {
		// UnixNano correctly represents the time, so use a zero-alloc representation.
		return TFlex{num: uint64(nsec), mate: timeLocation(v.Location())}
	}
	// Fall back to the general form.
	// Strip the monotonic portion to match the other representation.
	return TFlex{mate: timeTime(v.Round(0))}
}

// DurationValue returns a [TFlex] for a [time.Duration].
func DurationValue(v time.Duration) TFlex {
	return TFlex{num: uint64(v.Nanoseconds()), mate: KindDuration}
}

// GroupValue returns a new [TFlex] for a list of Attrs.
// The caller must not subsequently mutate the argument slice.
func GroupValue(as ...KVItem) TFlex {
	// Remove empty groups.
	// It is simpler overall to do this at construction than
	// to check each Group recursively for emptiness.
	if n := countEmptyGroups(as); n > 0 {
		as2 := make([]KVItem, 0, len(as)-n)
		for _, a := range as {
			if !a.Val.isEmptyGroup() {
				as2 = append(as2, a)
			}
		}
		as = as2
	}
	return TFlex{num: uint64(len(as)), mate: groupptr(unsafe.SliceData(as))}
}

// countEmptyGroups returns the number of empty group values in its argument.
func countEmptyGroups(as []KVItem) int {
	n := 0
	for _, a := range as {
		if a.Val.isEmptyGroup() {
			n++
		}
	}
	return n
}

// AnyValue returns a [TFlex] for the supplied value.
//
// If the supplied value is of type TFlex, it is returned
// unmodified.
//
// Given a value of one of Go's predeclared string, bool, or
// (non-complex) numeric types, AnyValue returns a TFlex of kind
// [KindStr], [KindBool], [KindU64], [KindI64], or [KindF64].
// The width of the original numeric type is not preserved.
//
// Given a [time.Time] or [time.Duration] value, AnyValue returns a TFlex of kind
// [KindTime] or [KindDuration]. The monotonic time is not preserved.
//
// For nil, or values of all other types, including named types whose
// underlying type is numeric, AnyValue returns a value of kind [KindAny].
func AnyValue(v any) TFlex {
	switch v := v.(type) {
	case string:
		return StringValue(v)
	case int:
		return Int64Value(int64(v))
	case uint:
		return Uint64Value(uint64(v))
	case int64:
		return Int64Value(v)
	case uint64:
		return Uint64Value(v)
	case bool:
		return BoolValue(v)
	case time.Duration:
		return DurationValue(v)
	case time.Time:
		return TimeValue(v)
	case uint8:
		return Uint64Value(uint64(v))
	case uint16:
		return Uint64Value(uint64(v))
	case uint32:
		return Uint64Value(uint64(v))
	case uintptr:
		return Uint64Value(uint64(v))
	case int8:
		return Int64Value(int64(v))
	case int16:
		return Int64Value(int64(v))
	case int32:
		return Int64Value(int64(v))
	case float64:
		return Float64Value(v)
	case float32:
		return Float64Value(float64(v))
	case []KVItem:
		return GroupValue(v...)
	case Kind:
		return TFlex{mate: kind(v)}
	case TFlex:
		return v
	default:
		return TFlex{mate: v}
	}
}

//////////////// Accessors

// Any returns v's value as an any.
func (v TFlex) Any() any {
	switch v.Kind() {
	case KindAny:
		if k, ok := v.mate.(kind); ok {
			return Kind(k)
		}
		return v.mate
	case KindLogValuer:
		return v.mate
	case KindGroup:
		return v.group()
	case KindI64:
		return int64(v.num)
	case KindU64:
		return v.num
	case KindF64:
		return v.float()
	case KindStr:
		return v.str()
	case KindBool:
		return v.bool()
	case KindDuration:
		return v.duration()
	case KindTime:
		return v.time()
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
}

// String returns TFlex's value as a string, formatted like [fmt.Sprint]. Unlike
// the methods Int64, Float64, and so on, which panic if v is of the
// wrong kind, String never panics.
func (v TFlex) String() string {
	if sp, ok := v.mate.(stringptr); ok {
		return unsafe.String(sp, v.num)
	}
	var buf []byte
	return string(v.append(buf))
}

func (v TFlex) str() string {
	return unsafe.String(v.mate.(stringptr), v.num)
}

// Int64 returns v's value as an int64. It panics
// if v is not a signed integer.
func (v TFlex) Int64() int64 {
	if g, w := v.Kind(), KindI64; g != w {
		panic(fmt.Sprintf("TFlex kind is %s, not %s", g, w))
	}
	return int64(v.num)
}

// Uint64 returns v's value as a uint64. It panics
// if v is not an unsigned integer.
func (v TFlex) Uint64() uint64 {
	if g, w := v.Kind(), KindU64; g != w {
		panic(fmt.Sprintf("TFlex kind is %s, not %s", g, w))
	}
	return v.num
}

// Bool returns v's value as a bool. It panics
// if v is not a bool.
func (v TFlex) Bool() bool {
	if g, w := v.Kind(), KindBool; g != w {
		panic(fmt.Sprintf("TFlex kind is %s, not %s", g, w))
	}
	return v.bool()
}

func (v TFlex) bool() bool {
	return v.num == 1
}

// Duration returns v's value as a [time.Duration]. It panics
// if v is not a time.Duration.
func (v TFlex) Duration() time.Duration {
	if g, w := v.Kind(), KindDuration; g != w {
		panic(fmt.Sprintf("TFlex kind is %s, not %s", g, w))
	}

	return v.duration()
}

func (v TFlex) duration() time.Duration {
	return time.Duration(int64(v.num))
}

// Float64 returns v's value as a float64. It panics
// if v is not a float64.
func (v TFlex) Float64() float64 {
	if g, w := v.Kind(), KindF64; g != w {
		panic(fmt.Sprintf("TFlex kind is %s, not %s", g, w))
	}

	return v.float()
}

func (v TFlex) float() float64 {
	return math.Float64frombits(v.num)
}

// Time returns v's value as a [time.Time]. It panics
// if v is not a time.Time.
func (v TFlex) Time() time.Time {
	if g, w := v.Kind(), KindTime; g != w {
		panic(fmt.Sprintf("TFlex kind is %s, not %s", g, w))
	}
	return v.time()
}

// See TimeValue to understand how times are represented.
func (v TFlex) time() time.Time {
	switch a := v.mate.(type) {
	case timeLocation:
		if a == nil {
			return time.Time{}
		}
		return time.Unix(0, int64(v.num)).In(a)
	case timeTime:
		return time.Time(a)
	default:
		panic(fmt.Sprintf("bad time type %T", v.mate))
	}
}

// LogValuer returns v's value as a LogValuer. It panics
// if v is not a LogValuer.
func (v TFlex) LogValuer() LogValuer {
	return v.mate.(LogValuer)
}

// Group returns v's value as a []KVItem.
// It panics if v's [Kind] is not [KindGroup].
func (v TFlex) Group() []KVItem {
	if sp, ok := v.mate.(groupptr); ok {
		return unsafe.Slice((*KVItem)(sp), v.num)
	}
	panic("Group: bad kind")
}

func (v TFlex) group() []KVItem {
	return unsafe.Slice((*KVItem)(v.mate.(groupptr)), v.num)
}

//////////////// Other

// Equal reports whether v and w represent the same Go value.
func (v TFlex) Equal(w TFlex) bool {
	k1 := v.Kind()
	k2 := w.Kind()
	if k1 != k2 {
		return false
	}
	switch k1 {
	case KindI64, KindU64, KindBool, KindDuration:
		return v.num == w.num
	case KindStr:
		return v.str() == w.str()
	case KindF64:
		return v.float() == w.float()
	case KindTime:
		return v.time().Equal(w.time())
	case KindAny, KindLogValuer:
		return v.mate == w.mate // may panic if non-comparable
	case KindGroup:
		return slices.EqualFunc(v.group(), w.group(), KVItem.Equal)
	default:
		panic(fmt.Sprintf("bad kind: %s", k1))
	}
}

// isEmptyGroup reports whether v is a group that has no attributes.
func (v TFlex) isEmptyGroup() bool {
	if v.Kind() != KindGroup {
		return false
	}
	// We do not need to recursively examine the group's Attrs for emptiness,
	// because GroupValue removed them when the group was constructed, and
	// groups are immutable.
	return len(v.group()) == 0
}

// append appends a text representation of v to dst.
// v is formatted as with fmt.Sprint.
func (v TFlex) append(dst []byte) []byte {
	switch v.Kind() {
	case KindStr:
		return append(dst, v.str()...)
	case KindI64:
		return strconv.AppendInt(dst, int64(v.num), 10)
	case KindU64:
		return strconv.AppendUint(dst, v.num, 10)
	case KindF64:
		return strconv.AppendFloat(dst, v.float(), 'g', -1, 64)
	case KindBool:
		return strconv.AppendBool(dst, v.bool())
	case KindDuration:
		return append(dst, v.duration().String()...)
	case KindTime:
		return append(dst, v.time().String()...)
	case KindGroup:
		return fmt.Append(dst, v.group())
	case KindAny, KindLogValuer:
		return fmt.Append(dst, v.mate)
	default:
		panic(fmt.Sprintf("bad kind: %s", v.Kind()))
	}
}

// A LogValuer is any Go value that can convert itself into a TFlex for logging.
//
// This mechanism may be used to defer expensive operations until they are
// needed, or to expand a single value into a sequence of components.
type LogValuer interface {
	LogValue() TFlex
}

const maxLogValues = 100

// Resolve repeatedly calls LogValue on v while it implements [LogValuer],
// and returns the result.
// If v resolves to a group, the group's attributes' values are not recursively
// resolved.
// If the number of LogValue calls exceeds a threshold, a TFlex containing an
// error is returned.
// Resolve's return value is guaranteed not to be of Kind [KindLogValuer].
func (v TFlex) Resolve() (rv TFlex) {
	orig := v
	defer func() {
		if r := recover(); r != nil {
			rv = AnyValue(fmt.Errorf("LogValue panicked\n%s", stack(3, 5)))
		}
	}()

	for i := 0; i < maxLogValues; i++ {
		if v.Kind() != KindLogValuer {
			return v
		}
		v = v.LogValuer().LogValue()
	}
	err := fmt.Errorf("LogValue called too many times on TFlex of type %T", orig.Any())
	return AnyValue(err)
}

func stack(skip, nFrames int) string {
	pcs := make([]uintptr, nFrames+1)
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return "(no stack)"
	}
	frames := runtime.CallersFrames(pcs[:n])
	var b strings.Builder
	i := 0
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&b, "called from %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
		i++
		if i >= nFrames {
			fmt.Fprintf(&b, "(rest of stack elided)\n")
			break
		}
	}
	return b.String()
}
