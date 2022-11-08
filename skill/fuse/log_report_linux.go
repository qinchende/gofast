// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
//go:build linux

package fuse

import (
	"flag"
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/sysx/host"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qinchende/gofast/skill/proc"
	"github.com/qinchende/gofast/skill/timex"
)

const (
	clusterNameKey  = "CLUSTER_NAME"
	testEnv         = "test.v"
	printTimeFormat = "2006-01-02 15:04:05"
)

var (
	reporter     = logx.Info
	lock         sync.RWMutex
	lessExecutor = exec.NewLess(time.Minute * 5)
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
			var sb strings.Builder
			fmt.Fprintf(&sb, "%s\n", timex.Time().Format(printTimeFormat))
			if len(clusterName) > 0 {
				fmt.Fprintf(&sb, "cluster: %s\n", clusterName)
			}
			fmt.Fprintf(&sb, "host: %s\n", host.Hostname())
			dp := atomic.SwapInt32(&dropped, 0)
			if dp > 0 {
				fmt.Fprintf(&sb, "dropped: %d\n", dp)
			}
			sb.WriteString(strings.TrimSpace(msg))
			fn(sb.String())
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
