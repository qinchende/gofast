package valid

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

var (
	fieldOptionError = "field %s has wrong valid setting"
)

// 解析字段配置的选项参数
func ParseOptions(field *reflect.StructField, str string) (*FieldOpts, error) {
	if str == "" {
		return nil, nil
	}

	items := strings.Split(str, ",")
	var fOpts FieldOpts
	var err error
	for _, segment := range items {
		item := strings.TrimSpace(segment)
		switch {
		case item == attrRequired:
			fOpts.Required = true
		default:
			kv := strings.Split(item, equalToken)
			if len(kv) != 2 {
				return nil, fmt.Errorf(fieldOptionError, field.Name)
			}

			switch {
			case kv[0] == attrRequired:
				if fOpts.Required, err = strconv.ParseBool(kv[1]); err != nil {
					return nil, fmt.Errorf(fieldOptionError, field.Name)
				}
			case kv[0] == attrEnum:
				fOpts.Enum = strings.Split(kv[1], itemSeparator)
			case kv[0] == attrDefault:
				fOpts.DefValue = strings.TrimSpace(kv[1])
				//fOpts.DefExist = true
			case kv[0] == attrRange:
				if fOpts.Range, err = parseNumberRange(kv[1]); err != nil {
					return nil, fmt.Errorf(fieldOptionError, field.Name)
				}
			case kv[0] == attrLength:
				if fOpts.Len, err = parseNumberRange(kv[1]); err != nil {
					return nil, fmt.Errorf(fieldOptionError, field.Name)
				}
			case kv[0] == attrRegex:
				fOpts.Regex = strings.TrimSpace(kv[1])
			case kv[0] == attrMatch:
				fOpts.Match = strings.TrimSpace(kv[1])
			}
		}
	}

	return &fOpts, nil
}

// support below notations:
// [:5] (:5] [:5) (:5)
// [1:] [1:) (1:] (1:)
// [1:5] [1:5) (1:5] (1:5)
func parseNumberRange(str string) (*numRange, error) {
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

	if len(fields[0]) == 0 && len(fields[1]) == 0 {
		return nil, errNumberRange
	}

	var left float64
	if len(fields[0]) > 0 {
		var err error
		if left, err = strconv.ParseFloat(fields[0], 64); err != nil {
			return nil, err
		}
	} else {
		left = -math.MaxFloat64
	}

	var right float64
	if len(fields[1]) > 0 {
		var err error
		if right, err = strconv.ParseFloat(fields[1], 64); err != nil {
			return nil, err
		}
	} else {
		right = math.MaxFloat64
	}

	return &numRange{
		left:     left,
		right:    right,
		lInclude: leftInclude,
		rInclude: rightInclude,
	}, nil
}
