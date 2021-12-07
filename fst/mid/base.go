// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

type funcServeHTTP func(w http.ResponseWriter, r *http.Request)

type FitHelper struct {
	nextHandler funcServeHTTP
}

func (fh *FitHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fh.nextHandler(w, r)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 高级功能：每项路由可选配置，精准控制
type RouteConfigs []*RouteConfig

type RouteConfig struct {
	fst.RouteIndex
	MaxReq        int32   `cnf:",def=1000000,range=[0:100000000]"` // 支持最大并发量
	MaxContentLen int64   `cnf:",def=0"`                           // 最大请求字节数，32MB（def=33554432）
	Timeout       int32   `cnf:",def=3000,range=[0:600000]"`       // 超时时间毫秒
	BreakRate     float32 `cnf:",def=3000,range=[0:600000]"`       // google sre算法K值敏感度，K 越小越容易丢请求，推荐 1.5-2 之间
}

func (rc *RouteConfig) AddToList(idx uint16) {
	rc.Idx = idx
	RConfigs = append(RConfigs, rc)

}

// 所有配置项汇总
var RConfigs RouteConfigs

// 对当前配置项，按照route索引顺序排序
func (rcs *RouteConfigs) Reordering(rtLen uint16) {
	old := RConfigs
	RConfigs = make(RouteConfigs, rtLen, rtLen)
	for i := 0; i < len(old); i++ {
		it := old[i]
		RConfigs[it.Idx] = it
	}

	defItem := RouteConfig{}
	for i := 0; i < len(RConfigs); i++ {
		if RConfigs[i] == nil {
			RConfigs[i] = &defItem
		}
	}
}
