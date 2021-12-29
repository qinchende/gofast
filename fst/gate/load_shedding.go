// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"sync/atomic"
	"time"

	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/stat"
)

type (
	// A SheddingStat is used to store the statistics for load shedding.
	SheddingStat struct {
		name  string
		total int64
		pass  int64
		drop  int64
	}

	snapshot struct {
		Total int64
		Pass  int64
		Drop  int64
	}
)

// NewSheddingStat returns a SheddingStat.
func NewSheddingStat(name string) *SheddingStat {
	st := &SheddingStat{
		name: name,
	}
	go st.run()
	return st
}

// IncrementTotal increments the total requests.
func (s *SheddingStat) Total() {
	atomic.AddInt64(&s.total, 1)
}

// IncrementPass increments the passed requests.
func (s *SheddingStat) Pass() {
	atomic.AddInt64(&s.pass, 1)
}

// IncrementDrop increments the dropped requests.
func (s *SheddingStat) Drop() {
	atomic.AddInt64(&s.drop, 1)
}

// 重置计数器
func (s *SheddingStat) reset() snapshot {
	return snapshot{
		Total: atomic.SwapInt64(&s.total, 0),
		Pass:  atomic.SwapInt64(&s.pass, 0),
		Drop:  atomic.SwapInt64(&s.drop, 0),
	}
}

// 单独的协程运行这个定时任务。启动定时日志输出
func (s *SheddingStat) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// 定时器，每分钟执行一次，死循环
	for range ticker.C {
		st := s.reset()
		if st.Total == 0 {
			continue
		}
		cpu := stat.CpuUsage()

		if st.Drop == 0 {
			logx.Statf("(%s)[1m], cpu: %d, total: %d, pass: %d, drop: %d", s.name, cpu, st.Total, st.Pass, st.Drop)
		} else {
			logx.Errorf("(%s)[1m], cpu: %d, total: %d, pass: %d, drop: %d", s.name, cpu, st.Total, st.Pass, st.Drop)
		}
	}
}
