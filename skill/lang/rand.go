package lang

import (
	cRand "crypto/rand"
	"io"
	mRand "math/rand"
	"time"
)

func GetRandomString(length int) string {
	return string(GetRandomBytes(length))
}

func GetRandomBytes(length int) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := mRand.New(mRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}

func GetRandomInt(max int) int {
	r := mRand.New(mRand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

// GetRandomKey creates a random key with the given length in bytes.
// On failure, returns nil.
//
// Callers should explicitly check for the possibility of a nil return, treat
// it as a failure of the system random number generator, and not continue.
func GetRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(cRand.Reader, k); err != nil {
		return nil
	}
	return k
}
