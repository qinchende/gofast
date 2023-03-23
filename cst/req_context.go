package cst

//  返回结构体
type Ret struct {
	Code int    // 返回编码
	Msg  string // 文本消息
	Data any    // 携带数据体
	Desc string // 描述，内部说明，不对外传递和显示
}

func (ret Ret) Callback() {}

// 上下文中用来保存解析到的请求数据，主要是KV形式
// 可能用map，也可能自定义数组等合适的数据结构存取。
type SuperKV interface {
	Get(k string) (any, bool)
}
