package wit

import (
	"github.com/qinchende/gofast/aid/jsonx"
	"github.com/qinchende/gofast/core/cst"
	"reflect"
)

type (
	KVItemGroup struct {
		KVItem      // KV键值对
		Group  uint // 分组标记
	}
	KVListGroup []KVItemGroup
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (it *KVItemGroup) SetGroup(flags uint) *KVItemGroup {
	it.Group = flags
	return it
}

func (it *KVItemGroup) IsGroup(flags uint) bool {
	return (it.Group & flags) > 0
}

func (it *KVItemGroup) JSON() any {
	hash := cst.KV{}
	if it.Val != nil {
		value := reflect.ValueOf(it.Val)
		switch value.Kind() {
		case reflect.Struct:
			return it.Val
		case reflect.Map:
			keys := value.MapKeys()
			for i := range keys {
				hash[keys[i].String()] = value.MapIndex(keys[i]).Interface()
			}
		default:
			hash["meta"] = it.Val
		}
	}
	return hash
}

func (it *KVItemGroup) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(it.JSON())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (list KVListGroup) ByGroup(typ uint) KVListGroup {
	if len(list) == 0 {
		return nil
	}
	if typ == 0 {
		return list // 默认初始值
	}
	var bsTmp KVListGroup
	for i := range list {
		if list[i].IsGroup(typ) {
			bsTmp = append(bsTmp, list[i])
		}
	}
	return bsTmp
}

// 知道第一个符合类型的项
func (list KVListGroup) FirstOne(typ uint) *KVItemGroup {
	for i := range list {
		if list[i].IsGroup(typ) {
			return &list[i]
		}
	}
	return nil
}

func (list KVListGroup) Last() *KVItemGroup {
	if length := len(list); length > 0 {
		return &list[length-1]
	}
	return nil
}

// 收集 items 中的 Msg 字段
func (list KVListGroup) CollectMessages() []string {
	if len(list) == 0 {
		return nil
	}
	msgStrings := make([]string, len(list), len(list))
	for i := range list {
		msgStrings[i] = *list[i].Val.(*string)
	}
	return msgStrings
}

func (list KVListGroup) JSON() any {
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

func (list KVListGroup) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(list.JSON())
}
