package fst

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/dts"
)

var rHandlers []*RHandler // 所有配置项汇总
type (
	newSuperKV func() cst.SuperKV
	RHandler   struct {
		rIndex    uint16     // 索引位置
		handler   CtxHandler // 处理函数
		pmsNew    newSuperKV // 解析到具体的struct对象
		pmsFields []string   // 从结构体类型解析出的字段，需要排序，相当于解析到 map
	}
)

// 添加一个路由属性对象
func (rh *RHandler) BindRoute(ri *RouteItem) {
	// 该路由还没有绑定任何处理函数
	if ri.routeIdx <= 0 && rh.handler != nil {
		ri.Handle(rh.handler)
	}
	// 如果不是有效的RouteItem
	if ri.routeIdx <= 0 {
		return
	}
	rh.rIndex = ri.routeIdx
	rHandlers = append(rHandlers, rh)
}

// 克隆对象
func (rh *RHandler) Clone() RouteAttrs {
	clone := *rh
	return &clone
}

// 构建所有路由的属性数组。没有指定的就用默认值填充。
func RebuildRHandlers(routesLen uint16) {
	raw := rHandlers
	rHandlers = make([]*RHandler, routesLen)

	for i := range raw {
		lang.SortByLen(raw[i].pmsFields)
		rHandlers[raw[i].rIndex] = raw[i]
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func WrapHandler(hd CtxHandler, fn newSuperKV, cls []string) *RHandler {
	return &RHandler{
		handler:   hd,
		pmsNew:    fn,
		pmsFields: cls,
	}
}

func ToSuperKV(v any) cst.SuperKV {
	return dts.AsSuperKV(v)
}

func NewSuperKV[T any]() cst.SuperKV {
	return dts.AsSuperKV(new(T))
}

func PmsAs[T any](c *Context) *T {
	return (*T)((c.Pms).(*dts.StructKV).Ptr)
}
