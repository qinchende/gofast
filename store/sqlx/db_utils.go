package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/logx"
	"strings"
)

func ErrPanic(err error) {
	if err != nil {
		logx.Stack(err.Error())
		panic(err)
	}
}

func ErrLog(err error) {
	if err != nil {
		logx.Stack(err.Error())
	}
}

func realSql(sqlStr string, args ...any) string {
	return fmt.Sprintf(strings.ReplaceAll(sqlStr, "?", "%#v"), args...)
}
