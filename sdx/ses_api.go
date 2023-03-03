// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/lang"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 还原 session ，验证合法性 （所有请求先经过这里验证 session 信息）
// 每一次的访问，都必须要有一个 token ，没有token的 访问将视为 非法.
// 第一次没有 token 的情况下，默认造一个 token
func SessBuilder(c *fst.Context) {
	// 不可重复执行 token 检查，Sess构造的过程
	if c.Sess != nil {
		return
	}

	// 每个请求对应的SESSION对象都是新创建的，线程安全。
	ss := new(CtxSession)
	c.Sess = ss
	ss.token, _ = c.GetString("tok")

	// 没有tok，赋予当前请求新的token，同时走后面的逻辑
	if len(ss.token) < 50 {
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
		ss.values = make(cst.KV)
		if err := ss.loadSessionFromRedis(c); err != nil {
			c.CarryAddMsg(err.Error())
			c.AbortFai(110, "Load session data from redis error.", nil)
		}
	}
}

// 验证请求是否经过了合法认证
func SessMustLogin(c *fst.Context) {
	uid := c.Sess.Get(MySessDB.UidField)
	if uid == nil || uid == "" {
		c.AbortFai(110, "User login auth error.", nil)
	}
}

// 设置当前登录账号的 uid
func SessSetUid(c *fst.Context, uid int64) {
	c.Sess.Set(MySessDB.UidField, uid)
}

// 获取登录账号的uid
func SessGetUid(c *fst.Context) (uid int64) {
	uid, _ = lang.ToInt64(c.Sess.Get(MySessDB.UidField))
	return
}

// 销毁当前 Session
func SessDestroy(c *fst.Context) {
	c.Sess.Destroy()
	c.Sess = nil
}

// 生成新的 token 信息
func SessRecreate(c *fst.Context) {
	ss := new(CtxSession)
	ss.rebuildToken(c)
	c.Sess = ss
}
