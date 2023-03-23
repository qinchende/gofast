package gson

import "github.com/qinchende/gofast/skill/lang"

type GsonRows struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]any
}

type GsonRow struct {
	Cls []string
	Row []any
}

// 实现接口 cst.SuperKV
func (gr *GsonRow) Get(k string) (v any, ok bool) {
	idx := lang.SearchSortStrings(gr.Cls, k)
	if idx < 0 {
		return nil, false
	}
	return gr.Row[idx], true // 有可能找到字段了，但是存的值是nil
}
