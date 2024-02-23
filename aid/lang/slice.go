package lang

import "math"

// 数组项排重
func RemoveRepeatItems[T comparable](list []T) []T {
	// 创建一个临时map用来存储数组元素，防止动态伸缩
	kvs := make(map[T]struct{}, int(math.Ceil(float64(len(list)*8)/6.5)))
	for i := range list {
		kvs[list[i]] = struct{}{}
	}

	newList := make([]T, len(kvs))
	idx := 0
	for key := range kvs {
		newList[idx] = key
		idx++
	}
	return newList
}
