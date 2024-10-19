// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/aid/jsonx"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/fst/httpx"
	"github.com/qinchende/gofast/store/bind"
	"net/http"
	"strings"
)

// UrlParam returns the value of the URL param.
//
//	app.Get("/user/:id", func(c *fst.Context) {
//	    // a GET request to /user/chende
//	    id := c.UrlParam("id") // id == "chende"
//	})
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

// 标准库解法
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 解析 Url 中的参数
func (c *Context) QueryValues() cst.WebKV {
	// 单独调用这个还是会解析一下Get请求中携带的URL参数，即使ParseForm已解析了一次URL参数
	val := c.queryCache()
	if val == nil {
		val = make(cst.WebKV)
		httpx.ParseQuery(val, c.Req.Raw.URL.RawQuery)
		if c.app.WebConfig.CacheQueryValues {
			c.setQueryCache(val)
		}
	}
	return val
}

// 解析所有 Post 数据到 PostForm对象中，同时将 PostForm 和 QueryForm 中的数据合并到 Form 中。
func (c *Context) ParseForm() error {
	if c.Req.Raw.PostForm == nil {
		// 如果解析出错，就当做解析不出参数，参数为空
		maxMemory := c.app.WebConfig.MaxMultipartBytes
		if err := c.Req.Raw.ParseMultipartForm(maxMemory); err != http.ErrNotMultipart {
			return err
		}
	}
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// add by sdx on 20210305
// c.Pms 中有提交的所有数据，以KV形式存在。我们需要用这个数据源绑定任意的struct对象
func (c *Context) Bind(dst any) error {
	return bind.BindKVX(dst, c.Pms, pBindOptions)
}

func (c *Context) BindAndValid(dst any) error {
	return bind.BindKVX(dst, c.Pms, pBindAndValidOptions)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// ##这个方法很重要##
// 框架每次都将请求所携带的相关数据解析之后加入统一的变量c.Pms中，这样对开发人员来说只需要关注c.Pms中有无自己想要的数据，
// 至于数据是通过什么形式提交上来的并不那么重要。
// 最常见的就是 GET中 Url-Query-Params + POST中 req.Body 携带的数据
func (c *Context) CollectPms() error {
	// 防止重复解析
	if c.Pms != nil {
		return nil
	}
	c.newPms() // 实现了cst.SuperKV的类型都可以

	urlParsed := false

	// Body bytes data [JSON or Form]
	ctType := c.Req.Raw.Header.Get(cst.HeaderContentType)
	if strings.HasPrefix(ctType, cst.MIMEAppJson) {
		// +++ JSON格式（可以解析GsonRows数据）
		// 这里可以将复杂的JSON格式数据，直接解析成对象
		if err := jsonx.DecodeRequest(c.Pms, c.Req.Raw); err != nil {
			return err
		}
	} else if strings.HasPrefix(ctType, cst.MIMEPostForm) || strings.HasPrefix(ctType, cst.MIMEMultiPostForm) {
		// +++ Form表单 或者 文件上传
		maxMemory := c.app.WebConfig.MaxMultipartBytes
		if err := httpx.ParseMultipartForm(c.Pms, c.Req.Raw, ctType, maxMemory); err != nil {
			return err
		}
		urlParsed = true
	}

	// Url query params
	if !urlParsed && len(c.Req.Raw.URL.RawQuery) > 0 {
		if c.app.WebConfig.CacheQueryValues {
			kvs := c.QueryValues()
			for key := range kvs {
				c.Pms.SetString(key, kvs[key])
			}
		} else {
			httpx.ParseQuery(c.Pms, c.Req.Raw.URL.RawQuery)
		}
	}

	// Url pattern matching params
	if c.app.WebConfig.ApplyUrlParamsToPms && len(*c.route.params) > 0 {
		kvs := *c.route.params
		for i := range kvs {
			c.Pms.SetString(kvs[i].Key, kvs[i].Value)
		}
	}

	// Note: 是否加入http协议中的 headers 参数？
	// 个人不喜欢，也不推荐用http header的方式，传递业务数据。有啥好处呢，欺骗防火墙？掩耳盗铃？
	// 虽然 http2 & http3 能够实现header的动态表压缩，但对应用来说，从大量header找到自己想要的业务数据并不高效
	// 头信息多了，会无形中增加net/http标准库的资源消耗
	// **如果现实需要，可以自定义中间件加以整合处理**

	return nil
}
