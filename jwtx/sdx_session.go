package jwtx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
)

// 采用 “闪电侠” session 方案的时候需要先初始化参数
func (ss *SdxSession) Init() {
	if SdxSS != nil {
		return
	}
	SdxSS = ss

	// 给默认值
	//if ss.TTL == 0 {
	//	ss.TTL = 3600 * 4 * time.Second // 默认4个小时
	//}
	//if ss.TTLNew == 0 {
	//	ss.TTLNew = 180 * time.Second // 默认三分钟
	//}
	if ss.Redis == nil {
		ss.Redis = redis.NewGoRedis(&ss.RedisConnCnf)
	}

	// 指定 保存session 的处理函数
	fst.CtxSessionCreateFun = newSdxToken
	fst.CtxSessionDestroyFun = destroySession
	fst.CtxSessionExpireFun = setSessionExpire
	fst.CtxSessionSaveFun = saveSessionToRedis

	// 初始化完成
	SdxSS.isReady = true
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 还原 session ，验证合法性
// 所有请求先经过这里验证 session 信息
// 每一次的访问，都必须要有一个 token ，没有token的 访问将视为 非法.
// 第一次没有 token 的情况下，默认造一个 token
func SdxSessBuilder(ctx *fst.Context) {
	// 不可重复执行 token 检查，Sess构造的过程
	if ctx.Sess != nil {
		return
	}

	// 每个请求对应的SESSION对象都是新创建的，线程安全。
	ctx.Sess = new(fst.CtxSession)
	ctx.Sess.Values = make(map[string]interface{})
	tok := ctx.Pms["tok"]
	if tok == nil {
		tok = ""
	}
	ctx.Sess.Token = tok.(string)

	// 没有tok，赋予当前请求token，同时走后面的逻辑
	if len(ctx.Sess.Token) < 50 {
		newSdxToken(ctx)
		return
	}

	// 有 tok ，解析出 [sid、hmac]
	reqSid, reqHmac := fetchSid(ctx.Sess.Token)
	//if err != nil {
	//	fst.RaisePanicErr(err)
	//}
	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确，需要分配新的Token。
	// 2. 过期，用当前Token重建Session记录。
	isValid := SdxSS.checkToken(reqSid, reqHmac, ctx)

	// 按照ip计算出当前hmac，和请求中的hmac相比较，看是否相等
	// 如果Sid验证通过
	if isValid || SdxSS.CheckTokenIP == false {
		ctx.Sess.Sid = reqSid
	}

	// 如果没有sid，就新生成一个
	if ctx.Sess.Sid == "" {
		newSdxToken(ctx)
	} else {
		SdxSS.loadSessionFromRedis(ctx) // 通过sid 到 redis 中获取当前 session
	}
}

// 验证是否登录
func SdxMustLogin(ctx *fst.Context) {
	uid := ctx.Sess.Get(SdxSS.AuthField)
	if uid == nil || uid == "" {
		ctx.Fai(110, "认证失败，请先登录。", fst.KV{})
	}
}

// 新生成一个SDX Session对象，生成新的tok
func newSdxToken(ctx *fst.Context) {
	sid, tok := SdxSS.newToken(ctx)
	if ctx.Sess == nil {
		ctx.Sess = new(fst.CtxSession)
	}
	ctx.Sess.TokenIsNew = true
	ctx.Sess.Sid = sid
	ctx.Sess.Token = tok
	if ctx.Sess.Values == nil {
		ctx.Sess.Values = make(map[string]interface{})
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 需不需要安全级别更高的 IP 校验是个问题?
func (ss *SdxSession) newToken(ctx *fst.Context) (string, string) {
	return genToken(ss.Secret + ctx.ClientIP())
}

// 利用当前 sid 和 ctx 中包含的 request_ip | 计算出hmac值，然后和token中携带的 hmac值比较，来得出合法性
func (ss *SdxSession) checkToken(sid, sHmac string, ctx *fst.Context) bool {
	signSHA256 := genSignSHA256([]byte(sid), []byte(ss.Secret+ctx.ClientIP()))
	return sHmac == cleanString(signSHA256)
}
