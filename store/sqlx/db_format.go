package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/skill/hash"
	"strings"
	"time"
)

// 将SQL参数格式化，方便后面拼接SQL字符串
// 其实就是将所有参数几乎全部转换成数值或者字符串型
func formatArgs(args []any) []any {
	for idx, item := range args {
		switch item.(type) {
		case time.Time:
			args[idx] = item.(time.Time).Format(timeFormat)
		case *time.Time:
			args[idx] = item.(*time.Time).Format(timeFormat)
		}
	}
	return args
}

func realSql(sqlStr string, args ...any) string {
	return fmt.Sprintf(strings.ReplaceAll(sqlStr, "?", "%#v"), args...)
}

func sqlHash(sqlStr string) string {
	return hash.Md5HexString(sqlStr)
}

func realSqlHash(sqlStr string, args ...any) string {
	sql := realSql(sqlStr, args...)
	return hash.Md5HexString(sql)
}
