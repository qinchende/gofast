// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/mapx"
	"net/http"
	"net/url"
	"strings"
)

// add by sdx on 20210305
// c.Pms 中有提交的所有数据，以KV形式存在。我们需要用这个数据源绑定任意的struct对象
func (c *Context) Bind(dst any) error {
	return mapx.BindKV(dst, c.Pms, mapx.LikeInput)
}

// UrlParam returns the value of the URL param.
// It is a shortcut for c.UrlParams.Value(key)
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.UrlParam("id") // id == "john"
//     })
func (c *Context) UrlParam(key string) string {
	return c.route.params.Value(key)
}

// 必须有指定参数，否则抛异常
func (c *Context) UrlParamMust(key string) string {
	return c.route.params.ValueMust(key)
}

func (c *Context) UrlParamOk(key string) (string, bool) {
	return c.route.params.Get(key)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 解析 Url 中的参数
func (c *Context) QueryValues() url.Values {
	// 单独调用这个还是会解析一下Get请求中携带的URL参数，即使ParseForm已解析了一次URL参数
	val := c.queryCache()
	if val == nil {
		val = c.Req.URL.Query()
		if c.myApp.WebConfig.CacheQueryValues {
			c.setQueryCache(val)
		}
	}
	return val
}

// 解析所有 Post 数据到 PostForm对象中，同时将 PostForm 和 QueryForm 中的数据合并到 Form 中。
func (c *Context) ParseForm() {
	if c.Req.PostForm == nil {
		// 如果解析出错，就当做解析不出参数，参数为空
		maxMemory := c.myApp.WebConfig.MaxMultipartBytes
		if err := c.Req.ParseMultipartForm(maxMemory); err != nil && err != http.ErrNotMultipart {
			logx.DebugF("parse multipart form error: %v", err)
		}
	}
}

// ##这个方法很重要##
// 框架每次都将请求所携带的相关数据解析之后加入统一的变量c.Pms中，这样对开发人员来说只需要关注c.Pms中有无自己想要的数据，
// 至于数据是通过什么形式提交上来的并不那么重要。
// 最常见的就是GET请求URL上的参数，POST请求中req.Body携带的信息
func (c *Context) CollectPms() error {
	// 防止重复解析
	if c.Pms != nil {
		return nil
	}
	c.Pms = make(cst.KV)
	c.Pms2 = c.getPms() // 这里可能是不同数据结果
	urlParsed := false

	ctType := c.Req.Header.Get(cst.HeaderContentType)
	if strings.HasPrefix(ctType, cst.MIMEAppJson) {
		if err := jsonx.UnmarshalFromReader(&c.Pms, c.Req.Body); err != nil {
			return err
		}
		for k := range c.Pms {
			c.Pms2.Set(k, c.Pms[k])
		}
	} else if strings.HasPrefix(ctType, cst.MIMEPostForm) || strings.HasPrefix(ctType, cst.MIMEMultiPostForm) {
		c.ParseForm()
		urlParsed = true
		applyUrlValue(c.Pms, c.Req.Form)
		applyUrlValue(c.Pms2, c.Req.Form)
	}

	// Url中带入的查询参数加入参数字典
	if !urlParsed {
		qu := c.QueryValues()
		applyUrlValue(c.Pms, qu)
		applyUrlValue(c.Pms2, qu)
	}

	// 将UrlParams加入参数字典
	if c.myApp.WebConfig.ApplyUrlParamsToPms && c.route.params != nil {
		kvs := *c.route.params
		for i := range kvs {
			//c.Pms[kvs[i].Key] = kvs[i].Value
			c.Pms.Set(kvs[i].Key, kvs[i].Value)
			c.Pms2.Set(kvs[i].Key, kvs[i].Value)
		}
	}

	return nil
}

// 上传的参数一般都是单一的，不需要 url.Values 中的 slice切片
func applyUrlValue(pms cst.SuperKV, webValues url.Values) {
	for key := range webValues {
		pms.Set(key, webValues[key][0])
	}
}
