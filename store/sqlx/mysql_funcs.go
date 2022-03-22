package sqlx

import "database/sql"

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func (conn *MysqlORM) Query(sql string, args ...interface{}) *sql.Rows {
	rows, err := conn.Client.Query(sql, args...)
	errPanic(err)
	return rows
}

func (conn *MysqlORM) QueryContext(sql string, args ...interface{}) *sql.Rows {
	rows, err := conn.Client.QueryContext(conn.Ctx, sql, args...)
	errPanic(err)
	return rows
}

func (conn *MysqlORM) Exec(sql string, args ...interface{}) sql.Result {
	result, err := conn.Client.Exec(sql, args...)
	errPanic(err)
	return result
}

func (conn *MysqlORM) ExecContext(sql string, args ...interface{}) sql.Result {
	result, err := conn.Client.ExecContext(conn.Ctx, sql, args...)
	errPanic(err)
	return result
}
