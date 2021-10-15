package door

import "time"

type ReqTime struct {
	Duration  time.Duration // 任务耗时
	RouterIdx int16         // 路由节点的index
	Drop      bool          // 是否是一个丢弃的任务
}

func AddItem(item ReqTime) {

}
