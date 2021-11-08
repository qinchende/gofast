package fst

//
//// 经过路由匹配算法之后，第一个到这里执行
//func theFirstBeforeHandler(c *Context) {
//	// 检查，熔断
//}
//
//// 全局：（已经路由匹配），请求最后需要经过的处理函数
//// 这里可以统计请求的访问频率以及系统的响应情况。
//func theLastAfterHandler(c *Context) {
//	//var nodeIdx int16 = -1
//	//if c.match.ptrNode != nil {
//	//	nodeIdx = c.match.ptrNode.routerIdx
//	//}
//	//door.Keeper.AddItem(door.ReqItem{
//	//	RouterIdx: nodeIdx,
//	//	Duration:  timex.Since(c.EnterTime),
//	//})
//}
