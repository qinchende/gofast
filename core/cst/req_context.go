package cst

//  返回结构体
type Ret struct {
	Code int    // 返回编码
	Msg  string // 文本消息
	Data any    // 携带数据体
	Desc string // 描述，内部说明，不对外传递和显示
}

func (ret Ret) Callback() {}
