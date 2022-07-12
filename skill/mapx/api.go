package mapx

import (
	"github.com/qinchende/gofast/cst"
	"io"
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

func DefOptions(opt *ApplyOptions) *ApplyOptions {
	if opt == nil {
		opt = &ApplyOptions{}
	}
	if opt.FieldTag == "" {
		opt.FieldTag = "pms"
	}
	if opt.ValidTag == "" {
		opt.ValidTag = "valid"
	}
	return opt
}

// 应用在大量解析数据记录的场景
func DataOptions() *ApplyOptions {
	return &ApplyOptions{
		FieldTag:    "pms",
		ValidTag:    "valid",
		CacheSchema: true,
		FieldDirect: false,
		NotSnake:    false,
		NotDefValue: false,
		NotValid:    true,
	}
}

// 应用在解析配置文件的场景
func ConfigOptions() *ApplyOptions {
	return &ApplyOptions{
		FieldTag:    "pms",
		ValidTag:    "valid",
		CacheSchema: false,
		FieldDirect: true,
		NotSnake:    false,
		NotDefValue: false,
		NotValid:    false,
	}
}

// cst.KV
func ApplyKV(dst any, kvs cst.KV, opts *ApplyOptions) error {
	if err := applyKVToStruct(dst, kvs, opts); err != nil {
		return err
	}
	if opts.NotValid == false {
		return Validate(dst)
	}
	return nil
}

// JSON
func ApplyJsonReader(dst any, reader io.Reader, opts *ApplyOptions) error {
	return DecodeJsonReader(dst, reader, opts)
}

func ApplyJsonBytes(dst any, content []byte, opts *ApplyOptions) error {
	return DecodeJsonBytes(dst, content, opts)
}

// Yaml
func ApplyYamlReader(dst any, reader io.Reader, opts *ApplyOptions) error {
	return DecodeYamlReader(dst, reader, opts)
}

func ApplyYamlBytes(dst any, content []byte, opts *ApplyOptions) error {
	return DecodeYamlBytes(dst, content, opts)
}
