package hash

import (
	"crypto/md5"
	"fmt"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/spaolacci/murmur3"
)

func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

func Md5(data []byte) []byte {
	digest := md5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

func Md5HexBytes(data []byte) string {
	return fmt.Sprintf("%x", Md5(data))
}

func Md5HexString(data string) string {
	return fmt.Sprintf("%x", Md5(lang.StringToBytes(data)))
}
