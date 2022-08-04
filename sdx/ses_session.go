package sdx

import (
	"github.com/qinchende/gofast/fst"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// TODO: 还原 session ，验证合法性
// 所有请求先经过这里验证 session 信息
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

	tok := c.Pms["tok"]
	if tok == nil {
		tok = ""
	}
	ss.Token = tok.(string)

	// 没有tok，赋予当前请求新的token，同时走后面的逻辑
	if len(ss.Token) < 50 {
		ss.rebuildToken(c)
		return
	}

	// 有 tok ，解析出 [guid、hmac]，实际上 token = [guid].[hmac]
	reqGuid, reqHmac := parseToken(ss.Token)
	// 传了 token 就要检查当前 token 合法性：
	// 1. 不正确，需要分配新的Token。
	// 2. 过期，用当前Token重建Session记录。
	isValid := checkToken(reqGuid, reqHmac, MySessDB.Secret+c.ClientIP())

	// 按照ip计算出当前hmac，和请求中的hmac相比较，看是否相等
	// 如果Guid验证通过
	if isValid || MySessDB.MustKeepIP == false {
		ss.Guid = reqGuid
	}

	// 如果没有Guid，就新生成一个
	if ss.Guid == "" {
		ss.rebuildToken(c)
	} else {
		ss.Values = make(map[string]any)
		ss.loadSessionFromRedis(c) // 通过Guid 到 redis 中获取当前 session
	}
}

// 验证是否登录
func MustLogin(c *fst.Context) {
	uid := c.Sess.Get(MySessDB.GuidField)
	if uid == nil || uid == "" {
		c.AbortHandlers()
		c.Fai(110, "认证失败，请先登录。", nil)
	}
}

// 销毁当前Session
func DestroySession(c *fst.Context) {
	c.Sess.Destroy()
	c.Sess = nil
}

func NewSession(c *fst.Context) {
	ss := new(CtxSession)
	c.Sess = ss

	ss.rebuildToken(c)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 新生成一个SDX Session对象，生成新的tok
func (ss *CtxSession) rebuildToken(c *fst.Context) {
	guid, tok := genToken(MySessDB.Secret + c.ClientIP())
	ss.Values = make(map[string]any)
	ss.Guid = guid
	ss.Token = tok
	ss.TokenIsNew = true
	ss.Saved = false
}
