// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
)

//type funcServeHTTP func(w http.ResponseWriter, r *http.Request)
//
//type FitHelper struct {
//	nextHandler funcServeHTTP
//}
//
//func (fh *FitHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	fh.nextHandler(w, r)
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 高级功能：每项路由可选配置，精准控制
type routeConfigs []*RConfig

// 所有配置项汇总
var RConfigs routeConfigs
var defConfig RConfig

type RConfig struct {
	fst.RouteIndex
	MaxLen  int64 `v:"def=0"` // 最大请求字节数，32MB（def=33554432）
	Timeout int32 `v:"def=0"` // 超时时间毫秒，默认去全局的超时时间

	//MaxReq    int32   `cnf:",def=1000000,range=[0:100000000]"` // 支持最大并发量 (对单个请求不支持这个参数，这个是由自适应降载逻辑自动判断的)
	//BreakRate float32 `cnf:",def=3000,range=[0:600000]"` // google sre算法K值敏感度，K 越小越容易丢请求，推荐 1.5-2 之间 （这个算法目前底层写死1.5，基本上通用了，不必每个路由单独设置）
}

func (rc *RConfig) AddToList(idx uint16) {
	rc.Idx = idx
	RConfigs = append(RConfigs, rc)
}

// 对当前配置项，按照route索引顺序排序
func (rcs *routeConfigs) Reordering(app *fst.GoFast, rtLen uint16) {
	old := RConfigs
	RConfigs = make(routeConfigs, rtLen, rtLen)
	for i := 0; i < len(old); i++ {
		it := old[i]
		RConfigs[it.Idx] = it
	}

	// 设置默认值
	defConfig.MaxLen = 0
	defConfig.Timeout = int32(app.SdxDefTimeout)
	//defConfig.MaxReq = 1000000
	//defConfig.BreakRate = 1.5

	for i := 0; i < len(RConfigs); i++ {
		if RConfigs[i] == nil {
			RConfigs[i] = &defConfig
		}
	}
}
