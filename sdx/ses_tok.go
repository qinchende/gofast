// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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
	ss.token, _ = c.GetString(PmsToken)

	// 没有tok，赋予当前请求新的token，同时走后面的逻辑
	if ss.token == "" {
		ss.rebuildToken(c)
		return
	}

	// 有 tok ，解析出 [guid、hmac]，实际上 token = [guid].[hmac]
	reqGuid, reqHmac := parseToken(ss.token)
	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确，需要分配新的Token。
	// 2. 过期，用当前Token重建Session记录。
	isValid := checkToken(reqGuid, reqHmac, MySessDB.Secret+c.ClientIP())

	// 按照ip计算出当前hmac，和请求中的hmac相比较，看是否相等
	// 如果Guid验证通过
	if isValid || MySessDB.MustKeepIP == false {
		ss.guid = reqGuid
	}

	// 如果没有Guid，就新生成一个
	if ss.guid == "" {
		ss.rebuildToken(c)
	} else {
		ss.values = make(cst.WebKV)
		if err := ss.loadSessionFromRedis(); err != nil {
			c.CarryMsg(err.Error())
			c.AbortFai(110, "Load session data from redis error.", nil)
		}
	}
}

//// 生成新的 token 信息
//func SessRecreate(c *fst.Context) {
//	ss := new(TokSession)
//	ss.rebuildToken(c)
//	c.Sess = ss
//}

// TokSession
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 默认将使用 Redis 存放 session 信息
// TODO: 注意，这个实现是非线程安全的
type TokSession struct {
	values     cst.WebKV // map[string]string
	guid       string    // redis session key
	token      string    // Sid
	tokenIsNew bool      // Sid is new
	changed    bool      // Session values changed
	saved      bool      // Whether it has been saved
}

// TokSession 需要实现 sessionKeeper 所有接口
var _ fst.SessionKeeper = &TokSession{}

func (ss *TokSession) Get(key string) (v string, ok bool) {
	v, ok = ss.values[key]
	return
}

func (ss *TokSession) GetValues() cst.WebKV {
	return ss.values
}

func (ss *TokSession) Set(key string, val string) {
	ss.saved = false
	ss.values[key] = val
}

func (ss *TokSession) SetValues(kvs cst.WebKV) {
	ss.saved = false
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
	if ss.saved == true {
		return
	}

	// 调用自定义函数保存当前session。保存失败就抛异常
	if _, err := ss.saveSessionToRedis(); err != nil {
		cst.Panic("Save session error. " + err.Error())
	} else {
		ss.saved = true
	}
}

//func (ss *TokSession) Saved() bool {
//	return ss.saved
//}

func (ss *TokSession) Del(key string) {
	delete(ss.values, key)
	ss.saved = false
}

func (ss *TokSession) ExpireS(ttl int32) {
	yn, err := ss.setSessionExpire(ttl)
	if yn == false || err != nil {
		cst.Panic("Session expire error.")
	}
}

func (ss *TokSession) TokenIsNew() bool {
	return ss.tokenIsNew
}

func (ss *TokSession) Token() string {
	return ss.token
}

func (ss *TokSession) Destroy() {
	if err := ss.destroySession(); err != nil {
		cst.Panic("Destroy session error. " + err.Error())
	}
	ss.resetSession()
}

func (ss *TokSession) Recreate(c *fst.Context) {
	ss.rebuildToken(c)
}

// 新生成一个SDX Session对象，生成新的tok
func (ss *TokSession) rebuildToken(c *fst.Context) {
	guid, tok := genToken(MySessDB.Secret + c.ClientIP())
	ss.saved = true // 意味着没有设置值的时候就不需要保存了
	ss.values = make(cst.WebKV, 8)
	ss.guid = guid
	ss.token = tok
	ss.tokenIsNew = true
}

// 重置session对象
func (ss *TokSession) resetSession() {
	ss.saved = true
	ss.values = nil
	ss.guid = ""
	ss.token = ""
	ss.tokenIsNew = false
}
