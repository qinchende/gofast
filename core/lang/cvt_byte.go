package lang

func ToBytes(stream any) (data []byte) {
	switch stream.(type) {
	case string:
		data = S2B(stream.(string))
	case []byte:
		data = stream.([]byte)
	default:
		str := ToString(stream)
		data = S2B(str)
	}
	return
}
