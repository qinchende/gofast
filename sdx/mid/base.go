// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import "github.com/qinchende/gofast/fst"

type (
	RAttrs struct {
		RIndex  uint16 // 此路由在路由数组中的索引值
		MaxLen  int64  `v:"def=0"` // 最大请求字节数，32MB（def=33554432）
		Timeout int32  `v:"def=0"` // 超时时间毫秒，默认去全局的超时时间

		//MaxReq    int32   `cnf:",def=1000000,range=[0:100000000]"` // 支持最大并发量 (对单个请求不支持这个参数，这个是由自适应降载逻辑自动判断的)
		//BreakRate float32 `cnf:",def=3000,range=[0:600000]"` // google sre算法K值敏感度，K 越小越容易丢请求，推荐 1.5-2 之间 （这个算法目前底层写死1.5，基本上通用了，不必每个路由单独设置）
	}
	attrsList []*RAttrs // 高级功能：每项路由可选配置，精准控制
)

var (
	RAttrsList attrsList // 所有配置项汇总
)

func (ras *RAttrs) SetRouteIndex(routeIdx uint16) {
	ras.RIndex = routeIdx
	RAttrsList = append(RAttrsList, ras)
}

// 对当前配置项，按照route索引顺序排序
func (*attrsList) Reordering(app *fst.GoFast, rtLen uint16) {
	old := RAttrsList
	RAttrsList = make(attrsList, rtLen, rtLen)
	for i := 0; i < len(old); i++ {
		it := old[i]
		RAttrsList[it.RIndex] = it
	}

	// 设置默认值
	var defAttrs RAttrs
	defAttrs.MaxLen = 0
	defAttrs.Timeout = int32(app.SdxConfig.DefTimeoutMS)
	//defAttrs.MaxReq = 1000000
	//defAttrs.BreakRate = 1.5

	for i := 0; i < len(RAttrsList); i++ {
		if RAttrsList[i] == nil {
			RAttrsList[i] = &defAttrs
		}
	}
}
