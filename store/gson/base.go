package gson

type Gson struct {
	Ct   int64
	Tt   int64
	Cls  []string
	Rows [][]any
}

type GsonOne struct {
	Cls []string
	Row []any
}
