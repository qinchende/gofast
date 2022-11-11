package fuse

import (
	"github.com/qinchende/gofast/skill/mathx"
	"github.com/qinchende/gofast/skill/timex"
	"strings"
	"sync"
)

const (
	numHistoryReasons = 3
	timeFormatReason  = "15:04:05#"
)

// 错误信息滑动窗口
type errorWindow struct {
	reasons [numHistoryReasons]string
	index   int
	count   int
	lock    sync.Mutex
}

func (ew *errorWindow) add(reason string) {
	ew.lock.Lock()
	ew.reasons[ew.index] = timex.Time().Format(timeFormatReason) + reason
	ew.index = (ew.index + 1) % numHistoryReasons
	ew.count = mathx.MinInt(ew.count+1, numHistoryReasons)
	ew.lock.Unlock()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var (
	tmpReasonMem = [numHistoryReasons]string{}
	tmpReasons   []string
)

// 注意：整体上这个是非线程安全的
func (ew *errorWindow) Errors() []string {
	count := 0
	ew.lock.Lock()
	for i := ew.index - 1; i >= ew.index-ew.count; i-- {
		tmpReasonMem[count] = ew.reasons[(i+numHistoryReasons)%numHistoryReasons]
		count++
	}
	ew.lock.Unlock()

	tmpReasons = tmpReasonMem[0:count]
	defer func() {
		tmpReasons = tmpReasons[0:0]
	}()
	return tmpReasons
}

func (ew *errorWindow) ErrorsJoin(sep string) string {
	return strings.Join(ew.Errors(), sep)
}
