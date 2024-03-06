// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"net"
	"strings"
)

func (c *Context) GetHeader(key string) string {
	return c.Req.Raw.Header.Get(key)
}

func (c *Context) SetHeader(key, value string) {
	c.Res.Header().Set(key, value)
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
func (c *Context) ClientIP() string {
	if c.app.WebConfig.ForwardedByClientIP {
		// clientIP := c.GetHeader("X-Forwarded-For")
		// clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		ip := strings.TrimSpace(c.GetHeader("X-Forwarded-For"))
		if ip == "" {
			ip = strings.TrimSpace(c.GetHeader("X-Real-Ip"))
		}
		if ip != "" {
			return ip
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Req.Raw.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	ctType := c.Req.Raw.Header.Get("Content-Type")
	for i := range ctType {
		if ctType[i] == ' ' || ctType[i] == ';' {
			return ctType[:i]
		}
	}
	return ctType
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.GetHeader("Connection")), "upgrade") &&
		strings.EqualFold(c.GetHeader("Upgrade"), "websocket") {
		return true
	}
	return false
}
