package lang

func ToBytes(stream any) (data []byte) {
	switch stream.(type) {
	case string:
		data = STB(stream.(string))
	case []byte:
		data = stream.([]byte)
	default:
		str := ToString(stream)
		data = STB(str)
	}
	return
}
