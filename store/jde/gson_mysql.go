package jde

import (
	"database/sql"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/core/pool"
)

// 直接将数据库查询结果转换成GsonStr
func EncodeGsonRowsFromSqlRows(sqlRows *sql.Rows, tt int64) (ct int64, gsonStr string) {
	dbColumns, _ := sqlRows.Columns()
	clsLen := len(dbColumns)
	scanValues := make([]any, clsLen)

	// 先计算数据库
	bs1 := pool.GetBytes()
	defer pool.FreeBytes(bs1)

	for sqlRows.Next() {
		ct++
		for i := 0; i < clsLen; i++ {
			scanValues[i] = &scanValues[i]
		}
		cst.PanicIfErr(sqlRows.Scan(scanValues...))
		encGsonRowFromValues(bs1, scanValues)
	}

	bs2 := pool.GetBytes()
	defer pool.FreeBytes(bs2)

	ret := encGsonRowsHeader(*bs2, ct, tt, dbColumns)

	// 3. 记录值
	ret = append(ret, (*bs1)...)
	if ct > 0 {
		ret = ret[:len(ret)-1]
	}

	// 4. 收尾
	ret = append(ret, "]]"...)
	gsonStr = string(ret)
	return
}
