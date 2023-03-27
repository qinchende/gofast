package fst

import "github.com/qinchende/gofast/skill/lang"

type (
	RAttrs struct {
		RIndex    uint16   // 索引位置
		PmsFields []string // 从结构体类型解析出的字段，需要排序
		Handler   CtxHandler
	}
	listAttrs []*RAttrs // 高级功能：每项路由可选配置，精准控制
)

var routesAttrs listAttrs // 所有配置项汇总

// 添加一个路由属性对象
func (ras *RAttrs) BindRoute(ri *RouteItem) {
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
func (ras *RAttrs) Clone() RouteAttrs {
	fls := make([]string, len(ras.PmsFields))
	copy(fls, ras.PmsFields)

	clone := &RAttrs{
		RIndex:    ras.RIndex,
		PmsFields: fls,
		Handler:   ras.Handler,
	}
	return clone
}

// 构建所有路由的属性数组。没有指定的就用默认值填充。
func (*listAttrs) Rebuild(routesLen uint16) {
	old := routesAttrs
	routesAttrs = make(listAttrs, routesLen)
	for i := range old {
		lang.SortByLen(old[i].PmsFields)
		routesAttrs[old[i].RIndex] = old[i]
	}
}
