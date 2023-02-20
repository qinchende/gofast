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
	dbStructOptions = &ApplyOptions{
		FieldTag:    cst.FieldTag,
		ValidTag:    cst.FieldValidTag,
		CacheSchema: true,
		FieldDirect: false,
		NotSnake:    false,
		NotDefValue: false,
		NotValid:    true,
	}

	// 应用在解析配置文件的场景
	configStructOptions = &ApplyOptions{
		FieldTag:    cst.FieldTag,
		ValidTag:    cst.FieldValidTag,
		CacheSchema: false,
		FieldDirect: true,
		NotSnake:    false,
		NotDefValue: false,
		NotValid:    false,
	}
)

const (
	LikeConfig int8 = iota
	LikeDB
)

func ApplyKV(dst any, kvs cst.KV, like int8) error {
	if like == LikeDB {
		return applyKVToStruct(dst, kvs, dbStructOptions)
	}
	return applyKVToStruct(dst, kvs, configStructOptions)
}

func ApplyKVX(dst any, kvs cst.KV, opts *ApplyOptions) error {
	return applyKVToStruct(dst, kvs, opts)
}

func ApplySlice(dst any, src any, like int8) error {
	if like == LikeDB {
		return applyList(dst, src, nil, dbStructOptions)
	}
	return applyList(dst, src, nil, configStructOptions)
}
func ApplySliceX(dst any, src any, opts *ApplyOptions) error {
	return applyList(dst, src, nil, opts)
}

// 根据结构体配置信息，优化字段值 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Optimize(dst any, like int8) error {
	switch like {
	case LikeDB:
		return optimizeStruct(dst, dbStructOptions)
	default:
		return optimizeStruct(dst, configStructOptions)
	}
}
func OptimizeX(dst any, opts *ApplyOptions) error {
	return optimizeStruct(dst, opts)
}
