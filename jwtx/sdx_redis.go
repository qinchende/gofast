package jwtx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/bytesconv"
	"github.com/qinchende/gofast/skill/json"
)

// 从 redis 中获取 当前 请求上下文的 session data.
// TODO: 有可能 session 是空的
func (ss *SdxSession) initCtxSess(ctx *fst.Context) {
	str, err := ss.Redis.Get(sdxSessKeyPrefix + ctx.Sess.Sid)
	if str == "" || err != nil {
		str = `{}`
	}
	err = json.Unmarshal(bytesconv.StringToBytes(str), &ctx.Sess.Values)
	if err != nil {
		ctx.FaiX(110, "获取SESSION失败，请重新访问系统。", fst.KV{})
	}
}

// 保存到 redis
func SaveSessionToRedis(sdx *fst.CtxSession) (string, error) {
	str, _ := json.Marshal(sdx.Values)
	ttl := ses.TTL
	if sdx.IsNew {
		ttl = ses.TTLNew
	}
	return ses.Redis.Set(sdxSessKeyPrefix+sdx.Sid, str, ttl)
}
