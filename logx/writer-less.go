package logx

import (
	"github.com/qinchende/gofast/skill/syncx"
	"github.com/qinchende/gofast/skill/timex"
	"io"
	"sync/atomic"
	"time"
)

type lessWriter struct {
	*limitedExecutor
	writer io.Writer
}

//
//func NewLessWriter(writer io.Writer, milliseconds int) *lessWriter {
//	return &lessWriter{
//		limitedExecutor: newLimitedExecutor(milliseconds),
//		writer:          writer,
//	}
//}

func (w *lessWriter) Write(p []byte) (n int, err error) {
	w.logOrDiscard(func() {
		w.writer.Write(p)
	})
	return len(p), nil
}

type limitedExecutor struct {
	threshold time.Duration
	lastTime  *syncx.AtomicDuration
	discarded uint32
}

func newLimitedExecutor(milliseconds int) *limitedExecutor {
	return &limitedExecutor{
		threshold: time.Duration(milliseconds) * time.Millisecond,
		lastTime:  syncx.NewAtomicDuration(),
	}
}

func (le *limitedExecutor) logOrDiscard(execute func()) {
	if le == nil || le.threshold <= 0 {
		execute()
		return
	}

	now := timex.Now()
	if now-le.lastTime.Load() <= le.threshold {
		atomic.AddUint32(&le.discarded, 1)
	} else {
		le.lastTime.Set(now)
		discarded := atomic.SwapUint32(&le.discarded, 0)
		if discarded > 0 {
			ErrorF("Discarded %d error messages", discarded)
		}
		execute()
	}
}
