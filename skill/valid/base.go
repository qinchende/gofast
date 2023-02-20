package valid

import "reflect"

const (
	attrRequired = "required" // 必填项
	attrDefault  = "def"
	attrEnum     = "enum"
	attrRange    = "range"
	attrLength   = "len"
	attrRegex    = "regex"
	attrMatch    = "match" // email,mobile,ipv4,ipv4:port,ipv6,id_card,url,file,base64,time,datetime

	// 常用关键字
	itemSeparator = "|"
	equalToken    = "="
)

type (
	FieldOpts struct {
		Range    *numRange // 数值取值范围
		Enum     []string  // 枚举值数组
		Len      *numRange // 字符串长度范围
		Regex    string    // 正则表达式
		Match    string    // 匹配某个内置的格式
		DefValue string    // 默认值
		Required bool      // 是否必传项

		SField *reflect.StructField // 原始值，方便后期自定义验证特殊Tag
	}

	numRange struct {
		min        float64 // 最小
		max        float64 // 最大
		includeMin bool    // 包括最小
		includeMax bool    // 包括最大
	}
)
