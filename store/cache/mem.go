package cache

var gfMemCache map[string]MemItem

func init() {
	gfMemCache = make(map[string]MemItem, 0)
}

type MemItem struct {
	expire uint64
	Val    any
}
