// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"github.com/qinchende/gofast/skill/sysx"
	"sync/atomic"
	"time"

	"github.com/qinchende/gofast/logx"
)

type (
	// A SheddingStat is used to store the statistics for load shedding.
	sheddingStat struct {
		total int64
		pass  int64
		drop  int64
	}
	//
	//snapshot struct {
	//	total int64
	//	pass  int64
	//	drop  int64
	//}
)

// NewSheddingStat returns a SheddingStat.
func createSheddingStat() *sheddingStat {
	st := &sheddingStat{}
	go st.run()
	return st
}

// IncrementTotal increments the total requests.
func (s *sheddingStat) Total() {
	atomic.AddInt64(&s.total, 1)
}

// IncrementPass increments the passed requests.
func (s *sheddingStat) Pass() {
	atomic.AddInt64(&s.pass, 1)
}

// IncrementDrop increments the dropped requests.
func (s *sheddingStat) Drop() {
	atomic.AddInt64(&s.drop, 1)
}

// 重置计数器
func (s *sheddingStat) reset() sheddingStat {
	return sheddingStat{
		total: atomic.SwapInt64(&s.total, 0),
		pass:  atomic.SwapInt64(&s.pass, 0),
		drop:  atomic.SwapInt64(&s.drop, 0),
	}
}

// 单独的协程运行这个定时任务。启动定时日志输出
func (s *sheddingStat) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// 定时器，每分钟执行一次，死循环
	for range ticker.C {
		st := s.reset()
		if st.total == 0 && st.pass == 0 && st.drop == 0 {
			continue
		}
		cpu := sysx.CpuSmoothUsage()
		logx.Statf("down-1m | cpu: %d, total: %d, pass: %d, drop: %d", cpu, st.total, st.pass, st.drop)

		//if st.drop == 0 {
		//	logx.Statf("down-1m | cpu: %d, total: %d, pass: %d, drop: %d", cpu, st.total, st.pass, st.drop)
		//} else {
		//	logx.Errorf("down-1m | cpu: %d, total: %d, pass: %d, drop: %d", cpu, st.total, st.pass, st.drop)
		//}
	}
}
