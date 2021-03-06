// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/cst"
	"strings"
	"time"
)

/************************************/
/*********** Context Pms ************/
/************************************/
func (c *Context) ParseRequestData() {
	if c.Pms != nil {
		return
	}
	c.Pms = make(map[string]interface{})
	isForm := false

	ctType := c.ReqRaw.Header.Get(cst.HeaderContentType)
	switch {
	case strings.HasPrefix(ctType, MIMEJSON):
		if err := c.BindJSON(&c.Pms); err != nil {
		}
	case strings.HasPrefix(ctType, MIMEAppXML), strings.HasPrefix(ctType, MIMETextXML):
		if err := c.BindXML(&c.Pms); err != nil {
		}
	case strings.HasPrefix(ctType, MIMEPOSTForm), strings.HasPrefix(ctType, MIMEMultiPOSTForm):
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

	return
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

/************************************/
/******* metadata management ********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

/************************************/
/********* error management *********/
/************************************/

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	c.Errors = append(c.Errors, parsedError)
	return parsedError
}

/************************************/
/***** golang.org/x/net/context *****/
/************************************/

// Deadline always returns that there is no deadline (ok==false),
// maybe you want to use Req.Context().Deadline() instead.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil (chan which will wait forever),
// if you want to abort your work when the connection was closed
// you should use Req.Context().Done() instead.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err always returns nil, maybe you want to use Req.Context().Err() instead.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.ReqRaw
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
