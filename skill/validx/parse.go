// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package validx

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	fieldOptionError = "field %s valid failed."
)

// 解析字段配置的选项参数
func ParseOptions(sField *reflect.StructField, str string) (*ValidOptions, error) {
	if str == "" {
		return nil, nil
	}

	var vOpts = ValidOptions{Name: sField.Name}
	var err error

	items := strings.Split(str, ",")
	for _, segment := range items {
		// item指的是k=v字符串
		item := strings.TrimSpace(segment)
		switch {
		case item == attrRequired:
			vOpts.Required = true
		default:
			// 解析成 [k,v]
			kv := strings.Split(item, equalToken)
			if len(kv) != 2 {
				return nil, fmt.Errorf(fieldOptionError, sField.Name)
			}

			switch {
			case kv[0] == attrRequired:
				if vOpts.Required, err = strconv.ParseBool(kv[1]); err != nil {
					return nil, fmt.Errorf(fieldOptionError, sField.Name)
				}
			case kv[0] == attrEnum:
				vOpts.Enum = strings.Split(kv[1], itemSeparator)
			case kv[0] == attrDefault:
				vOpts.DefValue = strings.TrimSpace(kv[1])
				//vOpts.DefExist = true
			case kv[0] == attrRange:
				isDur := false
				if sField.Type == cst.TypeDuration {
					isDur = true
				}
				if vOpts.Range, err = parseNumberRange(kv[1], isDur); err != nil {
					return nil, fmt.Errorf(fieldOptionError, sField.Name)
				}
			case kv[0] == attrLength:
				if vOpts.Len, err = parseNumberRange(kv[1], false); err != nil {
					return nil, fmt.Errorf(fieldOptionError, sField.Name)
				}
			case kv[0] == attrRegex:
				vOpts.Regex = strings.TrimSpace(kv[1])
			case kv[0] == attrMatch:
				vOpts.Match = strings.TrimSpace(kv[1])
			case kv[0] == attrTimeFmt:
				vOpts.TimeFmt = strings.TrimSpace(kv[1])
			}
		}
	}

	return &vOpts, nil
}

// support below notations:
// [:5] (:5] [:5) (:5)
// [1:] [1:) (1:] (1:)
// [1:5] [1:5) (1:5] (1:5)
func parseNumberRange(str string, isDur bool) (ret *numRange, err error) {
	if len(str) == 0 {
		return nil, errNumberRange
	}

	var leftInclude bool
	switch str[0] {
	case '[':
		leftInclude = true
	case '(':
		leftInclude = false
	default:
		return nil, errNumberRange
	}

	str = str[1:]
	if len(str) == 0 {
		return nil, errNumberRange
	}

	var rightInclude bool
	switch str[len(str)-1] {
	case ']':
		rightInclude = true
	case ')':
		rightInclude = false
	default:
		return nil, errNumberRange
	}

	str = str[:len(str)-1]
	fields := strings.Split(str, ":")
	if len(fields) != 2 {
		return nil, errNumberRange
	}

	var left, right float64
	if left, err = parseNumber(fields[0], isDur, -1); err != nil {
		return
	}
	if right, err = parseNumber(fields[1], isDur, 1); err != nil {
		return
	}

	return &numRange{
		min:        left,
		max:        right,
		includeMin: leftInclude,
		includeMax: rightInclude,
	}, nil
}

func parseNumber(str string, isDur bool, abs float64) (val float64, err error) {
	if len(str) > 0 {
		if isDur {
			if dur, errTmp := time.ParseDuration(str); errTmp != nil {
				return 0, errTmp
			} else {
				val = float64(dur)
			}
		} else if val, err = strconv.ParseFloat(str, 64); err != nil {
			return 0, err
		}
	} else {
		val = abs * math.MaxFloat64
	}
	return
}
