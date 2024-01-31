// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/randx"
)

// 生成新的 Guid
// 例如：YXRJT0l5ckpYNldBTjYzNHZw
func genSessGuid(cap int) []byte {
	sidLen := MySessDB.SidSize
	randomLen := int((sidLen*3 + 3) / 4)

	// base64编码需要字节 = MySessDB.SidSize or MySessDB.SidSize+1
	b64Len := base64Enc.EncodedLen(randomLen)

	minLen := b64Len + randomLen
	if cap < minLen {
		cap = minLen
	}

	// 一次性申请到运算过程中用到的所有内存
	buf := make([]byte, minLen, cap)
	base64Enc.Encode(buf[0:b64Len], randx.RandomFill(buf[b64Len:minLen]))

	return buf[0:sidLen]
}

// md5
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// md5是128bit的输出，而sha256是256bit输出
// 如果base64填充=，会得到一个长度为 24 的字符串
// 如果不填充=，会得到一个长度为 22 的字符串
func md5B64Str(data, secret []byte) string {
	buf := make([]byte, md5B64Len+md5Len)
	md5Fill(data, secret, buf[md5B64Len:md5B64Len])
	base64Enc.Encode(buf[:md5B64Len], buf[md5B64Len:])
	return lang.BTS(buf[:md5B64Len])
}

// 返回的md5值，存放在底层自己生成的字节切片对象
func md5Value(data, secret []byte) []byte {
	mac := hmac.New(md5.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}

// Note：必须 cap(buf) >= 16 ，md5值 append 在buf的后面，通常 len(buf) = 0
func md5Fill(data, secret, buf []byte) {
	mac := hmac.New(md5.New, secret)
	mac.Write(data)
	mac.Sum(buf)
}

// sha256
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 如果base64填充=，会得到一个长度为 44 的字符串
// 如果不填充=，会得到一个长度为 43 的字符串
func sha256Base64(data, secret []byte) string {
	return base64Enc.EncodeToString(sha256Value(data, secret))
}

func sha256Value(data, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}
