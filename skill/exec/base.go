package exec

import "time"

const defaultFlushInterval = time.Second

type Execute func(tasks []any)
