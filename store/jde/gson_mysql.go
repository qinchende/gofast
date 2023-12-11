package jde

import (
	"database/sql"
	"github.com/qinchende/gofast/core/pool"
	"github.com/qinchende/gofast/cst"
	"strconv"
)

// 直接将数据库查询结果转换成GsonStr
func EncodeGsonRowsFromSqlRows(sqlRows *sql.Rows, tt int64) (ct int64, str string) {
	dbColumns, _ := sqlRows.Columns()
	clsLen := len(dbColumns)
	scanValues := make([]any, clsLen)

	// 先计算数据库
	bf := pool.GetBytesNormal()
	defer pool.FreeBytes(bf)

	for sqlRows.Next() {
		ct++
		for i := 0; i < clsLen; i++ {
			scanValues[i] = &scanValues[i]
		}
		cst.PanicIfErr(sqlRows.Scan(scanValues...))
		encGsonRowFromValues(bf, scanValues)
	}

	bf2 := pool.GetBytesNormal()
	defer pool.FreeBytes(bf2)

	tp := *bf2
	tp = append(tp, '[')

	// 0. 当前记录数量
	tp = append(tp, strconv.FormatInt(int64(ct), 10)...)
	tp = append(tp, ',')
	// 1. 总记录数量
	tp = append(tp, strconv.FormatInt(tt, 10)...)
	tp = append(tp, ",["...)

	// 2. 字段
	for i := 0; i < len(dbColumns); i++ {
		if i != 0 {
			tp = append(tp, ',')
		}
		tp = append(tp, '"')
		tp = append(tp, dbColumns[i]...)
		tp = append(tp, '"')
	}
	tp = append(tp, "],["...)

	// 3. 记录值
	tp = append(tp, (*bf)...)
	if ct > 0 {
		tp = tp[:len(tp)-1]
	}

	// 4. 收尾
	tp = append(tp, "]]"...)
	str = string(tp)
	return
}
