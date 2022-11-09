package fuse

import (
	"fmt"
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
func (cb *autoBreaker) Name() string {
	return cb.name
}

func (cb *autoBreaker) Allow() error {
	return cb.logError(cb.throttle.allow())
}

func (cb *autoBreaker) Accept() {
	cb.throttle.markSuc()
}

func (cb *autoBreaker) Reject(reason string) {
	if reason != "" {
		cb.errWin.add(reason)
	}
	cb.throttle.markFai()
}

func (cb *autoBreaker) Do(req funcReq) error {
	return cb.logError(cb.throttle.doReq(req, nil, defAcceptable))
}

func (cb *autoBreaker) DoWithAcceptable(req funcReq, cpt funcAcceptable) error {
	return cb.logError(cb.throttle.doReq(req, nil, cpt))
}

func (cb *autoBreaker) DoWithFallback(req funcReq, fb funcFallback) error {
	return cb.logError(cb.throttle.doReq(req, fb, defAcceptable))
}

func (cb *autoBreaker) DoWithFallbackAcceptable(req funcReq, fb funcFallback, cpt funcAcceptable) error {
	return cb.logError(cb.throttle.doReq(req, fb, cpt))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// utils
func defAcceptable(err error) bool {
	return err == nil
}

func (cb *autoBreaker) logError(err error) error {
	if cb.showLog && err != nil {
		cb.reduceLog.DoOrNot(func(skip int32) {
			if err != ErrServiceUnavailable {
				return
			}
			logx.InfoReport(cst.KV{
				"msg": fmt.Sprintf("proc(%s/%d), callee: %s, breaker is open and requests dropped\nlast errors:\n%s",
					proc.ProcessName(), proc.Pid(), cb.name, cb.errWin),
				"skip": skip,
			})
		})
	}
	return err
}
