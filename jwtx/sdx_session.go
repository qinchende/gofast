package jwtx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/bytesconv"
	"github.com/qinchende/gofast/skill/json"
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

// TODO: 执行 session 验证
// 所有请求先经过这里验证 session 信息
// 每一次的访问，都必须要有一个 token ，没有token的 访问将视为 非法.
// 第一次没有 token 的情况下，默认造一个 token
func SdxSessHandler(ctx *fst.Context) {
	// 不可重复执行 token 检查，Sess构造的过程
	if ctx.Sess != nil {
		return
	}

	ctx.Sess = &fst.GFSession{Saved: true}
	tok := ctx.Pms["tok"]

	// 没有 tok，新建一个token，同时走后门的逻辑
	if tok == "" {
		sid, tok := ss.newToken(ctx)
		ctx.Sess.IsNew = true
		ctx.Sess.Sid = sid
		ctx.Sess.Token = tok
		ctx.Pms["tok"] = tok
		return
	}

	// 有 tok ，解析出 Sid
	reqSid, reqHash, err := fetchSid(tok)
	if err != nil {
		fst.RaisePanicErr(err)
	}

	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确， 需要分配新的Token。
	// 2. 过期，  用当前Token重建Session记录。
	isValid := ss.checkToken(reqSid, reqHash, ctx)

	// 如果验证通过
	if isValid {
		ss.getSessData(ctx)
	} else {
		fst.RaisePanic("check token error. ")
	}
}

// 验证是否登录
func SdxMustLoginHandler(ctx *fst.Context) {
	if ctx.Sess.Values[ss.SessKey] == "" {
		ctx.FaiMsg("not login ", "")
		fst.RaisePanic("not login")
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++==
// 从 redis 中获取 session data.
func (ss *SdxSession) getSessData(ctx *fst.Context) {
	str, err := ss.Redis.Get("tls:" + ctx.Sess.Sid)
	if str == "" || err != nil {
		str = "{}"
	}

	err = json.Unmarshal(bytesconv.StringToBytes(str), &ctx.Sess.Values)
	if err != nil {
		fst.RaisePanicErr(err)
	}
}

func (ss *SdxSession) newToken(ctx *fst.Context) (string, string) {
	return genToken(ss.Secret + ctx.ClientIP())
}

func (ss *SdxSession) checkToken(sid, hash string, ctx *fst.Context) bool {
	signSHA256 := genSignSHA256([]byte(sid), []byte(ss.Secret+ctx.ClientIP()))
	return hash == cleanString(signSHA256)
}

func fetchSid(tok string) (string, string, error) {
	start := strings.Index(tok, "t:")
	dot := strings.Index(tok, ".")
	if start != 0 || dot <= 0 {
		return "", "", errors.New("Can't find sid. ")
	}
	sid := tok[2:dot]
	if len(sid) <= 18 {
		return "", "", errors.New("Sid length error. ")
	}
	hash := tok[(dot + 1):]

	return sid, hash, nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++==
func genToken(secret string) (string, string) {
	sid := genSid(24)
	tok := "t:" + genSign(sid, secret)
	return sid, tok
}

// 按照指定长度length, 自动生成随机的Sid字符串，
func genSid(length int) string {
	src := lang.GetRandomBytes(length)
	sid := base64.StdEncoding.EncodeToString(src)
	sid = cleanString(sid)

	if length > len(sid) {
		length = len(sid)
	}
	return sid[:length]
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
