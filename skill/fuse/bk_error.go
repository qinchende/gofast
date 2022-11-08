package fuse

import (
	"fmt"
	"github.com/qinchende/gofast/skill/mathx"
	"github.com/qinchende/gofast/skill/timex"
	"strings"
	"sync"
)

const (
	numHistoryReasons = 5
	timeFormat        = "15:04:05"
)

type errorWindow struct {
	reasons [numHistoryReasons]string
	index   int
	count   int
	lock    sync.Mutex
}

func (ew *errorWindow) add(reason string) {
	ew.lock.Lock()
	ew.reasons[ew.index] = fmt.Sprintf("%s %s", timex.Time().Format(timeFormat), reason)
	ew.index = (ew.index + 1) % numHistoryReasons
	ew.count = mathx.MinInt(ew.count+1, numHistoryReasons)
	ew.lock.Unlock()
}

func (ew *errorWindow) String() string {
	var reasons []string

	ew.lock.Lock()
	// reverse order
	for i := ew.index - 1; i >= ew.index-ew.count; i-- {
		reasons = append(reasons, ew.reasons[(i+numHistoryReasons)%numHistoryReasons])
	}
	ew.lock.Unlock()

	return strings.Join(reasons, "\n")
}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//// 包了一个自定义的 reason 对象
//type myPromise struct {
//	promise promiseNoReason
//	errWin  *errorWindow
//}
//
//func (p myPromise) Accept() {
//	p.promise.Accept()
//}
//
//func (p myPromise) Reject(reason string) {
//	p.errWin.add(reason)
//	p.promise.Reject()
//}
