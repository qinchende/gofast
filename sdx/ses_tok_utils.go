// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/randx"
	"regexp"
	"strings"
	"time"
)

// crypto
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// tok=t:NFRRcE81WDFQSEZJQUptZkpJ.v9EN6bWz8KU6sKRrcEId1OKUKqYx0hed2zSpCQImvc
// 解析不到，都将返回空字符串
func parseToken(tok string) (string, string) {
	start := strings.Index(tok, MySessDB.PrefixToken)
	dot := strings.Index(tok, ".")
	// 格式明显不对，直接返回空
	if start != 0 || dot <= 0 {
		return "", ""
	}
	guid := tok[2:dot]
	if len(guid) <= 18 {
		return "", ""
	}
	sHmac := tok[(dot + 1):]
	return guid, sHmac
}

// 闪电侠Guid：为24位的字符串
func genToken(secret string) (string, string) {
	guid := genGuid(24)
	tok := MySessDB.PrefixToken + genSign(guid, secret)
	return guid, tok
}

// 按照指定长度length, 自动生成随机的Guid字符串，
func genGuid(length int) string {
	src := randx.RandomBytes(length)
	guid := base64.StdEncoding.EncodeToString(src)
	guid = cleanString(guid)

	if length > len(guid) {
		length = len(guid)
	}
	return guid[:length]
}

func genSign(val, secret string) string {
	signSHA256 := genSignSHA256(lang.STB(val), lang.STB(secret))
	return val + "." + cleanString(signSHA256)
}

func genSignSHA256(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)

	// toBase64
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func cleanString(src string) string {
	regExp := regexp.MustCompile("[+=]*")
	//regExp := regexp.MustCompile("[+=/]*")
	return regExp.ReplaceAllString(src, "")
}

// 利用当前 guid 和 c 中包含的 request_ip | 计算出hmac值，然后和token中携带的 hmac值比较，来得出合法性
func checkToken(guid, sHmac, secret string) bool {
	signSHA256 := genSignSHA256(lang.STB(guid), lang.STB(secret))
	return sHmac == cleanString(signSHA256)
}

// Redis
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 从 redis 中获取 当前 请求上下文的 session data.
// TODO: 有可能 session 是空的
func (ss *TokSession) loadSessionFromRedis() error {
	str, err := MySessDB.Redis.Get(MySessDB.PrefixSessKey + ss.guid)
	if str == "" || err != nil {
		str = "{}"
	}
	return jsonx.Unmarshal(&ss.values, lang.STB(str))
}

// 保存到 redis
func (ss *TokSession) saveSessionToRedis() (string, error) {
	str, _ := jsonx.Marshal(ss.values)
	ttl := MySessDB.TTL
	if _, ok := ss.values[MySessDB.UidField]; ss.tokenIsNew && !ok {
		ttl = MySessDB.TTLNew
	}
	return MySessDB.Redis.Set(MySessDB.PrefixSessKey+ss.guid, str, time.Duration(ttl)*time.Second)
}

// 设置Session过期时间
func (ss *TokSession) setSessionExpire(ttl int32) (bool, error) {
	if ttl <= 0 {
		ttl = MySessDB.TTL
	}
	return MySessDB.Redis.Expire(MySessDB.PrefixSessKey+ss.guid, time.Duration(ttl)*time.Second)
}

// TODO: 这里的函数很多都没有考虑发生错误的情况
func (ss *TokSession) destroySession() (err error) {
	_, err = MySessDB.Redis.Del(MySessDB.PrefixSessKey + ss.guid)
	return
}
