// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/store/dts"
)

const (
	PmsToken = "tok"
)

type BaseFields struct {
	Tok string `v:"len=[64:128]"`
}

type Prefixes struct {
	PrefixToken   string `v:"def=t:"`   // token 字符串的 前缀
	PrefixSessKey string `v:"def=ses:"` // session 的前缀
}

type SessionCnf struct {
	Prefixes
	RedisConn  redis.ConnCnf `v:""`                               // 用 Redis 做持久化
	UidField   string        `v:"def=uid"`                        // 标记当前登录用户字段是? 比如：user_id
	Secret     string        `v:"required,def=sdx"`               // token秘钥
	TTL        int32         `v:"def=14400,range=[0:2000000000]"` // session有效期 默认 3600*4 秒
	TTLNew     int32         `v:"def=180,range=[0:2000000000]"`   // 首次产生的session有效期 默认 60*3 秒
	MustKeepIP bool          `v:"def=true"`                       // 看是否检查 token ip 地址
}

// 参数配置，Redis实例等
type SessionDB struct {
	SessionCnf
	Redis *redis.GfRedis
}

var _BasePms = []string{PmsToken}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 构造给定对象的字段名数组，加上公共的字段
func PmsKeys(obj any) []string {
	ss := dts.SchemaAsReq(obj)
	cls := ss.Columns
	newCls := make([]string, len(cls)+len(_BasePms))
	copy(newCls, cls)
	copy(newCls[len(cls):], _BasePms)
	return newCls // TODO: 可能需要考虑排重
}

// 每个进程只有一个全局 SdxSS 配置对象
var MySessDB *SessionDB

// 采用 “闪电侠” session 方案的时候需要先初始化参数
func SetupSession(ss *SessionDB) {
	if MySessDB != nil {
		return
	}
	MySessDB = ss

	if ss.Redis == nil {
		ss.Redis = redis.NewGoRedis(&ss.RedisConn)
	}
}
