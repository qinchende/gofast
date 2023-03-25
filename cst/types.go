package cst

type (
	KV        map[string]any
	WebKV     map[string]string
	WebValues map[string][]string

	TypeError  error
	TypeInt    int
	TypeString string
)

func (kvs KV) Get(k string) (v any, ok bool) {
	v, ok = kvs[k]
	return
}

func (kvs KV) Set(k string, v any) {
	kvs[k] = v
}
