// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mapx

import "github.com/qinchende/gofast/cst"

type ApplyOptions struct {
	FieldTag     string // 解析字段名对应的Tag标签
	ValidTag     string // 验证合法性对应的Tag标签
	CacheSchema  bool   // 是否缓存schema，提高性能
	UseFieldName bool   // 直接用字段名称，而不是通过 reflect tag 取名称
	UseDefValue  bool   // 默认不使用默认值
	UseValid     bool   // 默认不验证值规范
}

var (
	// 应用在大量解析数据记录的场景
	dbStructOptions = &ApplyOptions{
		FieldTag:     cst.FieldTag,
		ValidTag:     cst.FieldValidTag,
		CacheSchema:  true,
		UseFieldName: false,
		UseDefValue:  true,
		UseValid:     false,
	}

	// 应用在解析配置文件的场景
	configStructOptions = &ApplyOptions{
		FieldTag:     cst.FieldTag,
		ValidTag:     cst.FieldValidTag,
		CacheSchema:  false,
		UseFieldName: true,
		UseDefValue:  true,
		UseValid:     true,
	}
)

const (
	LikeConfig int8 = iota // 采用解析配置文件的模式
	LikeDB                 // 采用解析数据库的模式
)

// 使用什么典型配置来解析验证数据
func matchOptions(like int8) (ao *ApplyOptions) {
	switch like {
	case LikeDB:
		ao = dbStructOptions
	case LikeConfig:
		ao = configStructOptions
	default:
		ao = configStructOptions
	}
	return
}
