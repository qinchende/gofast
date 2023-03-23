// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/cst"
)

type (
	Attrs struct {
		RIndex    uint16 `v:""`                       // 索引位置
		Priority  int16  `v:"def=500,range=[0:1000]"` // 业务优先级
		MaxLen    int64  `v:""`                       // 最大请求长度，0不限制
		TimeoutMS int32  `v:""`                       // 超时时间毫秒

		//MaxReq    int32   `cnf:",def=1000000,range=[0:100000000]"` // 支持最大并发量 (对单个请求不支持这个参数，这个是由自适应降载逻辑自动判断的)
		//BreakRate float32 `cnf:",def=3000,range=[0:600000]"` // google sre算法K值敏感度，K 越小越容易丢请求，推荐 1.5-2 之间 （这个算法目前底层写死1.5，基本上通用了，不必每个路由单独设置）
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
func (*listAttrs) Rebuild(routesLen uint16, cnf *cst.SdxConfig) {
	old := AttrsList
	AttrsList = make(listAttrs, routesLen)
	for _, it := range old {
		AttrsList[it.RIndex] = it
	}

	defAttrs := Attrs{
		MaxLen:    0,
		TimeoutMS: int32(cnf.DefTimeoutMS),
		//MaxReq:    1000000,
		//BreakRate: 1.5,
	}
	for idx, it := range AttrsList {
		if it == nil {
			AttrsList[idx] = &defAttrs
		}
	}
}
