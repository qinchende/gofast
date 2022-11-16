package fuse

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/exec"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/proc"
	"time"
)

type autoBreaker struct {
	throttle
	name      string
	errWin    *errorWindow
	reduceLog *exec.Reduce
	showLog   bool
}

func NewGBreaker(name string, showLog bool) Breaker {
	b := autoBreaker{
		name:    name,
		errWin:  new(errorWindow),
		showLog: showLog,
	}
	if showLog {
		b.reduceLog = exec.NewReduce(time.Second * 30)
	}
	if len(b.name) == 0 {
		b.name = lang.Rand()
	}
	b.throttle = newGoogleThrottle()
	return &b
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ab *autoBreaker) Name() string {
	return ab.name
}

func (ab *autoBreaker) Accept() {
	ab.throttle.markSuc()
}

func (ab *autoBreaker) Reject(reason string) {
	ab.errWin.add(reason)
	ab.throttle.markFai()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ab *autoBreaker) Allow() error {
	return ab.logError(ab.throttle.allow())
}

func (ab *autoBreaker) Do(req funcReq) error {
	return ab.logError(ab.throttle.doReq(req, nil, defAcceptable))
}

func (ab *autoBreaker) DoWithAcceptable(req funcReq, cpt funcAcceptable) error {
	return ab.logError(ab.throttle.doReq(req, nil, cpt))
}

func (ab *autoBreaker) DoWithFallback(req funcReq, fb funcFallback) error {
	return ab.logError(ab.throttle.doReq(req, fb, defAcceptable))
}

func (ab *autoBreaker) DoWithFallbackAcceptable(req funcReq, fb funcFallback, cpt funcAcceptable) error {
	return ab.logError(ab.throttle.doReq(req, fb, cpt))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// utils
func defAcceptable(err error) bool {
	return err == nil
}

func (ab *autoBreaker) logError(err error) error {
	if ab.showLog && err != nil {
		ab.reduceLog.DoOrNot(func(skip int32) {
			if err != ErrServiceUnavailable {
				return
			}
			logx.InfoReport(cst.KV{
				"typ":    logx.LogStatBreakerOpen.Type,
				"proc":   proc.ProcessName() + "/" + lang.ToString(proc.Pid()),
				"callee": ab.name,
				"skip":   skip,
				"msg":    ab.errWin.Errors(),
			})
		})
	}
	return err
}
