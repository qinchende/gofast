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
func (c *Context) ParseQuery() {
	// 单独调用这个还是会解析一下Get请求中携带的URL参数，即使ParseForm已解析了一次URL参数
	if c.queryCache == nil {
		c.queryCache = c.Req.URL.Query()
	}
}

// 解析所有 Post 数据到 PostForm对象中，同时将 PostForm 和 QueryForm 中的数据合并到 Form 中。
func (c *Context) ParseForm() {
	if c.formCache == nil {
		// 如果解析出错，就当做解析不出参数，参数为空
		maxMemory := c.myApp.WebConfig.MaxMultipartBytes
		if err := c.Req.ParseMultipartForm(maxMemory); err != nil && err != http.ErrNotMultipart {
			logx.DebugF("parse multipart form error: %v", err)
		}
		c.formCache = c.Req.PostForm
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
	c.Pms = c.getPms()
	urlParsed := false

	ctType := c.Req.Header.Get(cst.HeaderContentType)
	if strings.HasPrefix(ctType, cst.MIMEAppJson) {
		if err := jsonx.UnmarshalFromReader(&c.Pms, c.Req.Body); err != nil {
			return err
		}
	} else if strings.HasPrefix(ctType, cst.MIMEPostForm) || strings.HasPrefix(ctType, cst.MIMEMultiPostForm) {
		c.ParseForm()
		urlParsed = true
		applyUrlValue(c.Pms, c.Req.Form)
	}
	//else if strings.HasPrefix(ctType, cst.MIMEAppXml) || strings.HasPrefix(ctType, cst.MIMEXml) {
	//	if err := c.BindXML(&c.Pms); err != nil {
	//		return err
	//	}
	//}

	// Url中带入的查询参数加入参数字典
	if !urlParsed {
		c.ParseQuery()
		applyUrlValue(c.Pms, c.queryCache)
	}

	// 将UrlParams加入参数字典
	if c.myApp.WebConfig.ApplyUrlParamsToPms && c.route.params != nil {
		kvs := *c.route.params
		for i := range kvs {
			c.Pms[kvs[i].Key] = kvs[i].Value
		}
	}

	return nil
}

// 上传的参数一般都是单一的，不需要 url.Values 中的 slice切片
func applyUrlValue(pms cst.KV, webValues url.Values) {
	for key := range webValues {
		pms[key] = webValues[key][0]
	}
}

//// A. ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 解析 Url 中的参数
//func (c *Context) ParseQuery() {
//	// 单独调用这个还是会解析一下Get请求中携带的URL参数，即使ParseForm已解析了一次URL参数
//	if c.queryCache == nil {
//		c.queryCache = c.Req.URL.Query()
//	}
//}

//
//// Query returns the keyed url query value if it exists,
//// otherwise it returns an empty string `("")`.
//// It is shortcut for `c.Req.URL.Query().Get(key)`
////     GET /path?id=1234&name=Manu&value=
//// 	   c.Query("id") == "1234"
//// 	   c.Query("name") == "Manu"
//// 	   c.Query("value") == ""
//// 	   c.Query("wtf") == ""
//func (c *Context) Query(key string) string {
//	value, _ := c.Query2(key)
//	return value
//}
//
//// DefaultQuery returns the keyed url query value if it exists,
//// otherwise it returns the specified defaultValue string.
//// See: Query() and GetQuery() for further information.
////     GET /?name=Manu&lastname=
////     c.DefaultQuery("name", "unknown") == "Manu"
////     c.DefaultQuery("id", "none") == "none"
////     c.DefaultQuery("lastname", "none") == ""
//func (c *Context) QueryDef(key, defaultValue string) string {
//	if value, ok := c.Query2(key); ok {
//		return value
//	}
//	return defaultValue
//}
//
//// GetQuery is like Query(), it returns the keyed url query value
//// if it exists `(value, true)` (even when the value is an empty string),
//// otherwise it returns `("", false)`.
//// It is shortcut for `c.Req.URL.Query().Get(key)`
////     GET /?name=Manu&lastname=
////     ("Manu", true) == c.GetQuery("name")
////     ("", false) == c.GetQuery("id")
////     ("", true) == c.GetQuery("lastname")
//func (c *Context) Query2(key string) (string, bool) {
//	if values, ok := c.QueryArray2(key); ok {
//		return values[0], ok
//	}
//	return "", false
//}
//
//// QueryArray returns a slice of strings for a given query key.
//// The length of the slice depends on the number of params with the given key.
//func (c *Context) QueryArray(key string) []string {
//	values, _ := c.QueryArray2(key)
//	return values
//}
//
//// GetQueryArray returns a slice of strings for a given query key, plus
//// a boolean value whether at least one value exists for the given key.
//func (c *Context) QueryArray2(key string) ([]string, bool) {
//	c.ParseQuery()
//	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
//		return values, true
//	}
//	return []string{}, false
//}
//
//// QueryMap returns a map for a given query key.
//func (c *Context) QueryMap(key string) map[string]string {
//	kvs, _ := c.QueryMap2(key)
//	return kvs
//}
//
//// GetQueryMap returns a map for a given query key, plus a boolean value
//// whether at least one value exists for the given key.
//func (c *Context) QueryMap2(key string) (map[string]string, bool) {
//	c.ParseQuery()
//	return c.get(c.queryCache, key)
//}

//// B. ++++++++++++++++++++++++++++
//// 解析所有 Post 数据到 PostForm对象中，同时将 PostForm 和 QueryForm 中的数据合并到 Form 中。
//func (c *Context) ParseForm() {
//	if c.formCache == nil {
//		if err := c.Req.ParseMultipartForm(c.myApp.WebConfig.MaxMultipartBytes); err != nil && err != http.ErrNotMultipart {
//			logx.DebugF("error on parse multipart form array: %v", err)
//		}
//		c.formCache = c.Req.PostForm
//	}
//}

//
//// POST Content-type:
//// application/x-www-form-urlencoded：默认的编码方式。但是在用文本的传输和MP3等大型文件的时候，使用这种编码就显得效率低下。
//// multipart/form-data：指定传输数据为二进制类型，比如图片、mp3、文件。（注意：这个时候有 boundary=--xxxx 参数）
//// text/plain：纯文体的传输。空格转换为 “+” 加号，但不对特殊字符编码。
//// PostForm returns the specified key from a POST urlencoded form or multipart form
//// when it exists, otherwise it returns an empty string `("")`.
//func (c *Context) PostForm(key string) string {
//	value, _ := c.PostForm2(key)
//	return value
//}
//
//// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
//// when it exists, otherwise it returns the specified defaultValue string.
//// See: PostForm() and GetPostForm() for further information.
//func (c *Context) PostFormDef(key, defaultValue string) string {
//	if value, ok := c.PostForm2(key); ok {
//		return value
//	}
//	return defaultValue
//}
//
//// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
//// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
//// otherwise it returns ("", false).
//// For example, during a PATCH request to update the user's email:
////     email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
//// 	   email=                  -->  ("", true) := GetPostForm("email") // set email to ""
////                             -->  ("", false) := GetPostForm("email") // do nothing with email
//func (c *Context) PostForm2(key string) (string, bool) {
//	if values, ok := c.PostFormArray2(key); ok {
//		return values[0], ok
//	}
//	return "", false
//}
//
//// PostFormArray returns a slice of strings for a given form key.
//// The length of the slice depends on the number of params with the given key.
//func (c *Context) PostFormArray(key string) []string {
//	values, _ := c.PostFormArray2(key)
//	return values
//}
//
//// GetPostFormArray returns a slice of strings for a given form key, plus
//// a boolean value whether at least one value exists for the given key.
//func (c *Context) PostFormArray2(key string) ([]string, bool) {
//	c.ParseForm()
//	if values := c.formCache[key]; len(values) > 0 {
//		return values, true
//	}
//	return []string{}, false
//}
//
//// PostFormMap returns a map for a given form key.
//func (c *Context) PostFormMap(key string) map[string]string {
//	kvs, _ := c.PostFormMap2(key)
//	return kvs
//}
//
//// GetPostFormMap returns a map for a given form key, plus a boolean value
//// whether at least one value exists for the given key.
//func (c *Context) PostFormMap2(key string) (map[string]string, bool) {
//	c.ParseForm()
//	return c.get(c.formCache, key)
//}
//
//// C. ++++++++++++++++++++++++++++
//// get is an internal method and returns a map which satisfy conditions.
//func (c *Context) get(m map[string][]string, key string) (map[string]string, bool) {
//	kvs := make(map[string]string)
//	exist := false
//	for k, v := range m {
//		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
//			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
//				exist = true
//				kvs[k[i+1:][:j]] = v[0]
//			}
//		}
//	}
//	return kvs, exist
//}
