package sdx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/qinchende/gofast/skill/lang"
	"regexp"
	"strings"
)

// tok=t:NFRRcE81WDFQSEZJQUptZkpJ.v9EN6bWz8KU6sKRrcEId1OKUKqYx0hed2zSpCQImvc
func fetchGuid(tok string) (string, string) {
	start := strings.Index(tok, sdxTokenPrefix)
	dot := strings.Index(tok, ".")
	// 格式明显不对，直接返回空
	if start != 0 || dot <= 0 {
		// return "", "", errors.New("Can't parse guid. ")
		return "", ""
	}
	guid := tok[2:dot]
	if len(guid) <= 18 {
		// return "", "", errors.New("guid length error. ")
		return "", ""
	}
	sHmac := tok[(dot + 1):]
	return guid, sHmac
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++==
// 闪电侠Guid：为24位的字符串
func genToken(secret string) (string, string) {
	guid := genGuid(24)
	tok := sdxTokenPrefix + genSign(guid, secret)
	return guid, tok
}

// 按照指定长度length, 自动生成随机的Guid字符串，
func genGuid(length int) string {
	src := lang.GetRandomBytes(length)
	guid := base64.StdEncoding.EncodeToString(src)
	guid = cleanString(guid)

	if length > len(guid) {
		length = len(guid)
	}
	return guid[:length]
}

func genSign(val, secret string) string {
	signSHA256 := genSignSHA256([]byte(val), []byte(secret))
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
	signSHA256 := genSignSHA256([]byte(guid), []byte(secret))
	return sHmac == cleanString(signSHA256)
}
