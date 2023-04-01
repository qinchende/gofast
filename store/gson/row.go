package gson

import "github.com/qinchende/gofast/skill/lang"

// GsonField Data Type
const (
	Any int = iota
	String
	Int
	Float64
)

type GsonField struct {
	DType  int
	Name   string
	ValStr string
	ValAny any
}

type GsonRow struct {
	Cls []string
	Row []any
	Val []string
}

// TODO: 在这种模式下，GsonRow中的Cls必须是已经按照字符串长度从小到大排好序的
// 实现接口 cst.SuperKV
func (gr *GsonRow) Get(k string) (v any, ok bool) {
	idx := lang.SearchSortStrings(gr.Cls, k)
	if idx < 0 {
		return nil, false
	}
	return gr.Row[idx], true // 有可能找到字段了，但是存的值是nil
}

func (gr *GsonRow) Del(k string) {
	idx := lang.SearchSortStrings(gr.Cls, k)
	if idx >= 0 {
		gr.Row[idx] = nil
	}
}

func (gr *GsonRow) Set(k string, v any) {
	idx := lang.SearchSortStrings(gr.Cls, k)
	if idx >= 0 {
		gr.Row[idx] = v
	}
}

func (gr *GsonRow) SetString(k string, v string) {
	idx := lang.SearchSortStrings(gr.Cls, k)
	if idx >= 0 {
		gr.Val[idx] = v
		//gr.Row[idx] = gr.Val[idx]
	}
}

func (gr *GsonRow) Len() int {
	return len(gr.Cls)
}

// 高级功能 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (gr *GsonRow) KeyIndex(k string) int {
	return lang.SearchSortStrings(gr.Cls, k)
}

func (gr *GsonRow) GetKey(idx int) string {
	if idx < 0 || idx > gr.Len() {
		return ""
	}
	return gr.Cls[idx]
}

func (gr *GsonRow) GetValue(idx int) any {
	if idx < 0 || idx > gr.Len() {
		return nil
	}
	return gr.Row[idx]
}

func (gr *GsonRow) SetByIndex(idx int, v any) {
	if idx < 0 || idx > gr.Len() {
		return
	}
	gr.Row[idx] = v
}
