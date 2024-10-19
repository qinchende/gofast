// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/aid/jsonx"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/lang"
	"github.com/qinchende/gofast/fst/render"
	"net/http"
	"strings"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast JSON render
// JSON是GoFast默认的返回格式，一等公民。所以默认函数命名没有给出JSON字样

const (
	statusSuc  = "suc"
	statusFai  = "fai"
	dataField  = "data"
	tokenField = "tok"
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

func (c *Context) FaiCode(code int) {
	c.Fai(code, "", nil)
}

func (c *Context) FaiRet(ret *cst.Ret) {
	c.Fai(ret.Code, ret.Msg, ret.Data)
}

func (c *Context) Fai(code int, msg string, data any) {
	c.kvSucFai(statusFai, code, msg, data)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 自定义三目运算，返货成功或失败信息
func (c *Context) IfSucFai(ifTrue bool, suc, fai any) {
	if ifTrue {
		switch suc.(type) {
		case string:
			c.SucMsg(suc.(string))
		case *cst.Ret:
			c.SucRet(suc.(*cst.Ret))
		case int:
			c.SucCode(suc.(int))
		default:
			c.SucData(suc)
		}
	} else {
		switch fai.(type) {
		case string:
			c.FaiMsg(fai.(string))
		case *cst.Ret:
			c.FaiRet(fai.(*cst.Ret))
		case error:
			c.FaiErr(fai.(error))
		case int:
			c.FaiCode(fai.(int))
		default:
			c.FaiData(fai)
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// +++++
func (c *Context) SucMsg(msg string) {
	c.Suc(1, msg, nil)
}

func (c *Context) SucData(data any) {
	c.Suc(1, "", data)
}

func (c *Context) SucCode(code int) {
	c.Suc(code, "", nil)
}

func (c *Context) SucRet(ret *cst.Ret) {
	c.Suc(ret.Code, ret.Msg, ret.Data)
}

func (c *Context) Suc(code int, msg string, data any) {
	c.kvSucFai(statusSuc, code, msg, data)
}

func (c *Context) kvSucFai(status string, code int, msg string, data any) {
	jsonData := cst.KV{
		"status": status,
		"code":   code,
		"msg":    msg,
	}
	if data != nil {
		jsonData[dataField] = data
	}

	if c.Sess != nil && c.Sess.TokenIsNew() {
		jsonData[tokenField] = c.Sess.Token()
	}

	c.Json(http.StatusOK, jsonData)
}

// 主要用于JSON片段值 val 的返回，返回之前加入通用参数
func (c *Context) SucJsonPart(key, val string) {
	var buf strings.Builder
	buf.Grow(128)

	buf.WriteString(`{"status":"`)
	buf.WriteString(statusSuc)
	buf.WriteString(`","code":1,`)
	if c.Sess != nil && c.Sess.TokenIsNew() {
		buf.WriteByte('"')
		buf.WriteString(tokenField)
		buf.WriteString(`":"`)
		buf.WriteString(c.Sess.Token())
		buf.WriteString(`",`)
	}
	buf.WriteString(`"`)
	buf.WriteString(key)
	buf.WriteString(`":`)
	buf.WriteString(val)
	buf.WriteByte('}')

	c.String(http.StatusOK, buf.String())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Abort系列函数都将终止当前 handlers 的执行
// 立即返回错误，跳过后面的执行链
func (c *Context) AbortFai(code int, msg string, data any) {
	jsonData := cst.KV{
		"status": statusFai,
		"code":   code,
		"msg":    msg,
	}
	if data != nil {
		jsonData[dataField] = data
	}
	bytes, _ := jsonx.Marshal(jsonData)
	c.AbortDirect(http.StatusOK, bytes)
}

func (c *Context) AbortRet(ret *cst.Ret) {
	c.AbortFai(ret.Code, ret.Msg, ret.Data)
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
	http.ServeFile(c.Res, c.Req.Raw, filepath)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Render writes the response headers and calls render.Render to render data.
// 返回数据的接口
// 可以自定义扩展自己需要的Render
func (c *Context) Render(resStatus int, r render.Render) {
	if c.tryToRender() == false {
		return
	}

	c.Res.WriteHeader(resStatus)
	// 如果指定的返回状态，不能返回数据内容。需要特殊处理
	if !bodyAllowedForStatus(resStatus) {
		r.WriteContentType(c.Res)
		return
	}

	// 返回结果先写入缓存
	if err := r.Write(c.Res); err != nil {
		panic(err)
	}

	if c.route.ptrNode.hasBeforeSend {
		c.execBeforeSendHandlers() // 可以抛出异常，终止 Send data
	}

	// 一般来说，有session的请求，需要更新时间搓，保存session状态
	if c.Sess != nil {
		c.Sess.Save()
	}

	// ** now really send response data **
	if _, err := c.Res.Send(); err != nil {
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
	_ = c.Res.SendHijack(resStatus, lang.ToBytes(stream))
}

func (c *Context) AbortRedirect(resStatus int, redirectUrl string) {
	if c.tryToRender() == false {
		return
	}
	c.execIdx = maxRouteHandlers
	c.Res.SendHijackRedirect(c.Req.Raw, resStatus, redirectUrl)
}

// 这个是为超时返回准备的特殊方法，一般不要使用
func (c *Context) RenderTimeout(resStatus int, hint any) bool {
	return c.Res.SendByTimeoutGoroutine(resStatus, lang.ToBytes(hint))
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
	c.lock.Lock()
	if c.rendered {
		c.lock.Unlock()
		logx.Warn("Double render, the call canceled.")
		return false
	}
	c.rendered = true
	c.lock.Unlock()
	return true
}
