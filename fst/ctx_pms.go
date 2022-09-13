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

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (c *Context) Param(key string) string {
	return c.match.params.ByName(key)
}

// add by sdx on 20210305
// 就当 c.Pms (c.ReqRaw.Form) 中的是 JSON 对象，我们需要用这个数据源绑定任意的对象
func (c *Context) BindPms(dst any) error {
	return mapx.ApplyKVOfData(dst, c.Pms)
}

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
		if err := jsonx.UnmarshalFromReader(&c.Pms, c.ReqRaw.Body); err != nil {
			return err
		}
	//case strings.HasPrefix(ctType, cst.MIMEAppXml), strings.HasPrefix(ctType, cst.MIMEXml):
	//	if err := c.BindXML(&c.Pms); err != nil {
	//		return err
	//	}
	case strings.HasPrefix(ctType, cst.MIMEPostForm), strings.HasPrefix(ctType, cst.MIMEMultiPostForm):
		c.ParseForm()
		urlParsed = true
		applyUrlValue(c.Pms, c.ReqRaw.Form)
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

// A. ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 解析 Url 中的参数
func (c *Context) ParseQuery() {
	if c.queryCache == nil {
		c.queryCache = c.ReqRaw.URL.Query()
	}
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.ReqRaw.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.Query("id") == "1234"
// 	   c.Query("name") == "Manu"
// 	   c.Query("value") == ""
// 	   c.Query("wtf") == ""
func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
//     GET /?name=Manu&lastname=
//     c.DefaultQuery("name", "unknown") == "Manu"
//     c.DefaultQuery("id", "none") == "none"
//     c.DefaultQuery("lastname", "none") == ""
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
// It is shortcut for `c.ReqRaw.URL.Query().Get(key)`
//     GET /?name=Manu&lastname=
//     ("Manu", true) == c.GetQuery("name")
//     ("", false) == c.GetQuery("id")
//     ("", true) == c.GetQuery("lastname")
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.ParseQuery()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// QueryMap returns a map for a given query key.
func (c *Context) QueryMap(key string) map[string]string {
	kvs, _ := c.GetQueryMap(key)
	return kvs
}

// GetQueryMap returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.ParseQuery()
	return c.get(c.queryCache, key)
}

// POST Content-type:
// application/x-www-form-urlencoded：默认的编码方式。但是在用文本的传输和MP3等大型文件的时候，使用这种编码就显得效率低下。
// multipart/form-data：指定传输数据为二进制类型，比如图片、mp3、文件。（注意：这个时候有 boundary=--xxxx 参数）
// text/plain：纯文体的传输。空格转换为 “+” 加号，但不对特殊字符编码。
// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
func (c *Context) DefaultPostForm(key, defaultValue string) string {
	if value, ok := c.GetPostForm(key); ok {
		return value
	}
	return defaultValue
}

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
// For example, during a PATCH request to update the user's email:
//     email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
// 	   email=                  -->  ("", true) := GetPostForm("email") // set email to ""
//                             -->  ("", false) := GetPostForm("email") // do nothing with email
func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) PostFormArray(key string) []string {
	values, _ := c.GetPostFormArray(key)
	return values
}

// B. ++++++++++++++++++++++++++++
// 解析所有 Post 数据到 PostForm对象中，同时将 PostForm 和 QueryForm 中的数据合并到 Form 中。
func (c *Context) ParseForm() {
	if c.formCache == nil {
		//c.formCache = make(url.Values)
		if err := c.ReqRaw.ParseMultipartForm(c.myApp.MaxMultipartMemory); err != nil {
			if err != http.ErrNotMultipart {
				logx.DebugF("error on parse multipart form array: %v", err)
			}
		}
		c.formCache = c.ReqRaw.PostForm
	}
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
	c.ParseForm()
	if values := c.formCache[key]; len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// PostFormMap returns a map for a given form key.
func (c *Context) PostFormMap(key string) map[string]string {
	kvs, _ := c.GetPostFormMap(key)
	return kvs
}

// GetPostFormMap returns a map for a given form key, plus a boolean value
// whether at least one value exists for the given key.
func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	c.ParseForm()
	return c.get(c.formCache, key)
}

// get is an internal method and returns a map which satisfy conditions.
func (c *Context) get(m map[string][]string, key string) (map[string]string, bool) {
	kvs := make(map[string]string)
	exist := false
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				kvs[k[i+1:][:j]] = v[0]
			}
		}
	}
	return kvs, exist
}
