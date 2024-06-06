// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package lang

import (
	"errors"
	"sort"
)

var (
	ErrInvalidStartPosition = errors.New("start position is invalid")
	ErrInvalidStopPosition  = errors.New("stop position is invalid")
)

func SortByLen(keys []string) {
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) < len(keys[j])
	})
}

// list必须是按字符串长度从小到大排序好的数组，而且不能有空字符串，数据量不可太大
// 匹配到就返回索引，没找到就返回-1
//

func SearchSorted(items []string, str string) int {
	for i := 0; i < len(items); i++ {
		item := items[i]

		if len(item) < len(str) {
			continue
		}
		if len(item) > len(str) {
			break
		}

		// 如果首尾字符相同，然后再全量比较
		if item[0] == str[0] {
			if item[len(item)-1] == str[len(str)-1] {
				if item == str {
					return i
				}
			}
		}
	}
	return -1
}

// 直接跳过前面不可能匹配的项目，加快检索速度
//

func SearchSortedSkip(items []string, step int, str string) int {
	for ; step < len(items); step++ {
		item := items[step]

		if item[0] == str[0] {
			if item[len(item)-1] == str[len(str)-1] {
				if item == str {
					return step
				}
			}
		}

		if len(item) > len(str) {
			break
		}
	}
	return -1
}

func Contains(list []string, str string) bool {
	for i := range list {
		if list[i] == str {
			return true
		}
	}
	return false
}

func Filter(s string, filter func(r rune) bool) string {
	var n int
	chars := []rune(s)
	for i, x := range chars {
		if n < i {
			chars[n] = x
		}
		if !filter(x) {
			n++
		}
	}

	return string(chars[:n])
}

func HasEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) == 0 {
			return true
		}
	}

	return false
}

func NotEmpty(args ...string) bool {
	return !HasEmpty(args...)
}

func Remove(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)

	for _, str := range strs {
		var n int
		for _, v := range out {
			if v != str {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}

	return out
}

func Reverse(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

// Substr returns runes between start and stop [start, stop) regardless of the chars are ascii or utf8
func Substr(str string, start int, stop int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", ErrInvalidStartPosition
	}

	if stop < 0 || stop > length {
		return "", ErrInvalidStopPosition
	}

	return string(rs[start:stop]), nil
}

func TakeOne(valid, or string) string {
	if len(valid) > 0 {
		return valid
	} else {
		return or
	}
}

func TakeWithPriority(fns ...func() string) string {
	for _, fn := range fns {
		val := fn()
		if len(val) > 0 {
			return val
		}
	}

	return ""
}

func Union(first, second []string) []string {
	set := make(map[string]PlaceholderType)

	for _, each := range first {
		set[each] = Placeholder
	}
	for _, each := range second {
		set[each] = Placeholder
	}

	merged := make([]string, 0, len(set))
	for k := range set {
		merged = append(merged, k)
	}

	return merged
}
