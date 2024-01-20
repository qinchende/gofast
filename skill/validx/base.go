// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package validx

const (
	attrRequired = "required" // 必填项
	attrDefault  = "def"
	attrEnum     = "enum"
	attrRange    = "range"
	attrLength   = "len"
	attrRegex    = "regex"
	attrMatch    = "match" // email,mobile,ipv4,ipv4:port,ipv6,id_card,url,file,base64,time
	attrTimeFmt  = "time_fmt"

	// 常用关键字
	itemSeparator = "|"
	equalToken    = "="
)

type (
	ValidOptions struct {
		Name     string    // 字段名称
		Range    *numRange // 数值取值范围
		Enum     []string  // 枚举值数组
		Len      *numRange // 字符串长度范围
		Regex    string    // 正则表达式
		Match    string    // 匹配某个内置的格式
		TimeFmt  string    // 时间格式化
		DefValue string    // 默认值
		Required bool      // 是否必传项
	}

	numRange struct {
		min        float64 // 最小
		max        float64 // 最大
		includeMin bool    // 包括最小
		includeMax bool    // 包括最大
	}
)
