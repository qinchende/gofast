package fuse

//
//import (
//	"fmt"
//	"github.com/qinchende/gofast/skill/proc"
//)
//
//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//type loggedThrottle struct {
//	name   string
//	errWin *errorWindow
//	throttleNoReason
//}
//
//func newLoggedThrottle(name string, t throttleNoReason) loggedThrottle {
//	return loggedThrottle{
//		name:             name,
//		throttleNoReason: t,
//		errWin:           new(errorWindow),
//	}
//}
//
//func (lt loggedThrottle) allow() (Promise, error) {
//	promise, err := lt.throttleNoReason.allow()
//	return myPromise{
//		promise: promise,
//		errWin:  lt.errWin,
//	}, lt.logError(err)
//}
//
//func (lt loggedThrottle) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
//	return lt.logError(lt.throttleNoReason.doReq(req, fallback, func(err error) bool {
//		accept := acceptable(err)
//		if !accept {
//			lt.errWin.add(err.Error())
//		}
//		return accept
//	}))
//}
//
//func (lt loggedThrottle) logError(err error) error {
//	if err == ErrServiceUnavailable {
//		// if circuit open, not possible to have empty error window
//		Report(fmt.Sprintf(
//			"proc(%s/%d), callee: %s, breaker is open and requests dropped\nlast errors:\n%s",
//			proc.ProcessName(), proc.Pid(), lt.name, lt.errWin))
//	}
//
//	return err
//}
