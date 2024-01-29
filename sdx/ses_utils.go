// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
)

// md5是128bit的输出，而sha256是256bit输出
// 如果base64填充=，会得到一个长度为 24 的字符串
// 如果不填充=，会得到一个长度为 22 的字符串
func md5Base64(data, secret []byte) string {
	mac := hmac.New(md5.New, secret)
	mac.Write(data)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// 如果base64填充=，会得到一个长度为 44 的字符串
// 如果不填充=，会得到一个长度为 43 的字符串
func sha256Base64(data, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(data)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
