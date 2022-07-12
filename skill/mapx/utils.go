package mapx

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func ifPanic(yn bool, text string) {
	if yn {
		panic(text)
	}
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ToString(v any) string {
	if v == nil {
		return ""
	}
	return sdxAsString(v)
}
