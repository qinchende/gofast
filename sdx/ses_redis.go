package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/stringx"
	"time"
)

// 从 redis 中获取 当前 请求上下文的 session data.
// TODO: 有可能 session 是空的
func (ss *CtxSession) loadSessionFromRedis(c *fst.Context) {
	str, err := MySessDB.Redis.Get(sdxSessKeyPrefix + ss.Guid)
	if str == "" || err != nil {
		str = `{}`
	}
	err = jsonx.Unmarshal(&ss.Values, stringx.StringToBytes(str))
	if err != nil {
		c.AbortHandlers()
		c.Fai(110, "获取SESSION失败，请重新访问系统。", fst.KV{})
	}
}

// 保存到 redis
func (ss *CtxSession) saveSessionToRedis() (string, error) {
	str, _ := jsonx.Marshal(ss.Values)
	ttl := MySessDB.TTL
	if ss.TokenIsNew && ss.Values[MySessDB.GuidField] == nil {
		ttl = MySessDB.TTLNew
	}
	return MySessDB.Redis.Set(sdxSessKeyPrefix+ss.Guid, str, time.Duration(ttl)*time.Second)
}

// 设置Session过期时间
func (ss *CtxSession) setSessionExpire(ttl int32) (bool, error) {
	if ttl <= 0 {
		ttl = MySessDB.TTL
	}
	return MySessDB.Redis.Expire(sdxSessKeyPrefix+ss.Guid, time.Duration(ttl)*time.Second)
}

// TODO: 这里的函数很多都没有考虑发生错误的情况
func (ss *CtxSession) destroySession() {
	_, _ = MySessDB.Redis.Del(sdxSessKeyPrefix + ss.Guid)
}
