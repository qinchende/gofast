package sqlx

import (
	"database/sql"
	"github.com/qinchende/gofast/logx"
)

func (conn *MysqlORM) QueryRaw(sql string, args ...interface{}) *sql.Rows {
	logx.DebugPrint(sql)
	rows, err := conn.Client.Query(sql, args...)
	errPanic(err)
	return rows
}

func (conn *MysqlORM) QueryRawContext(sql string, args ...interface{}) *sql.Rows {
	logx.DebugPrint(sql)
	rows, err := conn.Client.QueryContext(conn.Ctx, sql, args...)
	errPanic(err)
	return rows
}

func (conn *MysqlORM) Exec(sql string, args ...interface{}) sql.Result {
	logx.DebugPrint(sql)
	result, err := conn.Client.Exec(sql, args...)
	errPanic(err)
	return result
}

func (conn *MysqlORM) ExecContext(sql string, args ...interface{}) sql.Result {
	logx.DebugPrint(sql)
	result, err := conn.Client.ExecContext(conn.Ctx, sql, args...)
	errPanic(err)
	return result
}
