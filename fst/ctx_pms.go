// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
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
	c.Pms = make(map[string]interface{})
	isForm := false

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
		isForm = true
		for key, val := range c.ReqRaw.Form {
			c.Pms[key] = val[0]
		}
	default:
	}

	if !isForm {
		c.ParseQuery()
		for key, val := range c.queryCache {
			c.Pms[key] = val[0]
		}
	}

	return nil
}

// 启用这个模块之后，gin 的 binding 特性就不能使用了，因为无法读取body内容了。
func (c *Context) GenPmsByJSONBody() {
	if c.Pms != nil {
		return
	}
	c.Pms = make(map[string]interface{})
	if err := c.BindJSON(&c.Pms); err != nil {
	}

	c.ParseQuery()
	for key, val := range c.queryCache {
		c.Pms[key] = val[0]
	}
}

func (c *Context) GenPmsByFormBody() {
	if c.Pms != nil {
		return
	}
	c.ParseForm()
	c.Pms = make(map[string]interface{}, len(c.ReqRaw.Form))
	for key, val := range c.ReqRaw.Form {
		c.Pms[key] = val[0]
	}
}

func (c *Context) GenPmsByXMLBody() {
	if c.Pms != nil {
		return
	}
	c.Pms = make(map[string]interface{})
	if err := c.BindXML(&c.Pms); err != nil {
	}

	c.ParseQuery()
	for key, val := range c.queryCache {
		c.Pms[key] = val[0]
	}
}

// 如果没有匹配路由，需要一些初始化
func (c *Context) GetPms(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Pms[key]
	c.mu.RUnlock()
	return
}
