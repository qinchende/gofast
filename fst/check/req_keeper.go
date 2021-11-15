// Copyright 2021 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package check

import (
	"github.com/qinchende/gofast/skill/executors"
	"os"
)

type RequestKeeper struct {
	executor  *executors.IntervalExecutor
	container *reqContainer
}

func CreateReqKeeper(fp FuncGetPath) *RequestKeeper {
	container := &reqContainer{
		getPath: fp,
		//name: name,
		pid: os.Getpid(),
	}

	return &RequestKeeper{
		executor:  executors.NewIntervalExecutor(LogInterval, container),
		container: container,
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
