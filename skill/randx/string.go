package randx

import (
	"github.com/qinchende/gofast/skill/lang"
)

const (
	sOnlyNumbers = "123456789"
	sLetters     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// 数字验证码
func RandomNumbers(size int) string {
	return lang.BTS(randomBytes(sOnlyNumbers, size))
}

func RandomString(size int) string {
	return lang.BTS(randomBytes(sLetters, size))
}

func RandomBytes(size int) []byte {
	return randomBytes(sLetters, size)
}

// 参数bs必须是一个带长度的切片内存空间
func RandomFill(bs []byte) []byte {
	for i := 0; i < len(bs); i++ {
		bs[i] = sLetters[seed.Intn(len(sLetters))]
	}
	return bs
}

func randomBytes(source string, size int) []byte {
	ret := make([]byte, size, size)
	for i := 0; i < size; i++ {
		ret[i] = source[seed.Intn(len(source))]
	}
	return ret
}

//// GetRandomKey creates a random key with the given size in bytes.
//// On failure, returns nil.
////
//// Callers should explicitly check for the possibility of a nil return, treat
//// it as a failure of the system random number generator, and not continue.
//func RandomKey(size int) []byte {
//	k := make([]byte, size)
//	if _, err := io.ReadFull(cRand.Reader, k); err != nil {
//		return nil
//	}
//	return k
//}
