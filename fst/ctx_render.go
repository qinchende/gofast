// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/render"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/stringx"
	"net/http"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast JSON render
// JSON是GoFast默认的返回格式，一等公民。所以默认函数命名没有给出JSON字样

const (
	statusSuc string = "suc"
	statusFai string = "fai"
)

func (c *Context) FaiErr(err error) {
	c.Fai(0, err.Error(), nil)
}

func (c *Context) FaiStr(msg string) {
	c.Fai(0, msg, nil)
}

func (c *Context) FaiKV(data KV) {
	c.Fai(0, "", data)
}

func (c *Context) Fai(code int32, msg string, data any) {
	c.kvSucFai(statusFai, code, msg, data)
}

// +++++
func (c *Context) SucStr(msg string) {
	c.Suc(0, msg, nil)
}

func (c *Context) SucKV(data KV) {
	c.Suc(0, "", data)
}

func (c *Context) Suc(code int32, msg string, data any) {
	c.kvSucFai(statusSuc, code, msg, data)
}

func (c *Context) kvSucFai(status string, code int32, msg string, data any) {
	jsonData := KV{
		"status": status,
		"code":   code,
		"msg":    msg,
	}
	if data != nil {
		jsonData["data"] = data
	}

	if c.Sess != nil && c.Sess.SidIsNew() {
		jsonData["tok"] = c.Sess.Sid()
	}

	c.Json(http.StatusOK, jsonData)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Abort系列函数都将终止当前 handlers 的执行

func (c *Context) AbortHandlers() {
	c.execIdx = maxRouteHandlers
}

// 自定义返回结果和状态
func (c *Context) AbortFaiStr(msg string) {
	bytes, _ := jsonx.Marshal(KV{
		"status": statusFai,
		"code":   -1,
		"msg":    msg,
	})
	c.AbortDirectBytes(http.StatusOK, bytes)
}

// 强行终止处理，返回指定结果，不执行Render
func (c *Context) AbortDirect(resStatus int, msg string) {
	c.execIdx = maxRouteHandlers
	_ = c.ResWrap.SendHijack(resStatus, stringx.StringToBytes(msg))
}

func (c *Context) AbortDirectBytes(resStatus int, bytes []byte) {
	c.execIdx = maxRouteHandlers
	_ = c.ResWrap.SendHijack(resStatus, bytes)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) Json(resStatus int, obj any) {
	c.Render(resStatus, render.JSON{Data: obj})
}

// String writes the given string into the response body.
func (c *Context) String(resStatus int, format string, values ...any) {
	c.Render(resStatus, render.String{Format: format, Data: values})
}

// File writes the specified file into the body stream in a efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.ResWrap, c.ReqRaw, filepath)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Render writes the response headers and calls render.Render to render data.
// 返回数据的接口
// 可以自定义扩展自己需要的Render
func (c *Context) Render(resStatus int, r render.Render) {
	// NOTE: 要避免 double render。只执行第一次Render的结果，后面的Render直接丢弃
	if c.rendered {
		logx.Warn("Double render, the call canceled.")
		return
	}
	// Render之前加入对应的 render 数据
	c.rendered = true

	// add preSend & afterSend events by sdx on 2021.01.06
	c.execPreSendHandlers()

	c.ResWrap.WriteHeader(resStatus)
	// 如果指定的返回状态，不能返回数据内容。需要特殊处理
	if !bodyAllowedForStatus(resStatus) {
		r.WriteContentType(c.ResWrap)
		return
	}

	// TODO: render之前，统一保存 session
	if c.Sess != nil {
		c.Sess.Save()
	}

	// 返回结果先写入缓存
	if err := r.Write(c.ResWrap); err != nil {
		panic(err)
	}
	// really send response data
	if _, err := c.ResWrap.Send(); err != nil {
		panic(err)
	}

	// add preSend & afterSend events by sdx on 2021.01.06
	c.execAfterSendHandlers()

	// NOTE(by chende 2021.10.28): 下面的这一堆代码需要删除掉，否则中间件不会往下执行了。
	// 到这里其实也意味着调用链到这里就中断了。不需要再执行其它处理函数。
	// 调用链是：[before(s)->handler(s)->after(s)]其中任何地方执行了Render，后面的函数都将不再调用。
	// 但是 preSend 和 afterSend 函数将继续执行。
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
