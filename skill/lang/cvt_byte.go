package lang

func ToBytes(stream any) (data []byte) {
	switch stream.(type) {
	case string:
		data = StringToBytes(stream.(string))
	case []byte:
		data = stream.([]byte)
	default:
		str := ToString(stream)
		data = StringToBytes(str)
	}
	return
}
