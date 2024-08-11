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
func (kv KV) Get(k string) (v any, ok bool) {
	v, ok = kv[k]
	return
}

func (kv KV) Set(k string, v any) {
	kv[k] = v
}

func (kv KV) Del(k string) {
	delete(kv, k)
}

func (kv KV) Len() int {
	return len(kv)
}

func (kv KV) GetString(k string) (v string, ok bool) {
	tmp, ok := kv[k]
	v = tmp.(string)
	return
}

func (kv KV) SetString(k string, v string) {
	kv[k] = v
}

// WebKV
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (wkv WebKV) Get(k string) (v any, ok bool) {
	v, ok = wkv[k]
	return
}

func (wkv WebKV) Set(k string, v any) {
	wkv[k] = v.(string)
}

func (wkv WebKV) Del(k string) {
	delete(wkv, k)
}

func (wkv WebKV) Len() int {
	return len(wkv)
}

func (wkv WebKV) GetString(k string) (v string, ok bool) {
	v, ok = wkv[k]
	return
}

func (wkv WebKV) SetString(k string, v string) {
	wkv[k] = v
}
