package jwtx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/lang"
	"regexp"
	"strings"
)

type SdxSessConfig struct {
	SessKey string `json:",optional"`
	Secret  string `json:",optional"`

	SessTTL    int `json:",optional"`
	SessTTLNew int `json:",optional"`
}

type SdxSession struct {
	SdxSessConfig
	Redis *redis.GoRedisX
}

var ss *SdxSession

func InitSdxRedis(i *SdxSession) {
	ss = i
	if ss.SessTTL == 0 {
		ss.SessTTL = 3600 * 4 // 默认4个小时
	}
	if ss.SessTTLNew == 0 {
		ss.SessTTLNew = 180 // 默认三分钟
	}
}

// 所有请求先经过这里验证 session 信息
func SdxSessHandler(ctx *fst.Context) {
	// 不可重复执行 token 检查，Sess构造的过程
	if ctx.Sess != nil {
		return
	}

	ctx.Sess = &fst.GFSession{Saved: true}
	tok := ctx.Pms["tok"]

	// 没有 tok
	if tok == "" {
		uid, tok := ss.newToken(ctx)
		ctx.Sess.IsNew = true
		ctx.Sess.Uid = uid
		ctx.Sess.Token = tok
		ctx.Pms["tok"] = tok
		return
	}

	//// 有 tok
	//// 解析出UID
	//if uid, err := getUid(tok); err != nil {
	//
	//}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++==
func (ss *SdxSession) getSessData(ctx *fst.Context) {
}

func (ss *SdxSession) newToken(ctx *fst.Context) (string, string) {
	return genToken(ss.Secret + ctx.ClientIP())
}

func genToken(secret string) (string, string) {
	uid := genUid(24)
	tok := "t:" + genSign(uid, secret)
	return uid, tok
}

// 按照指定长度lth, 自动生成随机的Uid字符串，
func genUid(lth int) string {
	src := lang.GetRandomBytes(lth)
	uid := base64.StdEncoding.EncodeToString(src)
	uid = dropChars(uid)

	if lth > len(uid) {
		lth = len(uid)
	}
	return uid[:lth]
}

func genSign(val, secret string) string {
	signSHA256 := genSignSHA256([]byte(val), []byte(secret))
	return val + "." + dropChars(signSHA256)
}

func genSignSHA256(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)

	// toBase64
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func getUid(tok string) (string, error) {
	start := strings.Index(tok, "t:")
	dot := strings.Index(tok, ".")
	if start != 0 || dot <= 0 {
		return "", errors.New("Can't find uid. ")
	}
	uid := tok[2:dot]
	if len(uid) <= 18 {
		return "", errors.New("Uid length error. ")
	}
	return uid, nil
}

func dropChars(src string, charts ...string) string {
	var str string
	if charts == nil {
		str = "/[+=]/g"
	} else {
		str = charts[0]
	}
	regExp, _ := regexp.Compile(str)
	return regExp.ReplaceAllString(src, "")
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// redis
//type session struct {
//	name    string
//	request *http.Request
//	store   Store
//	session *sessions.Session
//	saved   bool
//	writer  http.ResponseWriter
//}

func (ss *SdxSession) Get(ctx *fst.Context) {

}
func (ss *SdxSession) Set(ctx *fst.Context) {

}

func (ss *SdxSession) Save(ctx *fst.Context) {

}

func (ss *SdxSession) Delete(ctx *fst.Context) {

}
