package gson

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
	Target   any
	Tt       int64
	FlsStr   string
	FlsIdxes []uint8
}
