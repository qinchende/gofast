package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/hashx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
	"time"
)

func panicIfSqlErr(err error) {
	if err != nil {
		logx.Error("sqlx: " + err.Error())
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

func checkDestType(dest any) (*orm.TableSchema, reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Pointer {
		cst.PanicString("Target object must be pointer.")
	}
	sliceType := dTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		cst.PanicString("Target object must be slice.")
	}
	ts := orm.SchemaByType(dTyp)

	isPtr := false
	isKV := false
	recordType := sliceType.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if recordType.Kind() == reflect.Pointer {
		isPtr = true
		recordType = recordType.Elem()
	} else {
		typName := recordType.Name()
		// Note: 不要小看这里的if-else直接比较，往往比很多调用库函数高效多了
		if typName == "cst.KV" || typName == "KV" {
			isKV = true
		}
	}

	return ts, sliceType, recordType, isPtr, isKV
}
