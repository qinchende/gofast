// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fuse

import (
	"github.com/qinchende/gofast/aid/mathx"
	"strings"
	"time"
)

const (
	numHistoryReasons = 3
	timeFormatReason  = "15:04:05#"
)

// 错误信息滑动窗口
type errorWindow struct {
	reasonsTime [numHistoryReasons]time.Time
	reasons     [numHistoryReasons]string
	index       int
	count       int
	// add by cd.net on 20221116 没必要加锁，量大的时候影响性能
	//lock        sync.Mutex
}

func (ew *errorWindow) add(reason string) {
	//ew.lock.Lock()
	ew.reasonsTime[ew.index] = time.Now()
	ew.reasons[ew.index] = reason
	ew.index = (ew.index + 1) % numHistoryReasons
	ew.count = mathx.MinInt(ew.count+1, numHistoryReasons)
	//ew.lock.Unlock()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var (
	tmpReasonMem = [numHistoryReasons]string{}
	tmpReasons   []string
)

// 注意：整体上这个是非线程安全的
func (ew *errorWindow) Errors() []string {
	count := 0
	//ew.lock.Lock()
	for i := ew.index - 1; i >= ew.index-ew.count; i-- {
		idx := (i + numHistoryReasons) % numHistoryReasons
		tmpReasonMem[count] = ew.reasonsTime[idx].Format(timeFormatReason) + ew.reasons[idx]
		count++
	}
	//ew.lock.Unlock()

	tmpReasons = tmpReasonMem[0:count]
	defer func() {
		tmpReasons = tmpReasons[0:0]
	}()
	return tmpReasons
}

func (ew *errorWindow) ErrorsJoin(sep string) string {
	return strings.Join(ew.Errors(), sep)
}
