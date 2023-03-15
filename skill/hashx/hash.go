package hashx

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/spaolacci/murmur3"
)

func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

func Sum64(data []byte) uint64 {
	return murmur3.Sum64(data)
}

func Sum32(data []byte) uint32 {
	return murmur3.Sum32(data)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Md5(data []byte) []byte {
	digest := md5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

func Md5HexBytes(data []byte) string {
	src := Md5(data)
	dst := make([]byte, len(src)*2)
	hex.Encode(dst, src)
	return lang.BTS(dst)
}

func Md5HexString(data string) string {
	return Md5HexBytes(lang.STB(data))
}

// ++++++++++++++++++++++++++++++++
func Sha1(data []byte) []byte {
	digest := sha1.New()
	digest.Write(data)
	return digest.Sum(nil)
}

func Sha1HexBytes(data []byte) string {
	src := Sha1(data)
	dst := make([]byte, len(src)*2)
	hex.Encode(dst, src)
	return lang.BTS(dst)
}

func Sha1HexString(data string) string {
	return Sha1HexBytes(lang.STB(data))
}

func Sha256(data []byte) []byte {
	digest := sha256.New()
	digest.Write(data)
	return digest.Sum(nil)
}

// ++++++++++++++++++++++++++++++++
func HmacSha256(key, data []byte) []byte {
	digest := hmac.New(sha256.New, key)
	digest.Write(data)
	return digest.Sum(nil)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func Md5(data string) string {
//	md5 := md5.New()
//	md5.Write([]byte(data))
//	md5Data := md5.Sum([]byte(""))
//	return hex.EncodeToString(md5Data)
//}

//func Hmac(key, data string) string {
//	hmac := hmac.New(md5.New, []byte(key))
//	hmac.Write([]byte(data))
//	return hex.EncodeToString(hmac.Sum([]byte("")))
//}

//func Sha1(data string) string {
//	sha1 := sha1.New()
//	sha1.Write([]byte(data))
//	return hex.EncodeToString(sha1.Sum([]byte("")))
//}
