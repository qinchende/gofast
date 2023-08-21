// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package dts

import "github.com/qinchende/gofast/cst"

type BindOptions struct {
	FieldTag     string // 解析字段名对应的Tag标签
	ValidTag     string // 验证合法性对应的Tag标签
	CacheSchema  bool   // 是否缓存schema，提高性能
	UseFieldName bool   // 直接用字段名称，而不是通过 reflect tag 取名称
	UseDefValue  bool   // 默认不使用默认值
	UseValid     bool   // 默认不验证值规范

	model int8
}

// 内置几种典型的数据解析模式，当然可以根据需要自定义
var (
	// 应用在大量解析数据记录的场景
	dbStructOptions = &BindOptions{
		model: AsDB,

		FieldTag:     cst.FieldTag,
		ValidTag:     cst.FieldValidTag,
		CacheSchema:  true,
		UseFieldName: false,
		UseDefValue:  true,
		UseValid:     false,
	}

	// 应用在解析配置文件的场景
	reqStructOptions = &BindOptions{
		model: AsReq,

		FieldTag:     cst.FieldTag,
		ValidTag:     cst.FieldValidTag,
		CacheSchema:  true,
		UseFieldName: false,
		UseDefValue:  true,
		UseValid:     true,
	}

	// 应用在解析配置文件的场景
	cfgStructOptions = &BindOptions{
		model: AsConfig,

		FieldTag:     cst.FieldTag,
		ValidTag:     cst.FieldValidTag,
		CacheSchema:  true,
		UseFieldName: true,
		UseDefValue:  true,
		UseValid:     true,
	}
)

const (
	AsConfig int8 = iota // 采用解析配置文件的模式
	AsReq                // 采用解析输入表单的模式
	AsDB                 // 采用解析MySQL记录的模式
)

// 使用什么典型模式来解析验证数据
func AsOptions(model int8) (opt *BindOptions) {
	switch model {
	case AsDB:
		opt = dbStructOptions
	case AsReq:
		opt = reqStructOptions
	case AsConfig:
		opt = cfgStructOptions
	default:
		opt = cfgStructOptions
	}
	return
}
