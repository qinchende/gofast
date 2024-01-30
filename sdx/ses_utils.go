// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"github.com/qinchende/gofast/skill/randx"
)

// md5是128bit的输出，而sha256是256bit输出
// 如果base64填充=，会得到一个长度为 24 的字符串
// 如果不填充=，会得到一个长度为 22 的字符串
func md5Base64(data, secret []byte) string {
	return base64.RawURLEncoding.EncodeToString(md5Value(data, secret))
}

func md5Value(data, secret []byte) []byte {
	mac := hmac.New(md5.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}

// 如果base64填充=，会得到一个长度为 44 的字符串
// 如果不填充=，会得到一个长度为 43 的字符串
func sha256Base64(data, secret []byte) string {
	return base64.RawURLEncoding.EncodeToString(sha256Value(data, secret))
}

func sha256Value(data, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return mac.Sum(nil)
}

// 闪电侠Guid
// 例如：YXRJT0l5ckpYNldBTjYzNHZw
func genSessGuid(cap int) []byte {
	strLen := int((MySessDB.SidSize*3 + 3) / 4)
	encLen := base64.RawURLEncoding.EncodedLen(strLen)
	if cap < encLen {
		cap = encLen
	}

	buf := make([]byte, encLen, cap)
	base64.RawURLEncoding.Encode(buf, randx.RandomBytes(strLen))

	return buf[0:MySessDB.SidSize]
}
