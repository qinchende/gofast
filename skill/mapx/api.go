package mapx

import (
	"github.com/qinchende/gofast/cst"
)

type ApplyOptions struct {
	FieldTag    string // 解析字段名对应的Tag标签
	ValidTag    string // 验证合法性对应的Tag标签
	CacheSchema bool   // 是否缓存schema，提高性能
	FieldDirect bool   // 忽略FieldTag，直接取字段名称
	NotSnake    bool   // 默认转换成snake模式
	NotDefValue bool   // 默认使用默认值
	NotValid    bool   // 默认解析后就验证
}

var (
	// 应用在大量解析数据记录的场景
	dataOptions = &ApplyOptions{
		FieldTag:    cst.FieldTag,
		ValidTag:    cst.FieldValidTag,
		CacheSchema: true,
		FieldDirect: false,
		NotSnake:    false,
		NotDefValue: false,
		NotValid:    true,
	}

	// 应用在解析配置文件的场景
	configOptions = &ApplyOptions{
		FieldTag:    cst.FieldTag,
		ValidTag:    cst.FieldValidTag,
		CacheSchema: false,
		FieldDirect: true,
		NotSnake:    false,
		NotDefValue: false,
		NotValid:    false,
	}
)

// cst.KV
func ApplyKV(dst any, kvs cst.KV, opts *ApplyOptions) error {
	return applyKVToStruct(dst, kvs, opts)
}

func ApplyKVOfConfig(dst any, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, configOptions)
}

func ApplyKVOfData(dst any, kvs cst.KV) error {
	return applyKVToStruct(dst, kvs, dataOptions)
}

func ApplySliceOfConfig(dst any, src any) error {
	return applyList(dst, src, nil, configOptions)
}

func ApplySliceOfData(dst any, src any) error {
	return applyList(dst, src, nil, dataOptions)
}
