package jde

func DecodeGsonRowsFromBytes(v any, source []byte) error {
	//return decodeFromString(v, source)
	return nil
}

func DecodeGsonRowsFromString(v any, source string) error {
	return decGsonRowsFromString(v, source)
}
