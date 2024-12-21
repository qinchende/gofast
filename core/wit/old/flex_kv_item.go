// Copyright 2024 GoFast Author(sdx: http://chende.ren). All rights reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// Note: this is copy from stand-lib ~/log/slog
package old

import (
	"time"
)

const badKey = "!BADKEY"

// An KVItem is a key-value pair.
type KVItem struct {
	Key string
	Val TFlex
}

// String returns an KVItem for a string value.
func String(key, value string) KVItem {
	return KVItem{key, StringValue(value)}
}

// Int64 returns an KVItem for an int64.
func Int64(key string, value int64) KVItem {
	return KVItem{key, Int64Value(value)}
}

// Int converts an int to an int64 and returns
// an KVItem with that value.
func Int(key string, value int) KVItem {
	return Int64(key, int64(value))
}

// Uint64 returns an KVItem for a uint64.
func Uint64(key string, v uint64) KVItem {
	return KVItem{key, Uint64Value(v)}
}

// Float64 returns an KVItem for a floating-point number.
func Float64(key string, v float64) KVItem {
	return KVItem{key, Float64Value(v)}
}

// Bool returns an KVItem for a bool.
func Bool(key string, v bool) KVItem {
	return KVItem{key, BoolValue(v)}
}

// Time returns an KVItem for a [time.Time].
// It discards the monotonic portion.
func Time(key string, v time.Time) KVItem {
	return KVItem{key, TimeValue(v)}
}

// Duration returns an KVItem for a [time.Duration].
func Duration(key string, v time.Duration) KVItem {
	return KVItem{key, DurationValue(v)}
}

// Group returns an KVItem for a Group [TFlex].
// The first argument is the key; the remaining arguments
// are converted to Attrs as in [Logger.Log].
//
// Use Group to collect several key-value pairs under a single
// key on a log line, or as the result of LogValue
// in order to log a single value as multiple Attrs.
func Group(key string, args ...any) KVItem {
	return KVItem{key, GroupValue(argsToAttrSlice(args)...)}
}

func argsToAttrSlice(args []any) []KVItem {
	var (
		attr  KVItem
		attrs []KVItem
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

// Any returns an KVItem for the supplied value.
// See [AnyValue] for how values are treated.
func Any(key string, value any) KVItem {
	return KVItem{key, AnyValue(value)}
}

// Equal reports whether a and b have equal keys and values.
func (a KVItem) Equal(b KVItem) bool {
	return a.Key == b.Key && a.Val.Equal(b.Val)
}

func (a KVItem) String() string {
	return a.Key + "=" + a.Val.String()
}

// isEmpty reports whether a has an empty key and a nil value.
// That can be written as KVItem{} or Any("", nil).
func (a KVItem) isEmpty() bool {
	return a.Key == "" && a.Val.num == 0 && a.Val.mate == nil
}

// argsToAttr turns a prefix of the nonempty args slice into an KVItem
// and returns the unconsumed portion of the slice.
// If args[0] is an KVItem, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToAttr(args []any) (KVItem, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return String(badKey, x), nil
		}
		return Any(x, args[1]), args[2:]

	case KVItem:
		return x, args[1:]

	default:
		return Any(badKey, x), args[1:]
	}
}
