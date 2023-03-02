package randx

import (
	"github.com/qinchende/gofast/skill/lang"
)

const (
	sOnlyNumbers = "123456789"
	sLetters     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// 数字验证码
func RandomNumbers(length int) string {
	return lang.BytesToString(randomBytes(sOnlyNumbers, length))
}

func RandomString(length int) string {
	return lang.BytesToString(randomBytes(sLetters, length))
}

func RandomBytes(length int) []byte {
	return randomBytes(sLetters, length)
}

func randomBytes(source string, length int) []byte {
	ret := make([]byte, length, length)
	for i := 0; i < length; i++ {
		ret[i] = source[seed.Intn(len(source))]
	}
	return ret
}

//// GetRandomKey creates a random key with the given length in bytes.
//// On failure, returns nil.
////
//// Callers should explicitly check for the possibility of a nil return, treat
//// it as a failure of the system random number generator, and not continue.
//func RandomKey(length int) []byte {
//	k := make([]byte, length)
//	if _, err := io.ReadFull(cRand.Reader, k); err != nil {
//		return nil
//	}
//	return k
//}
