package gson

// 编码示例：
// [1,1,["name","account","age","login","mobile","tok"],[["b{m}c","bmcrmb",37,true,"1344466338783","t:Q0J44CM3R4dHhqWDZZM2944FbTZr"]]]
// Note：
// 1. 中间不含空格，其它特征符合JSON规范
// 2. 第一项是当前包含记录数，第二项是记录总数（用于分页），第三项是数据字段名，第四项是二维数组，存放每条记录对应前面字段的数据。
type GsonRows struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]any
}

// 用来Decoder ++++++++++++++++
// 当解码一个复合数据时，可以支持Gson片段
type RowsDecPet struct {
	Target any
	Ct     int64
	Tt     int64
}

// 解析GsonRows的返回结果
type RowsDecRet struct {
	Err  error
	Ct   int64
	Tt   int64
	Scan int
}

// 用来Encode ++++++++++++++++
type RowsEncPet struct {
	List      any     // 对象列表
	Tt        int64   // 无分页的总数，分页使用
	FieldsStr string  // GsonRows 字段 用 逗号拼接好的字符串
	FieldsIdx []uint8 // GsonRows 数据中 Fields 对应在 struct 的索引
}
