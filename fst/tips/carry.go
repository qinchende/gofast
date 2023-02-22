// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package tips

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
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
func (b *CarryItem) SetType(flags CarryType) *CarryItem {
	b.Type = flags
	return b
}

func (b *CarryItem) SetMeta(data any) *CarryItem {
	b.Meta = data
	return b
}

func (b *CarryItem) IsType(flags CarryType) bool {
	return (b.Type & flags) > 0
}

func (b *CarryItem) JSON() any {
	hash := cst.KV{}
	if b.Meta != nil {
		value := reflect.ValueOf(b.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return b.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				hash[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			hash["meta"] = b.Meta
		}
	}
	return hash
}

func (b *CarryItem) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(b.JSON())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (bs CarryList) ByType(typ CarryType) CarryList {
	if len(bs) == 0 {
		return nil
	}
	if typ == CarryTypeAny {
		return bs
	}
	var bsTmp CarryList
	for _, b := range bs {
		if b.IsType(typ) {
			bsTmp = append(bsTmp, b)
		}
	}
	return bsTmp
}

func (bs CarryList) Last() *CarryItem {
	if length := len(bs); length > 0 {
		return bs[length-1]
	}
	return nil
}

// 收集所有Basket中的Msg
func (bs CarryList) CollectMessages() []string {
	if len(bs) == 0 {
		return nil
	}
	msgStrings := make([]string, len(bs), len(bs))
	for i, b := range bs {
		msgStrings[i] = b.Msg
	}
	return msgStrings
}

func (bs CarryList) JSON() any {
	switch len(bs) {
	case 0:
		return nil
	case 1:
		return bs.Last().JSON()
	default:
		json := make([]any, len(bs), len(bs))
		for i, b := range bs {
			json[i] = b.JSON()
		}
		return json
	}
}

func (bs CarryList) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(bs.JSON())
}
