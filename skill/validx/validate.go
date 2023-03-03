// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package validx

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"regexp"
	"time"
)

var (
	errNumberRange = errors.New("wrong number range setting")
)

var regexMap = map[string]*regexp.Regexp{
	"email":     regexp.MustCompile(emailRegexString),
	"mobile":    regexp.MustCompile(chinaMobileRegexString),
	"ipv4":      regexp.MustCompile(ipv4RegexString),
	"ipv4:port": regexp.MustCompile(ipv4PortRegexString),
	"base64":    regexp.MustCompile(base64RegexString),
	"base64URL": regexp.MustCompile(base64URLRegexString),
}

func IsMobile(str string) bool {
	return regexMap["mobile"].MatchString(str)
}

func IsEmail(str string) bool {
	return regexMap["email"].MatchString(str)
}

func IsIPv4(str string) bool {
	return regexMap["ipv4"].MatchString(str)
}

func IsIPv4Port(str string) bool {
	return regexMap["ipv4:port"].MatchString(str)
}

func IsBase64(str string) bool {
	return regexMap["base64"].MatchString(str)
}

func IsLenRange(str string, min, max int) bool {
	return len(str) >= min && len(str) <= max
}

func IsNumRange(num float64, min, max float64) bool {
	return num >= min && num <= max
}

// 验证结构体字段值，是否符合指定规范
func ValidateField(fValue *reflect.Value, vOpts *ValidOptions) (err error) {
	if vOpts == nil {
		return nil
	}

	switch fValue.Kind() {
	case reflect.String:
		str := fValue.String()
		// 字符串长度
		if vOpts.Len != nil {
			if err = checkNumberRange(float64(len(str)), vOpts.Len); err != nil {
				return err
			}
		}
		// 检查是否符合枚举
		if vOpts.Enum != nil && !lang.Contains(vOpts.Enum, str) {
			return fmt.Errorf(`value "%s" not in "%v"`, str, vOpts.Enum)
		}
		// 否则常见的正则表达式
		if len(vOpts.Match) > 0 {
			reg := regexMap[vOpts.Match]
			if reg != nil {
				// 通用正则验证
				// case "email":
				// case "mobile":
				// case "ipv4":
				// case "ipv4:port":
				// case "base64":
				// case "base64URL":
				if reg.MatchString(str) == false {
					return fmt.Errorf(`The value "%s" not like "%s"`, str, vOpts.Match)
				}
			} else {
				// 自定义验证逻辑
				switch vOpts.Match {
				case "ipv6":
				case "id_card":
				case "url":
				case "file":
				case "time":
					if len(vOpts.TimeFmt) > 0 {
						if _, err := time.Parse(vOpts.TimeFmt, str); err != nil {
							return fmt.Errorf(`value "%s" can't match time format "%s"`, str, vOpts.TimeFmt)
						}
					}
				}
			}
		}
		// 自定义正则验证
		if len(vOpts.Regex) > 0 {
			if regexp.MustCompile(vOpts.Regex).MatchString(str) == false {
				return fmt.Errorf(`value "%s" can't match "%s"`, str, vOpts.Regex)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = checkNumberRange(float64(fValue.Int()), vOpts.Range)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		err = checkNumberRange(float64(fValue.Uint()), vOpts.Range)
	case reflect.Float32, reflect.Float64:
		err = checkNumberRange(fValue.Float(), vOpts.Range)
	case reflect.Struct:
		// todo: 如果是 time.Time 类型如何处理
	}
	return
}

// 长度范围验证
func checkNumberRange(fv float64, nr *numRange) error {
	if nr == nil {
		return nil
	}
	if (nr.includeMin && fv < nr.min) || (!nr.includeMin && fv <= nr.min) {
		return errNumberRange
	}
	if (nr.includeMax && fv > nr.max) || (!nr.includeMax && fv >= nr.max) {
		return errNumberRange
	}
	return nil
}
