// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/fst/render"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"net/http"
)

//  返回结构体
type Ret struct {
	Code int32  // 返回编码
	Msg  string // 文本消息
	Data any    // 携带数据体
	Desc string // 描述，内部说明，不对外传递和显示
}

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

func (c *Context) FaiMsg(msg string) {
	c.Fai(0, msg, nil)
}

func (c *Context) FaiData(data any) {
	c.Fai(0, "", data)
}

func (c *Context) FaiCode(code int32) {
	c.Fai(code, "", nil)
}

func (c *Context) FaiRet(ret *Ret) {
	c.Fai(ret.Code, ret.Msg, ret.Data)
}

func (c *Context) Fai(code int32, msg string, data any) {
	c.kvSucFai(statusFai, code, msg, data)
}

// 简易的抛出异常的方式，终止执行链，返回错误
func (c *Context) FaiPanicIf(yes bool, val any) {
	if !yes {
		return
	}

	switch val.(type) {
	case string:
		panic(cst.GFFaiString(val.(string)))
	case int:
		panic(cst.GFFaiInt(val.(int)))
	case error:
		panic(cst.GFError(val.(error)))
	default:
		str := lang.ToString(val)
		panic(cst.GFFaiString(str))
	}
}

// +++++
func (c *Context) SucMsg(msg string) {
	c.Suc(1, msg, nil)
}

func (c *Context) SucData(data any) {
	c.Suc(1, "", data)
}

func (c *Context) SucCode(code int32) {
	c.Suc(code, "", nil)
}

func (c *Context) SucRet(ret *Ret) {
	c.Suc(ret.Code, ret.Msg, ret.Data)
}

func (c *Context) Suc(code int32, msg string, data any) {
	c.kvSucFai(statusSuc, code, msg, data)
}

func (c *Context) kvSucFai(status string, code int32, msg string, data any) {
	jsonData := cst.KV{
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
// 立即返回错误，跳过后面的执行链
func (c *Context) AbortFai(code int, msg string) {
	bytes, _ := jsonx.Marshal(cst.KV{
		"status": statusFai,
		"code":   code,
		"msg":    msg,
	})
	c.AbortDirect(http.StatusOK, bytes)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) Json(resStatus int, obj any) {
	c.Render(resStatus, render.JSON{Data: obj})
}

// String writes the given string into the response body.
func (c *Context) String(resStatus int, format string, values ...any) {
	c.Render(resStatus, render.Text{Format: format, Data: values})
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
	if c.tryToRender() == false {
		return
	}

	c.ResWrap.WriteHeader(resStatus)
	// 如果指定的返回状态，不能返回数据内容。需要特殊处理
	if !bodyAllowedForStatus(resStatus) {
		r.WriteContentType(c.ResWrap)
		return
	}

	// 返回结果先写入缓存
	if err := r.Write(c.ResWrap); err != nil {
		panic(err)
	}

	if c.route.ptrNode.hasBeforeSend {
		c.execBeforeSendHandlers() // 可以抛出异常，终止 Send data
	}
	if c.Sess != nil {
		_ = c.Sess.Save()
	}
	if _, err := c.ResWrap.Send(); err != nil { // really send response data
		panic(err)
	}
	if c.route.ptrNode.hasAfterSend {
		c.execAfterSendHandlers()
	}
}

// 强行终止处理，立即返回指定结果，不执行Render
func (c *Context) AbortDirect(resStatus int, stream any) {
	if c.tryToRender() == false {
		return
	}
	c.execIdx = maxRouteHandlers
	_ = c.ResWrap.SendHijack(resStatus, lang.ToBytes(stream))
}

func (c *Context) AbortRedirect(resStatus int, redirectUrl string) {
	if c.tryToRender() == false {
		return
	}
	c.execIdx = maxRouteHandlers
	c.ResWrap.SendHijackRedirect(c.ReqRaw, resStatus, redirectUrl)
}

// 这个是为超时返回准备的特殊方法，一般不要使用
func (c *Context) ReturnTimeout(resStatus int, hint any) {
	c.ResWrap.sendByTimeoutGoroutine(resStatus, lang.ToBytes(hint))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

// NOTE: 要避免 double render。只执行第一次Render的结果，后面的Render直接丢弃
func (c *Context) tryToRender() bool {
	c.rwLock.Lock()
	if c.rendered {
		c.rwLock.Unlock()
		logx.Warn("Double render, the call canceled.")
		return false
	}
	c.rendered = true
	c.rwLock.Unlock()
	return true
}
