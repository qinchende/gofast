//go:build linux
// +build linux

package logx

import (
	"flag"
	"fmt"
	"github.com/qinchende/gofast/skill/sysx/host"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qinchende/gofast/skill/executors"
	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/timex"
)

const (
	clusterNameKey = "CLUSTER_NAME"
	testEnv        = "test.v"
	timeFormat     = "2006-01-02 15:04:05"
)

var (
	reporter     = Info
	lock         sync.RWMutex
	lessExecutor = executors.NewLessExecutor(time.Minute * 5)
	dropped      int32
	clusterName  = proc.Env(clusterNameKey)
)

func init() {
	if flag.Lookup(testEnv) != nil {
		SetReporter(nil)
	}
}

func Report(msg string) {
	lock.RLock()
	fn := reporter
	lock.RUnlock()

	if fn != nil {
		reported := lessExecutor.DoOrDiscard(func() {
			var builder strings.Builder
			fmt.Fprintf(&builder, "%s\n", timex.Time().Format(timeFormat))
			if len(clusterName) > 0 {
				fmt.Fprintf(&builder, "cluster: %s\n", clusterName)
			}
			fmt.Fprintf(&builder, "host: %s\n", host.Hostname())
			dp := atomic.SwapInt32(&dropped, 0)
			if dp > 0 {
				fmt.Fprintf(&builder, "dropped: %d\n", dp)
			}
			builder.WriteString(strings.TrimSpace(msg))
			fn(builder.String())
		})
		if !reported {
			atomic.AddInt32(&dropped, 1)
		}
	}
}

func SetReporter(fn func(string)) {
	lock.Lock()
	defer lock.Unlock()
	reporter = fn
}
