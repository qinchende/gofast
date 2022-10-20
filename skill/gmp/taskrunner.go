package gmp

import (
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/lang"
)

type TaskRunner struct {
	limitChan chan lang.PlaceholderType
}

func NewTaskRunner(concurrency int) *TaskRunner {
	return &TaskRunner{
		limitChan: make(chan lang.PlaceholderType, concurrency),
	}
}

func (rp *TaskRunner) Schedule(task func()) {
	rp.limitChan <- lang.Placeholder

	go func() {
		//defer rescue.Recover(func() {
		//	<-rp.limitChan
		//})

		defer func() {
			<-rp.limitChan

			if p := recover(); p != nil {
				logx.Stacks(p)
			}
		}()

		task()
	}()
}
