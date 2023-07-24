package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
)

var routesAttrs []*RHandler // 所有配置项汇总
type (
	RHandler struct {
		RIndex    uint16             // 索引位置
		PmsFields []string           // 从结构体类型解析出的字段，需要排序，相当于解析到 map
		PmsNew    func() cst.SuperKV // 解析到具体的struct对象
		Handler   CtxHandler         // 处理函数
	}
)

// 添加一个路由属性对象
func (ras *RHandler) BindRoute(ri *RouteItem) {
	if ri.Index() <= 0 && ras.Handler != nil {
		ri.Handle(ras.Handler)
	}
	// 如果不是有效的RouteItem
	if ri.Index() <= 0 {
		return
	}
	ras.RIndex = ri.Index()
	routesAttrs = append(routesAttrs, ras)
}

// 克隆对象
func (ras *RHandler) Clone() RouteAttrs {
	clone := *ras
	return &clone
}

// 构建所有路由的属性数组。没有指定的就用默认值填充。
func RebuildRHandlers(routesLen uint16) {
	old := routesAttrs
	routesAttrs = make([]*RHandler, routesLen)
	for i := range old {
		lang.SortByLen(old[i].PmsFields)
		routesAttrs[old[i].RIndex] = old[i]
	}
}
