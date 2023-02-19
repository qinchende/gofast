package cst

// 整个项目建议统一对象的解析标准，尽量不出现json,xml,form等一系列标签
const (
	FieldTag      = "pms" // 字段名称tag(首要)
	FieldTagDB    = "dbf" // 字段名称tag(次优先级，主要用于DB表结构)
	FieldValidTag = "v"   // 验证字段tag
)

//var FieldTagOthers = []string{"json", "form", "xml"}
