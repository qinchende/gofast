package exec

import "time"

const defaultFlushInterval = time.Second

type FuncExecute func(tasks []any)
