package gson

import "github.com/qinchende/gofast/skill/lang"

// 每个字段值占用4个字，32字节
type FValue struct {
	Val any    // 任意类型值
	str string // 只记录字符串值
}

type GsonRow struct {
	Cls []string // 字段
	Row []FValue // 对应值，如果是nil，证明没有匹配到
	//values []string   // 真正的值
}

// TODO: 在这种模式下，GsonRow中的Cls必须是已经按照字符串长度从小到大排好序的
// 实现接口 cst.SuperKV
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gr *GsonRow) Set(k string, v any) {
	idx := lang.SearchSorted(gr.Cls, k)
	if idx >= 0 {
		gr.Row[idx].Val = v
	}
}

func (gr *GsonRow) Get(k string) (v any, ok bool) {
	idx := lang.SearchSorted(gr.Cls, k)
	if idx < 0 {
		return nil, false
	}
	return gr.Row[idx].Val, true // 有可能找到字段了，但是存的值是nil
}

func (gr *GsonRow) Del(k string) {
	idx := lang.SearchSorted(gr.Cls, k)
	if idx >= 0 {
		gr.Row[idx].Val = nil
	}
}

func (gr *GsonRow) Len() int {
	return len(gr.Cls)
}

// 绕一圈，主要是为了避免对象分配，提高性能。
func (gr *GsonRow) SetString(k string, v string) {
	idx := lang.SearchSorted(gr.Cls, k)
	if idx >= 0 {
		gr.Row[idx].str = v
		gr.Row[idx].Val = &gr.Row[idx].str
	}
}

func (gr *GsonRow) GetString(k string) (v string, ok bool) {
	idx := lang.SearchSorted(gr.Cls, k)
	if idx < 0 || gr.Row[idx].Val == nil {
		return "", false
	}
	switch gr.Row[idx].Val.(type) {
	case string:
		return gr.Row[idx].Val.(string), true
	}
	return "", false
}

// GsonRow 特有的高级功能 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 初始化内存空间
func (gr *GsonRow) Init(cls []string) {
	gr.Cls = cls
	gr.Row = make([]FValue, gr.Len())
	//gr.values = make([]string, gr.Len())
}

func (gr *GsonRow) KeyIndex(k string) int {
	return lang.SearchSorted(gr.Cls, k)
}

func (gr *GsonRow) GetKeyByIndex(idx int) string {
	if idx < 0 || idx > gr.Len() {
		return ""
	}
	return gr.Cls[idx]
}

func (gr *GsonRow) GetValue(idx int) any {
	if idx < 0 || idx > gr.Len() {
		return nil
	}
	return gr.Row[idx].Val
}

func (gr *GsonRow) SetByIndex(idx int, v any) {
	if idx < 0 || idx > gr.Len() {
		return
	}
	gr.Row[idx].Val = v
}

func (gr *GsonRow) SetStringByIndex(idx int, v string) {
	if idx < 0 || idx > gr.Len() {
		return
	}
	gr.Row[idx].str = v
	gr.Row[idx].Val = &gr.Row[idx].str
}

//// GsonField Data Type
//const (
//	Any int = iota
//	String
//	Int
//	Float64
//)
//
//type GsonField struct {
//	Typ int
//	Key string
//	Val any
//	str string
//}
