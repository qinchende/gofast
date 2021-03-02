package jwtx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
	"time"
)

type SdxSessConfig struct {
	Redis   redis.ConnConfig `json:",optional"` // 用 Redis 做持久化
	SessKey string           `json:",optional"` // 用户信息的主键
	Secret  string           `json:",optional"` // token秘钥
	TTL     time.Duration    `json:",optional"` // session有效期 默认 3600*4 秒
	TTLNew  time.Duration    `json:",optional"` // 首次产生的session有效期 默认 60*3 秒
}

type SdxSession struct {
	SdxSessConfig
	Redis *redis.GoRedisX
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var ss *SdxSession

func InitSdxRedis(sdx *SdxSession) {
	if ss != nil {
		return
	}
	ss = sdx
	if ss.TTL == 0 {
		ss.TTL = 3600 * 4 // 默认4个小时
	}
	if ss.TTLNew == 0 {
		ss.TTLNew = 180 // 默认三分钟
	}
	if ss.SessKey == "" {
		ss.SessKey = "cus_id"
	}

	// 设置保存session的逻辑
	fst.CtxSessionSaveFun = SaveRedis
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

	//ctx.Sess = &fst.GFSession{Saved: true}
	// 初始化的时候就需要查询一次 redis 得到session
	ctx.Sess = &fst.CtxSession{Saved: true}
	//InitRedis(ctx.Sess)

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
	ctx.Sess.Sid = reqSid

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
	if ctx.Sess.Values == nil {
		ctx.FaiMsg("No session fond.", "")
		fst.RaisePanic("No session fond. ")
	}

	uid := ctx.Sess.Get(ss.SessKey)
	//uid := ss.Get(ctx)
	if uid == nil || uid == "" {
		//ctx.FaiMsg("not login ", "")
		//fst.RaisePanic("not login")
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (ss *SdxSession) newToken(ctx *fst.Context) (string, string) {
	return genToken(ss.Secret + ctx.ClientIP())
}

func (ss *SdxSession) checkToken(sid, hash string, ctx *fst.Context) bool {
	signSHA256 := genSignSHA256([]byte(sid), []byte(ss.Secret+ctx.ClientIP()))
	return hash == cleanString(signSHA256)
}
