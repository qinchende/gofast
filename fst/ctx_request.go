// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// 请求从net/http包传出来之后，需要在框架中转换成我们自己的Request对象
type Request struct {
	RawReq *http.Request

	Errors errorMsgs
	gftApp *GoFast
	gftCtx *Context
	fitIdx int
}

func (r *Request) requestHeader(key string) string {
	return r.RawReq.Header.Get(key)
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func (r *Request) ClientIP() string {
	if r.gftApp.ForwardedByClientIP {
		clientIP := r.requestHeader("X-Forwarded-For")
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		if clientIP == "" {
			clientIP = strings.TrimSpace(r.requestHeader("X-Real-Ip"))
		}
		if clientIP != "" {
			return clientIP
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RawReq.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func (r *Request) Error(err error) *Error {
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

	r.Errors = append(r.Errors, parsedError)
	return parsedError
}

func (r *Request) Errorf(format string, v ...interface{}) {
	r.Error(errors.New(fmt.Sprintf(format, v...)))
}
