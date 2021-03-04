package jwtx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/fst"
	"time"
)

type SdxSessConfig struct {
	Redis        redis.ConnConfig `json:",optional"`                // 用 Redis 做持久化
	CheckTokenIP bool             `json:",optional,default=true"`   // 看是否检查 token ip 地址
	SessKey      string           `json:",optional,default=cus_id"` // 用户信息的主键
	Secret       string           `json:",optional"`                // token秘钥
	TTL          time.Duration    `json:",optional"`                // session有效期 默认 3600*4 秒
	TTLNew       time.Duration    `json:",optional"`                // 首次产生的session有效期 默认 60*3 秒
}

// 每个进程只有一个全局 sdx session 配置对象
var ss *SdxSession

// 参数配置，Redis实例等
type SdxSession struct {
	SdxSessConfig
	Redis *redis.GoRedisX
}

func InitSdxSession(sdx *SdxSession) {
	if ss != nil {
		return
	}
	ss = sdx
	if ss.TTL == 0 {
		ss.TTL = 3600 * 4 * time.Second // 默认4个小时
	}
	if ss.TTLNew == 0 {
		ss.TTLNew = 180 * time.Second // 默认三分钟
	}
	//if ss.SessKey == "" {
	//	ss.SessKey = "cus_id"
	//}

	// 指定 保存session 的处理函数
	fst.CtxSessionSaveFun = SaveSessionToRedis
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

	ctx.Sess = &fst.CtxSession{Saved: false, IsNew: false}
	tok := ctx.Pms["tok"]

	// 没有 tok，新建一个token，假装当前请求是有token的，同时走后面的逻辑
	if tok == "" {
		ss.initNewToken(ctx)
		ss.initCtxSess(ctx)
		return
	}

	// 有 tok ，解析出 [sid、hmac]
	reqSid, reqHmac := fetchSid(tok)
	//if err != nil {
	//	fst.RaisePanicErr(err)
	//}
	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确，需要分配新的Token。
	// 2. 过期，用当前Token重建Session记录。
	isValid := ss.checkToken(reqSid, reqHmac, ctx)

	// 按照ip计算出当前hmac，和请求中的hmac相比较，看是否相等
	// 如果Sid验证通过
	if isValid || ss.CheckTokenIP == false {
		ctx.Sess.Sid = reqSid
	}

	// 如果没有sid，就新生成一个
	if ctx.Sess.Sid == "" {
		ss.initNewToken(ctx)
	}
	// 通过sid 到 redis 中获取当前 session
	ss.initCtxSess(ctx)
}

// 验证是否登录
func SdxMustLogin(ctx *fst.Context) {
	uid := ctx.Sess.Get(ss.SessKey)
	if uid == nil || uid == "" {
		ctx.FaiX(110, "认证失败，请先登录。", fst.KV{})
		return
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 生成新 token，初始化当前 session
func (ss *SdxSession) initNewToken(ctx *fst.Context) {
	sid, tok := ss.newToken(ctx)
	ctx.Sess.IsNew = true
	ctx.Sess.Sid = sid
	ctx.Sess.Token = tok
	ctx.Sess.Values = make(map[string]interface{})
	// ctx.Pms["tok"] = tok
}

// TODO: 需不需要安全级别更高的 IP 校验是个问题?
func (ss *SdxSession) newToken(ctx *fst.Context) (string, string) {
	return genToken(ss.Secret + ctx.ClientIP())
}

// 利用当前 sid 和 ctx 中包含的 request_ip | 计算出hmac值，然后和token中携带的 hmac值比较，来得出合法性
func (ss *SdxSession) checkToken(sid, sHmac string, ctx *fst.Context) bool {
	signSHA256 := genSignSHA256([]byte(sid), []byte(ss.Secret+ctx.ClientIP()))
	return sHmac == cleanString(signSHA256)
}
