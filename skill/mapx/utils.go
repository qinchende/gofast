package mapx

import "errors"

var (
	errNotKVType    = errors.New("only map-like configs supported")
	errNotArrayType = errors.New("only array-like configs supported")
)

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

func Repr(v any) string {
	if v == nil {
		return ""
	}
	return sdxAsString(v)
}
