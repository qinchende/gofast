// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"encoding/base64"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/randx"
	"github.com/qinchende/gofast/store/jde"
	"time"
)

// 还原 session ，验证合法性 （所有请求先经过这里验证 session 信息）
// 每一次的访问，都必须要有一个 token ，没有token的 访问将视为 非法.
// 第一次没有 token 的情况下，默认造一个 token
func TokSessBuilder(c *fst.Context) {
	// 不可重复执行 token 检查，Sess构造的过程
	if c.Sess != nil {
		return
	}

	// 每个请求对应的SESSION对象都是新创建的，线程安全。
	ss := new(TokSession)
	c.Sess = ss
	ss.raw, _ = c.GetString(PmsToken)

	// 没有tok，赋予当前请求新的token，同时走后面的逻辑
	if ss.raw == "" {
		ss.createNewToken()
		return
	}

	// 有 tok ，解析出 [guid、hmac]，实际上 token = [guid].[hmac]
	reqGuid, reqHmac := parseToken(ss.raw)
	if reqGuid == "" {
		ss.createNewToken()
		return
	}

	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确，需要分配新的Token。
	// 2. 正确或者过期，利用当前sid重建Session记录。
	isValid := checkToken(reqGuid, MySessDB.Secret, reqHmac)
	if !isValid {
		ss.createNewToken()
		return
	}

	ss.guid = reqGuid
	if err := ss.loadSessFromRedis(); err != nil {
		c.CarryMsg(err.Error())
		c.AbortFai(110, "Load session data from redis error.", nil)
	}
}

// TokSession
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 默认将使用 Redis 存放 session 信息
// TODO: 注意，这个实现是非线程安全的
type TokSession struct {
	raw        string    // raw token
	guid       string    // unique key
	values     cst.WebKV // map[string]string
	ttl        uint32    // 过期时间，秒
	tokenIsNew bool      // just new token?
	changed    bool      // 值是否改变
	saved      bool      // Whether it has been saved
}

// TokSession 需要实现 sessionKeeper 所有接口
var _ fst.SessionKeeper = &TokSession{}
var _TokSessionInitializer TokSession

func (ss *TokSession) Get(key string) (v string, ok bool) {
	v, ok = ss.values[key]
	return
}

func (ss *TokSession) GetValues() cst.WebKV {
	return ss.values
}

func (ss *TokSession) Set(key string, val string) {
	if ss.values == nil {
		ss.values = make(cst.WebKV)
	}
	ss.changed = true
	ss.values[key] = val
}

func (ss *TokSession) SetValues(kvs cst.WebKV) {
	if ss.values == nil {
		ss.values = make(cst.WebKV)
	}
	ss.changed = true
	for k, v := range kvs {
		ss.values[k] = v
	}
}

func (ss *TokSession) SetUid(uid string) {
	ss.Set(MySessDB.UidField, uid)
}

func (ss *TokSession) GetUid() (uid string) {
	uid, _ = ss.Get(MySessDB.UidField)
	return
}

func (ss *TokSession) Save() {
	// 如果已经保存了，不会重复保存
	if ss.saved == true && ss.changed == false {
		return
	}
	// 调用自定义函数保存当前session。保存失败就抛异常
	if err := ss.saveSessToRedis(); err != nil {
		cst.Panic("Save session error. " + err.Error())
	} else {
		ss.saved = true
		ss.changed = false
	}
}

func (ss *TokSession) Del(key string) {
	delete(ss.values, key)
	ss.changed = true
}

// 注意ttl参数的单位是：秒
func (ss *TokSession) ExpireS(ttl uint32) {
	ss.ttl = ttl
	ss.saved = false
}

func (ss *TokSession) TokenIsNew() bool {
	return ss.tokenIsNew
}

func (ss *TokSession) Token() string {
	return ss.raw
}

func (ss *TokSession) Destroy() {
	if err := ss.delSessRedis(); err != nil {
		cst.Panic("Destroy session error. " + err.Error())
	}
	*ss = _TokSessionInitializer
}

func (ss *TokSession) Recreate() {
	ss.createNewToken()
}

// 新生成一个SDX Session对象，生成新的tok
func (ss *TokSession) createNewToken() {
	ss.guid, ss.raw = genToken(MySessDB.Secret)
	ss.tokenIsNew = true
	// 意味着没有设置值的时候不需要保存，新的token传给前端即可
	// 可减轻大量首次请求，保存占用大量数据库资源
	ss.saved = true
	ss.changed = false
}

// crypto
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// raw token = [guid].[hmac]
// hmac is just md5 hash value
// tok=YXRJT0l5ckpYNldBTjYzNHZw.aSeGhf7Nhar08YcyuQmlgw
// 解析不到，都将返回空字符串
func parseToken(tok string) (string, string) {
	dot := int(MySessDB.SidSize)
	// 格式明显不对，直接返回空
	if len(tok) < dot || tok[dot] != '.' {
		return "", ""
	}
	return tok[:dot], tok[(dot + 1):]
}

// 利用当前 guid 和 c 中包含的 request_ip | 计算出hmac值，然后和token中携带的 hmac值比较，来得出合法性
func checkToken(guid, secret, sHmac string) bool {
	md5Val := md5Base64(lang.STB(guid), lang.STB(secret))
	return sHmac == md5Val
}

// YXRJT0l5ckpYNldBTjYzNHZw.aSeGhf7Nhar08YcyuQmlgw
// 闪电侠Guid
func genToken(secret string) (string, string) {
	size := int((MySessDB.SidSize*3 + 3) / 4)
	// TODO：要想办法保证sid的唯一性
	sid := randx.RandomBytes(size)
	guid := base64.RawURLEncoding.EncodeToString(sid)
	guid = guid[:MySessDB.SidSize] // 要确保guid的长度一致性

	md5Val := md5Base64(lang.STB(guid), lang.STB(secret))
	return guid, guid + "." + md5Val
}

// Redis
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 从 redis 中获取 当前 请求上下文的 session data.
func (ss *TokSession) loadSessFromRedis() error {
	str, err := MySessDB.Redis.Get(ss.redisKey())
	if str == "" || err != nil {
		// Note: 此时Redis为空，或者没这个Key，或者数据库连接错误。后面默认是不需要再保存的
		ss.saved = true
		return nil
	}
	if ss.values == nil {
		ss.values = make(cst.WebKV)
	}
	return jde.DecodeBytes(&ss.values, lang.STB(str))
}

func (ss *TokSession) saveSessToRedis() (err error) {
	ttl := MySessDB.TTL
	if ss.ttl > 0 {
		ttl = time.Duration(ss.ttl) * time.Second
	}

	if ss.changed {
		str := ""
		if str, err = jde.EncodeToString(ss.values); err == nil {
			_, err = MySessDB.Redis.Set(ss.redisKey(), str, ttl)
		}
	} else {
		_, err = MySessDB.Redis.Expire(ss.redisKey(), ttl)
	}
	return
}

func (ss *TokSession) delSessRedis() (err error) {
	_, err = MySessDB.Redis.Del(ss.redisKey())
	return
}

func (ss *TokSession) redisKey() string {
	if len(MySessDB.PrefixSessKey) == 0 {
		return ss.guid
	}
	return MySessDB.PrefixSessKey + ss.guid
}
