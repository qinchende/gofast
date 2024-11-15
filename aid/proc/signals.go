//go:build linux || darwin
// +build linux darwin

package proc

import (
	"github.com/qinchende/gofast/aid/logx"
	"os"
	"os/signal"
	"syscall"
)

const timeFormat = "0102150405"

func init() {
	go func() {
		var profiler Stopper

		// https://golang.org/pkg/os/signal/#Notify
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				dumpGoroutines()
			case syscall.SIGUSR2:
				if profiler == nil {
					profiler = StartProfile()
				} else {
					profiler.Stop()
					profiler = nil
				}
			case syscall.SIGTERM:
				gracefulStop(signals)
			default:
				logx.Err("Got unregistered signal:", v)
			}
		}
	}()
}
