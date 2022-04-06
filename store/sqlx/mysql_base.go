package sqlx

import (
	"database/sql"
	"github.com/qinchende/gofast/logx"
)

func (conn *MysqlORM) QuerySql(sql string, args ...interface{}) *sql.Rows {
	logx.DebugPrint(sql)

	if conn.tx == nil {
		rows, err := conn.Reader.Query(sql, args...)
		errPanic(err)
		return rows
	} else {
		rows, err := conn.tx.Query(sql, args...)
		errPanic(err)
		return rows
	}
}

func (conn *MysqlORM) QuerySqlContext(sql string, args ...interface{}) *sql.Rows {
	logx.DebugPrint(sql)

	if conn.tx == nil {
		rows, err := conn.Reader.QueryContext(conn.Ctx, sql, args...)
		errPanic(err)
		return rows
	} else {
		rows, err := conn.tx.QueryContext(conn.Ctx, sql, args...)
		errPanic(err)
		return rows
	}
}

func (conn *MysqlORM) Exec(sql string, args ...interface{}) sql.Result {
	logx.DebugPrint(sql)

	if conn.tx == nil {
		result, err := conn.Writer.Exec(sql, args...)
		errPanic(err)
		return result
	} else {
		result, err := conn.tx.Exec(sql, args...)
		errPanic(err)
		return result
	}
}

func (conn *MysqlORM) ExecContext(sql string, args ...interface{}) sql.Result {
	logx.DebugPrint(sql)

	if conn.tx == nil {
		result, err := conn.Writer.ExecContext(conn.Ctx, sql, args...)
		errPanic(err)
		return result
	} else {
		result, err := conn.tx.ExecContext(conn.Ctx, sql, args...)
		errPanic(err)
		return result
	}
}

func (conn *MysqlORM) Begin() *MysqlORM {
	tx, err := conn.Writer.Begin()
	errPanic(err)
	return &MysqlORM{tx: tx}
}

func (conn *MysqlORM) Commit() error {
	return conn.tx.Commit()
}

func (conn *MysqlORM) Rollback() error {
	return conn.tx.Rollback()
}
