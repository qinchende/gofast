// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import "errors"

type SessionKeeper interface {
	New() *CtxSession
	Get(string) interface{}
	Set(string, interface{})
	Save()
	Delete(string)
}

// GoFast框架的 Context Session
// 默认将使用 Redis 存放 分布式 session 信息
type CtxSession struct {
	Sid        string
	Token      string
	TokenIsNew bool
	Saved      bool
	Values     map[string]interface{}
}

// CtxSession 需要实现 SessionKeeper 所有接口
var _ SessionKeeper = &CtxSession{}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *CtxSession) New() *CtxSession {
	return &CtxSession{}
}

func (ss *CtxSession) Get(key string) interface{} {
	if ss.Values == nil {
		return nil
	}
	return ss.Values[key]
}

func (ss *CtxSession) Set(key string, val interface{}) {
	ss.Saved = false
	ss.Values[key] = val
}

// TODO: 你可以自定义实现这个save方法
var CtxSessionSaveFun = func(ss *CtxSession) (string, error) {
	return "", errors.New("Error. ")
}

func (ss *CtxSession) Save() {
	// 如果已经保存了，不会重复保存
	if ss.Saved == true {
		return
	}
	// 调用自定义函数保存当前 session
	_, err := CtxSessionSaveFun(ss)

	// TODO: 如果保存失败怎么办？目前是抛异常，本次请求直接返回错误。
	if err != nil {
		RaisePanic("Save session error.")
	} else {
		ss.Saved = true
	}
}

func (ss *CtxSession) Delete(key string) {
	delete(ss.Values, key)
	ss.Saved = false
}
