package fst

import (
	"reflect"
)

type (
	Attrs struct {
		RIndex    uint16       `v:""` // 索引位置
		PmsStruct reflect.Type `v:""` // 收集参数的结构体类型
		PmsFields []string     `v:""` // 从结构体类型解析出的字段，需要排序
	}
	listAttrs []*Attrs // 高级功能：每项路由可选配置，精准控制
)

var AttrsList listAttrs // 所有配置项汇总

// 添加一个路由属性对象
func (ras *Attrs) SetIndex(routeIdx uint16) {
	ras.RIndex = routeIdx
	AttrsList = append(AttrsList, ras)
}

// 构建所有路由的属性数组。没有指定的就用默认值填充。
func (*listAttrs) rebuild(routesLen uint16) {
	old := AttrsList
	AttrsList = make(listAttrs, routesLen)
	for _, it := range old {
		AttrsList[it.RIndex] = it
	}

	defAttrs := Attrs{}
	for idx, it := range AttrsList {
		if it == nil {
			AttrsList[idx] = &defAttrs
		}
	}
}
