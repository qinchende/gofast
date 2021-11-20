// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package check

import (
	"github.com/qinchende/gofast/skill/breaker"
	"github.com/qinchende/gofast/skill/executors"
	"os"
	"strconv"
)

type RequestKeeper struct {
	executor  *executors.IntervalExecutor
	container *reqContainer
	Breakers  []breaker.Breaker
}

func CreateReqKeeper(name string, len uint16, fp FuncGetPath) *RequestKeeper {
	container := &reqContainer{
		name:    name,
		pid:     os.Getpid(),
		getPath: fp,
	}

	// 建好 breakers
	bks := make([]breaker.Breaker, len)
	for i := 0; i < int(len); i++ {
		bks = append(bks, breaker.NewBreaker(breaker.WithName(strconv.Itoa(i))))
	}

	return &RequestKeeper{
		executor:  executors.NewIntervalExecutor(LogInterval, container),
		container: container,
		Breakers:  bks,
	}
}

func (ct *RequestKeeper) AddItem(item ReqItem) {
	ct.executor.Add(item)
}

// AddDrop adds a drop to m.
func (ct *RequestKeeper) AddDrop() {
	ct.executor.Add(ReqItem{
		Drop: true,
	})
}
