package lang

import "math"

var Placeholder PlaceholderType

type (
	PlaceholderType = struct{}
)

// 用泛型的方式获取一个值的地址
func Ptr[T any](x T) *T {
	return &x
}

func Round64(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func Round32(f float32, n int) float32 {
	return float32(Round64(float64(f), n))
}
