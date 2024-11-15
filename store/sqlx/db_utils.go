// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/aid/hashx"
	"github.com/qinchende/gofast/aid/logx"
	"github.com/qinchende/gofast/core/cst"
	"strings"
	"time"
)

func panicIfSqlErr(err error) {
	if err != nil {
		logx.Err().Msg("sqlx: " + err.Error())
		cst.PanicIfErr(err)
	}
}

// 将SQL参数格式化，方便后面拼接SQL字符串
// 其实就是将所有参数几乎全部转换成数值或者字符串型
func formatArgs(args []any) []any {
	for i, v := range args {
		switch v.(type) {
		case time.Time:
			args[i] = v.(time.Time).Format(cst.TimeFmtYmdHms)
		case *time.Time:
			args[i] = v.(*time.Time).Format(cst.TimeFmtYmdHms)
		}
	}
	return args
}

func realSql(sqlStr string, args ...any) string {
	return fmt.Sprintf(strings.ReplaceAll(sqlStr, "?", "%#v"), args...)
}

func sqlHash(sqlStr string) string {
	return hashx.Md5HexString(sqlStr)
}

func realSqlHash(sqlStr string, args ...any) string {
	sql := realSql(sqlStr, args...)
	return hashx.Md5HexString(sql)
}
