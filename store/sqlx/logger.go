package sqlx

import "github.com/qinchende/gofast/logx"

func errPanic(err error) {
	if err != nil {
		logx.Error(err)
		panic(err)
	}
}

func errLog(err error) {
	if err != nil {
		logx.Error(err)
	}
}
