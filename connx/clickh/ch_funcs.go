package clickh

import (
	"database/sql"
)

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func (mx *ClickHouseX) Query(sql string, args ...any) *sql.Rows {
	rows, err := mx.Cli.Query(sql, args...)
	errPanic(err)
	return rows
}

func (mx *ClickHouseX) QueryContext(sql string, args ...any) *sql.Rows {
	rows, err := mx.Cli.QueryContext(mx.Ctx, sql, args...)
	errPanic(err)
	return rows
}

func (mx *ClickHouseX) Exec(sql string, args ...any) sql.Result {
	result, err := mx.Cli.Exec(sql, args...)
	errPanic(err)
	return result
}

func (mx *ClickHouseX) ExecContext(sql string, args ...any) sql.Result {
	result, err := mx.Cli.ExecContext(mx.Ctx, sql, args...)
	errPanic(err)
	return result
}
