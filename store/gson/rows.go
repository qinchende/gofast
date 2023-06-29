package gson

type GsonRows struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]any
}

type RowsRet struct {
	Err  error
	Ct   int64
	Tt   int64
	Scan int
}

type RowsPet struct {
	Ct     int64
	Tt     int64
	Target any
}
