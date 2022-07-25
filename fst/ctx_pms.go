// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"net/url"
	"strings"
)

/************************************/
/*********** Context Pms ************/
/************************************/
func (c *Context) ParseRequestData() error {
	// 防止重复解析
	if c.Pms != nil {
		return nil
	}
	c.Pms = make(cst.KV)
	urlParsed := false

	ctType := c.ReqRaw.Header.Get(cst.HeaderContentType)
	switch {
	case strings.HasPrefix(ctType, cst.MIMEAppJson):
		if err := c.BindJSON(&c.Pms); err != nil {
			return err
		}
	case strings.HasPrefix(ctType, cst.MIMEAppXml), strings.HasPrefix(ctType, cst.MIMEXml):
		if err := c.BindXML(&c.Pms); err != nil {
			return err
		}
	case strings.HasPrefix(ctType, cst.MIMEPostForm), strings.HasPrefix(ctType, cst.MIMEMultiPostForm):
		c.ParseForm()
		urlParsed = true
		applyUrlValue(c.Pms, c.ReqRaw.Form)
	default:
	}

	if !urlParsed {
		c.ParseQuery()
		applyUrlValue(c.Pms, c.queryCache)
	}

	return nil
}

// 上传的参数一般都是单一的，不需要 url.Values 中的 slice切片
func applyUrlValue(pms cst.KV, values url.Values) {
	for key, val := range values {
		if len(val) > 1 {
			pms[key] = val
		} else {
			pms[key] = val[0]
		}
	}
}

//// 如果没有匹配路由，需要一些初始化
//func (c *Context) GetPms(key string) (val any, ok bool) {
//	c.mu.RLock()
//	val, ok = c.Pms[key]
//	c.mu.RUnlock()
//	return
//}

//// 启用这个模块之后，gin 的 binding 特性就不能使用了，因为无法读取body内容了。
//func (c *Context) GenPmsByJSONBody() {
//	if c.Pms != nil {
//		return
//	}
//	c.Pms = make(cst.KV)
//	if err := c.BindJSON(&c.Pms); err != nil {
//	}
//
//	c.ParseQuery()
//	for key, val := range c.queryCache {
//		c.Pms[key] = val[0]
//	}
//}
//
//func (c *Context) GenPmsByFormBody() {
//	if c.Pms != nil {
//		return
//	}
//	c.ParseForm()
//	c.Pms = make(cst.KV, len(c.ReqRaw.Form))
//	for key, val := range c.ReqRaw.Form {
//		c.Pms[key] = val[0]
//	}
//}
//
//func (c *Context) GenPmsByXMLBody() {
//	if c.Pms != nil {
//		return
//	}
//	c.Pms = make(cst.KV)
//	if err := c.BindXML(&c.Pms); err != nil {
//	}
//
//	c.ParseQuery()
//	for key, val := range c.queryCache {
//		c.Pms[key] = val[0]
//	}
//}
