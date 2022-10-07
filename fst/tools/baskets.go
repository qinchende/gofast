// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package tools

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"reflect"
)

// BasketType is an unsigned int error code as defined in the GoFast spec.
type (
	BasketType uint //type BasketType uint
	Basket     struct {
		Type BasketType // 数据分类
		Msg  string     // 描述信息
		Meta any        // 详细数据
	}
	Baskets []*Basket
)

const (
	BasketTypeAny BasketType = 0
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (b *Basket) SetType(flags BasketType) *Basket {
	b.Type = flags
	return b
}

func (b *Basket) SetMeta(data any) *Basket {
	b.Meta = data
	return b
}

func (b *Basket) IsType(flags BasketType) bool {
	return (b.Type & flags) > 0
}

func (b *Basket) JSON() any {
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

func (b *Basket) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(b.JSON())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (bs Baskets) ByType(typ BasketType) Baskets {
	if len(bs) == 0 {
		return nil
	}
	if typ == BasketTypeAny {
		return bs
	}
	var bsTmp Baskets
	for _, b := range bs {
		if b.IsType(typ) {
			bsTmp = append(bsTmp, b)
		}
	}
	return bsTmp
}

func (bs Baskets) Last() *Basket {
	if length := len(bs); length > 0 {
		return bs[length-1]
	}
	return nil
}

// 收集所有Basket中的Msg
func (bs Baskets) CollectMessages() []string {
	if len(bs) == 0 {
		return nil
	}
	msgStrings := make([]string, len(bs), len(bs))
	for i, b := range bs {
		msgStrings[i] = b.Msg
	}
	return msgStrings
}

func (bs Baskets) JSON() any {
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

func (bs Baskets) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(bs.JSON())
}
