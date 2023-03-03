// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"time"
)

// 从 redis 中获取 当前 请求上下文的 session data.
// TODO: 有可能 session 是空的
func (ss *CtxSession) loadSessionFromRedis(c *fst.Context) error {
	str, err := MySessDB.Redis.Get(sdxSessKeyPrefix + ss.guid)
	if str == "" || err != nil {
		str = "{}"
	}
	return jsonx.Unmarshal(&ss.values, lang.StringToBytes(str))
}

// 保存到 redis
func (ss *CtxSession) saveSessionToRedis() (string, error) {
	str, _ := jsonx.Marshal(ss.values)
	ttl := MySessDB.TTL
	if ss.tokenIsNew && ss.values[MySessDB.UidField] == nil {
		ttl = MySessDB.TTLNew
	}
	return MySessDB.Redis.Set(sdxSessKeyPrefix+ss.guid, str, time.Duration(ttl)*time.Second)
}

// 设置Session过期时间
func (ss *CtxSession) setSessionExpire(ttl int32) (bool, error) {
	if ttl <= 0 {
		ttl = MySessDB.TTL
	}
	return MySessDB.Redis.Expire(sdxSessKeyPrefix+ss.guid, time.Duration(ttl)*time.Second)
}

// TODO: 这里的函数很多都没有考虑发生错误的情况
func (ss *CtxSession) destroySession() (err error) {
	_, err = MySessDB.Redis.Del(sdxSessKeyPrefix + ss.guid)
	return
}
