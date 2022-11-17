// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
)

var (
	sdxTokenPrefix   = "t:"   // token 字符串的 前缀
	sdxSessKeyPrefix = "tls:" // session 的前缀
)

type RedisSessCnf struct {
	RedisConn  gfrds.ConnCnf `v:"required"`                       // 用 Redis 做持久化
	GuidField  string        `v:"def=uid"`                        // 标记当前登录用户字段是 user_id
	Secret     string        `v:"required;def=sdx"`               // token秘钥
	TTL        int32         `v:"def=14400,range=[0:2000000000]"` // session有效期 默认 3600*4 秒
	TTLNew     int32         `v:"def=180,range=[0:2000000000]"`   // 首次产生的session有效期 默认 60*3 秒
	MustKeepIP bool          `v:"def=true"`                       // 看是否检查 token ip 地址
}

// 参数配置，Redis实例等
type SessionDB struct {
	RedisSessCnf
	Redis *gfrds.GfRedis
	//isReady bool // 是否已经初始化
}

// 每个进程只有一个全局 SdxSS 配置对象
var MySess *SessionDB

// 采用 “闪电侠” session 方案的时候需要先初始化参数
func SetupSession(ss *SessionDB) {
	if MySess != nil {
		return
	}
	MySess = ss

	if ss.Redis == nil {
		ss.Redis = gfrds.NewGoRedis(&ss.RedisConn)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 默认将使用 Redis 存放 session 信息
// TODO: 注意，这个实现是非线程安全的
type CtxSession struct {
	values     cst.KV // map[string]interface{}
	guid       string // redis session key
	token      string // Sid
	tokenIsNew bool   // Sid is new
	saved      bool   // Whether it has been saved
}

// CtxSession 需要实现 sessionKeeper 所有接口
var _ fst.SessionKeeper = &CtxSession{}

func (ss *CtxSession) GetValues() cst.KV {
	return ss.values
}

func (ss *CtxSession) Get(key string) any {
	if ss.values == nil {
		return nil
	}
	return ss.values[key]
}

func (ss *CtxSession) Set(key string, val any) {
	ss.saved = false
	ss.values[key] = val
}

func (ss *CtxSession) SetKV(kvs cst.KV) {
	ss.saved = false
	//if ss.values == nil {
	//	logx.InfoF("%#v", ss)
	//	return
	//}
	for k, v := range kvs {
		ss.values[k] = v
	}
}

func (ss *CtxSession) Save() error {
	// 如果已经保存了，不会重复保存
	if ss.saved == true {
		return nil
	}
	// 调用自定义函数保存当前 session
	_, err := ss.saveSessionToRedis()

	// TODO: 如果保存失败怎么办？目前是抛异常，本次请求直接返回错误。
	if err != nil {
		fst.GFPanic("Save session error. " + err.Error())
	} else {
		ss.saved = true
	}
	return nil
}

func (ss *CtxSession) Saved() bool {
	return ss.saved
}

func (ss *CtxSession) Del(key string) {
	delete(ss.values, key)
	ss.saved = false
}

func (ss *CtxSession) Expire(ttl int32) {
	yn, err := ss.setSessionExpire(ttl)
	if yn == false || err != nil {
		fst.GFPanic("Session expire error.")
	}
}

func (ss *CtxSession) SidIsNew() bool {
	return ss.tokenIsNew
}

func (ss *CtxSession) Sid() string {
	return ss.token
}

func (ss *CtxSession) Destroy() {
	ss.destroySession()
	ss.resetSession()
}

func (ss *CtxSession) Recreate(c *fst.Context) {
	ss.rebuildToken(c)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 新生成一个SDX Session对象，生成新的tok
func (ss *CtxSession) rebuildToken(c *fst.Context) {
	guid, tok := genToken(MySess.Secret + c.ClientIP())
	ss.saved = true // 意味着没有设置值的时候就不需要保存了
	ss.values = make(map[string]any)
	ss.guid = guid
	ss.token = tok
	ss.tokenIsNew = true
}

// 重置session对象
func (ss *CtxSession) resetSession() {
	ss.saved = true
	ss.values = nil
	ss.guid = ""
	ss.token = ""
	ss.tokenIsNew = false
}
