package sdx

import (
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/fst"
)

var (
	sdxTokenPrefix   = "t:"   // token 字符串的 前缀
	sdxSessKeyPrefix = "tls:" // session 的前缀
)

type RedisSessCnf struct {
	RedisConn  gfrds.ConnCnf `v:"required"`                       // 用 Redis 做持久化
	GuidField  string        `v:"def=user_id"`                    // 标记当前登录用户字段是 user_id
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
var MySessDB *SessionDB

// 采用 “闪电侠” session 方案的时候需要先初始化参数
func SetupSession(ss *SessionDB) {
	if MySessDB != nil {
		return
	}
	MySessDB = ss

	if ss.Redis == nil {
		ss.Redis = gfrds.NewGoRedis(&ss.RedisConn)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 默认将使用 Redis 存放 session 信息
type CtxSession struct {
	Values     fst.KV // map[string]interface{}
	Guid       string // redis session key
	Token      string // Sid
	TokenIsNew bool   // Sid is new
	Saved      bool   // Whether it has been saved
}

// CtxSession 需要实现 sessionKeeper 所有接口
var _ fst.SessionKeeper = &CtxSession{}

func (ss *CtxSession) Get(key string) any {
	if ss.Values == nil {
		return nil
	}
	return ss.Values[key]
}

func (ss *CtxSession) Set(key string, val any) {
	ss.Saved = false
	ss.Values[key] = val
}

func (ss *CtxSession) SetKV(kvs fst.KV) {
	ss.Saved = false
	//if ss.Values == nil {
	//	logx.InfoF("%#v", ss)
	//	return
	//}
	for k, v := range kvs {
		ss.Values[k] = v
	}
}

func (ss *CtxSession) Save() {
	// 如果已经保存了，不会重复保存
	if ss.Saved == true {
		return
	}
	// 调用自定义函数保存当前 session
	_, err := ss.saveSessionToRedis()

	// TODO: 如果保存失败怎么办？目前是抛异常，本次请求直接返回错误。
	if err != nil {
		fst.GFPanic("Save session error.")
	} else {
		ss.Saved = true
	}
}

func (ss *CtxSession) Del(key string) {
	delete(ss.Values, key)
	ss.Saved = false
}

func (ss *CtxSession) Expire(ttl int32) {
	yn, err := ss.setSessionExpire(ttl)
	if yn == false || err != nil {
		fst.GFPanic("Session expire error.")
	}
}

func (ss *CtxSession) SidIsNew() bool {
	return ss.TokenIsNew
}

func (ss *CtxSession) Sid() string {
	return ss.Token
}

func (ss *CtxSession) Destroy() {
	ss.destroySession()
}
