package jwtx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/bytesconv"
	"github.com/qinchende/gofast/skill/json"
	"time"
)

// 从 redis 中获取 当前 请求上下文的 session data.
// TODO: 有可能 session 是空的
func (ss *SdxSession) loadSessionFromRedis(ctx *fst.Context) {
	str, err := ss.Redis.Get(sdxSessKeyPrefix + ctx.Sess.Sid)
	if str == "" || err != nil {
		str = `{}`
	}
	err = json.Unmarshal(bytesconv.StringToBytes(str), &ctx.Sess.Values)
	if err != nil {
		ctx.Fai(110, "获取SESSION失败，请重新访问系统。", fst.KV{})
	}
}

// 保存到 redis
func saveSessionToRedis(sdx *fst.CtxSession) (string, error) {
	str, _ := json.Marshal(sdx.Values)
	ttl := SdxSS.TTL
	if sdx.TokenIsNew && sdx.Values[SdxSS.AuthField] == nil {
		ttl = SdxSS.TTLNew
	}
	return SdxSS.Redis.Set(sdxSessKeyPrefix+sdx.Sid, str, time.Duration(ttl)*time.Second)
}

// 设置Session过期时间
func setSessionExpire(sdx *fst.CtxSession, ttl int32) (bool, error) {
	if ttl <= 0 {
		ttl = SdxSS.TTL
	}
	return SdxSS.Redis.Expire(sdxSessKeyPrefix+sdx.Sid, time.Duration(ttl)*time.Second)
}

// TODO: 这里的函数很多都没有考虑发生错误的情况
func destroySession(sdx *fst.CtxSession) {
	_, _ = SdxSS.Redis.Del(sdxSessKeyPrefix + sdx.Sid)
}
