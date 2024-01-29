// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/store/dts"
	"time"
)

const (
	PmsToken = "tok"
)

type BaseFields struct {
	Tok string `v:"len=[64:128]"`
}

var _BasePms = []string{PmsToken}

// 构造给定对象的字段名数组，加上公共的字段
func PmsKeys(obj any) []string {
	ss := dts.SchemaAsReq(obj)
	cls := ss.Columns
	newCls := make([]string, len(cls)+len(_BasePms))
	copy(newCls, cls)
	copy(newCls[len(cls):], _BasePms)
	return newCls // TODO: 可能需要考虑排重
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
type SessionCnf struct {
	//PrefixToken   string        `v:"def=t:"`                         // token 字符串的 前缀
	RedisConn     redis.ConnCnf `v:""`                           // 用 Redis 做持久化
	PrefixSessKey string        `v:"def=ses:"`                   // session 的前缀
	UidField      string        `v:"def=uid"`                    // 标记当前登录用户字段是? 比如：user_id
	Secret        string        `v:"required,def=sdx"`           // 用于计算token的秘钥
	SecretLast    string        `v:"def=sdx"`                    // 上一个密钥，可能在更换密钥时有用
	TTL           time.Duration `v:"def=14400s,range=[0s:240h]"` // session有效期 默认 3600*4 秒
	TTLNew        time.Duration `v:"def=180s,range=[0s:1h]"`     // 首次产生的session有效期 默认 60*3 秒
	SidSize       uint8         `v:"def=24"`                     // session id (uuid)长度
	MustKeepIP    bool          `v:"def=false"`                  // 看是否检查 token ip 地址
}

// 参数配置，Redis实例等
type SessionDB struct {
	SessionCnf
	Redis *redis.GfRedis
}

// 每个进程只有一个全局 SdxSS 配置对象
var MySessDB *SessionDB

// 采用 “闪电侠” session 方案的时候需要先初始化参数
func SetSessionDB(ss *SessionDB) {
	MySessDB = ss
	if ss.Redis == nil {
		ss.Redis = redis.NewGoRedis(&ss.RedisConn)
	}
}
