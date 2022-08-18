package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/stringx"
	"time"
)

// 从 redis 中获取 当前 请求上下文的 session data.
// TODO: 有可能 session 是空的
func (ss *CtxSession) loadSessionFromRedis(c *fst.Context) error {
	str, err := MySess.Redis.Get(sdxSessKeyPrefix + ss.guid)
	if str == "" || err != nil {
		str = "{}"
	}
	return jsonx.Unmarshal(&ss.values, stringx.StringToBytes(str))
}

// 保存到 redis
func (ss *CtxSession) saveSessionToRedis() (string, error) {
	str, _ := jsonx.Marshal(ss.values)
	ttl := MySess.TTL
	if ss.tokenIsNew && ss.values[MySess.GuidField] == nil {
		ttl = MySess.TTLNew
	}
	return MySess.Redis.Set(sdxSessKeyPrefix+ss.guid, str, time.Duration(ttl)*time.Second)
}

// 设置Session过期时间
func (ss *CtxSession) setSessionExpire(ttl int32) (bool, error) {
	if ttl <= 0 {
		ttl = MySess.TTL
	}
	return MySess.Redis.Expire(sdxSessKeyPrefix+ss.guid, time.Duration(ttl)*time.Second)
}

// TODO: 这里的函数很多都没有考虑发生错误的情况
func (ss *CtxSession) destroySession() {
	_, _ = MySess.Redis.Del(sdxSessKeyPrefix + ss.guid)
}
