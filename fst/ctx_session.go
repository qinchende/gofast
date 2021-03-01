// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

type SessionKeeper interface {
	Get(string)
	Set(string, interface{})
	Save()
	Delete(string)
}

type GFSession struct {
	Sid    string
	Token  string
	Values map[interface{}]interface{}
	IsNew  bool
	Saved  bool
}

// GFSession 需要实现 SessionKeeper 所有接口
var _ SessionKeeper = &GFSession{}

func (ss *GFSession) Get(key string) {

}

func (ss *GFSession) Set(key string, val interface{}) {
	ss.Saved = false
}

func (ss *GFSession) Save() {
	ss.Saved = true
}

func (ss *GFSession) Delete(key string) {
	ss.Saved = false
}
