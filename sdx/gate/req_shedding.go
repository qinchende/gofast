// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

import (
	"sync/atomic"
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
	go st.logPrintCpuShedding()
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
