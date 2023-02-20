package valid

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/skill/lang"
	"reflect"
	"regexp"
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

// 验证结构体字段值，是否符合指定规范
func ValidateField(fValue *reflect.Value, fOpts *FieldOpts) error {
	if fOpts == nil {
		return nil
	}

	var err error
	// 实体对象字段类型
	switch fValue.Kind() {
	case reflect.String:
		str := fValue.String()
		// 字符串长度
		if fOpts.Len != nil {
			if err = checkNumberRange(float64(len(str)), fOpts.Len); err != nil {
				return err
			}
		}
		// 检查是否符合枚举
		if fOpts.Enum != nil && !lang.Contains(fOpts.Enum, str) {
			return fmt.Errorf(`value "%s" not in "%v"`, str, fOpts.Enum)
		}
		// 否则常见的正则表达式
		if fOpts.Match != "" {
			reg := regexMap[fOpts.Match]
			if reg != nil {
				if reg.MatchString(str) == false {
					return fmt.Errorf(`value "%s" not like "%s"`, str, fOpts.Match)
				}
			} else {
				switch {
				//case str == "email":
				//case str == "mobile":
				//case str == "ipv4":
				//case str == "ipv4:port":
				//case str == "base64":
				//case str == "base64URL":
				case str == "ipv6":
				case str == "id_card":
				case str == "url":
				case str == "file":
				case str == "time":
				case str == "datetime":
				}
			}
		}
		// 自定义正则表达式
		if fOpts.Regex != "" {
			if regexp.MustCompile(fOpts.Regex).MatchString(str) == false {
				return fmt.Errorf(`value "%s" can't match "%s"`, str, fOpts.Regex)
			}
		}
	default:
		var f64 float64
		switch fValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f64 = float64(fValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			f64 = float64(fValue.Uint())
		case reflect.Float32, reflect.Float64:
			f64 = fValue.Float()
		}
		err = checkNumberRange(f64, fOpts.Range)
	}
	return err
}

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
