// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package gate

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

func CreateReqKeeper(name string, fp FuncGetPath) *RequestKeeper {
	container := &reqContainer{
		name:    name,
		pid:     os.Getpid(),
		getPath: fp,
	}

	return &RequestKeeper{
		executor:  executors.NewIntervalExecutor(LogInterval, container),
		container: container,
	}
}

func (rk *RequestKeeper) SetBreakers(length uint16) {
	rk.Breakers = make([]breaker.Breaker, 0, length)
	for i := 0; i < int(length); i++ {
		rk.Breakers = append(rk.Breakers, breaker.NewBreaker(breaker.WithName(strconv.Itoa(i))))
	}
}

func (rk *RequestKeeper) AddItem(item ReqItem) {
	rk.executor.Add(item)
}

// AddDrop adds a drop to m.
func (rk *RequestKeeper) AddDrop() {
	rk.executor.Add(ReqItem{
		Drop: true,
	})
}
