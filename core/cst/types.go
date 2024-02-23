package cst

type (
	KV    map[string]any
	WebKV map[string]string
	//WebValues map[string][]string

	TypeError  error
	TypeInt    int
	TypeString string
)

// 可能用map，也可能自定义数组等合适的数据结构存取。
// 比如上下文中用来保存解析到的请求数据，主要是KV形式
type SuperKV interface {
	Set(k string, v any)
	Get(k string) (any, bool)
	Del(k string)
	Len() int
	SetString(k string, v string)      // 如果Value是String，提高性能
	GetString(k string) (string, bool) // 如果Value是String，提高性能
}

// KV
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (kvs KV) Get(k string) (v any, ok bool) {
	v, ok = kvs[k]
	return
}

func (kvs KV) Set(k string, v any) {
	kvs[k] = v
}

func (kvs KV) Del(k string) {
	delete(kvs, k)
}

func (kvs KV) Len() int {
	return len(kvs)
}

func (kvs KV) GetString(k string) (v string, ok bool) {
	tmp, ok := kvs[k]
	v = tmp.(string)
	return
}

func (kvs KV) SetString(k string, v string) {
	kvs[k] = v
}

// WebKV
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (kvs WebKV) Get(k string) (v any, ok bool) {
	v, ok = kvs[k]
	return
}

func (kvs WebKV) Set(k string, v any) {
	kvs[k] = v.(string)
}

func (kvs WebKV) Del(k string) {
	delete(kvs, k)
}

func (kvs WebKV) Len() int {
	return len(kvs)
}

func (kvs WebKV) GetString(k string) (v string, ok bool) {
	v, ok = kvs[k]
	return
}

func (kvs WebKV) SetString(k string, v string) {
	kvs[k] = v
}
