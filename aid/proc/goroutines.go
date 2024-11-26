//go:build linux || darwin
// +build linux darwin

package proc

import (
	"fmt"
	"github.com/qinchende/gofast/aid/logx"
	"os"
	"path"
	"runtime/pprof"
	"syscall"
	"time"
)

const (
	goroutineProfile = "goroutine"
	debugLevel       = 2
)

func dumpGoroutines() {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	dumpFile := path.Join(os.TempDir(), fmt.Sprintf("%s-%d-goroutines-%s.dump",
		command, pid, time.Now().Format(timeFormat)))

	logx.Info().SendMsgF("Got dump goroutine signal, printing goroutine profile to %s", dumpFile)

	if f, err := os.Create(dumpFile); err != nil {
		logx.Err().SendMsgF("Failed to dump goroutine profile, error: %v", err)
	} else {
		defer f.Close()
		pprof.Lookup(goroutineProfile).WriteTo(f, debugLevel)
	}
}
