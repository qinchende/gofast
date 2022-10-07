package lang

var Placeholder PlaceholderType

type (
	PlaceholderType = struct{}
)

// 用泛型的方式获取一个值的地址
func Ptr[T any](x T) *T {
	return &x
}
