// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import "github.com/qinchende/gofast/fst"

func PmsParser(c *fst.Context) {
	c.PanicIfErr(c.CollectPms(), "解析请求数据出现错误")
}

// 验证请求是否经过了合法认证
func SessMustLogin(c *fst.Context) {
	if _, ok := c.Sess.Get(MySessDB.UidField); !ok {
		c.AbortFai(110, "User login auth error.", nil)
	}
}

// 验证Session是一个合法的来源
func SessIPValid(c *fst.Context) {
}

//func JwtAuthHandler(secret string) fst.CtxHandler {
//	return mid.JwtAuthHandler(secret)
//}
//func BuildPmsOfJson(ctx *fst.Context) {
//	ctx.GenPmsByJSONBody()
//}
//
//func BuildPmsOfXml(ctx *fst.Context) {
//	ctx.GenPmsByXMLBody()
//}
//
//func BuildPmsOfForm(ctx *fst.Context) {
//	ctx.GenPmsByFormBody()
//}
