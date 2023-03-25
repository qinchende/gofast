package fst

import "github.com/qinchende/gofast/skill/lang"

type (
	GAttrs struct {
		PmsFields []string
	}

	RAttrs struct {
		RIndex    uint16   // 索引位置
		PmsFields []string // 从结构体类型解析出的字段，需要排序
	}
	listAttrs []*RAttrs // 高级功能：每项路由可选配置，精准控制
)

var routesAttrs listAttrs // 所有配置项汇总

// 添加一个路由属性对象
func (ras *RAttrs) SetIndex(routeIdx uint16) {
	ras.RIndex = routeIdx
	routesAttrs = append(routesAttrs, ras)
}
func (ras *RAttrs) Clone() RouteAttrs {
	fls := make([]string, len(ras.PmsFields))
	copy(fls, ras.PmsFields)

	clone := &RAttrs{
		RIndex:    ras.RIndex,
		PmsFields: fls,
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

	// 没有的项看是否给默认值
	//for i := range routesAttrs {
	//	if it == nil {
	//		routesAttrs[idx] = &RAttrs{}
	//	}
	//}
}
