package gson

type GsonRows struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]any
}

type GsonRowsSuper struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]FValue
}
