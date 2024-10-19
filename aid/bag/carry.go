// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package bag

import (
	"github.com/qinchende/gofast/aid/jsonx"
	"github.com/qinchende/gofast/core/cst"
	"reflect"
)

type (
	CarryType uint
	CarryItem struct {
		Type CarryType // 数据分类
		Msg  string    // 描述信息
		Meta any       // 详细数据
	}
	CarryList []*CarryItem
)

const (
	CarryTypeAny CarryType = 0
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (it *CarryItem) SetType(flags CarryType) *CarryItem {
	it.Type = flags
	return it
}

func (it *CarryItem) SetMeta(data any) *CarryItem {
	it.Meta = data
	return it
}

func (it *CarryItem) IsType(flags CarryType) bool {
	return (it.Type & flags) > 0
}

func (it *CarryItem) JSON() any {
	hash := cst.KV{}
	if it.Meta != nil {
		value := reflect.ValueOf(it.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return it.Meta
		case reflect.Map:
			keys := value.MapKeys()
			for i := range keys {
				hash[keys[i].String()] = value.MapIndex(keys[i]).Interface()
			}
		default:
			hash["meta"] = it.Meta
		}
	}
	return hash
}

func (it *CarryItem) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(it.JSON())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (list CarryList) ByType(typ CarryType) CarryList {
	if len(list) == 0 {
		return nil
	}
	if typ == CarryTypeAny {
		return list
	}
	var bsTmp CarryList
	for i := range list {
		if list[i].IsType(typ) {
			bsTmp = append(bsTmp, list[i])
		}
	}
	return bsTmp
}

// 知道第一个符合类型的项
func (list CarryList) FirstOne(typ CarryType) *CarryItem {
	for i := range list {
		if list[i].IsType(typ) {
			return list[i]
		}
	}
	return nil
}

func (list CarryList) Last() *CarryItem {
	if length := len(list); length > 0 {
		return list[length-1]
	}
	return nil
}

// 收集 items 中的 Msg 字段
func (list CarryList) CollectMessages() []string {
	if len(list) == 0 {
		return nil
	}
	msgStrings := make([]string, len(list), len(list))
	for i := range list {
		msgStrings[i] = list[i].Msg
	}
	return msgStrings
}

func (list CarryList) JSON() any {
	switch len(list) {
	case 0:
		return nil
	case 1:
		return list.Last().JSON()
	default:
		json := make([]any, len(list), len(list))
		for i := range list {
			json[i] = list[i].JSON()
		}
		return json
	}
}

func (list CarryList) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(list.JSON())
}
