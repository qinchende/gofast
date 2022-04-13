package sqlx

import (
	"context"
	msql "database/sql"
	"github.com/qinchende/gofast/logx"
)

func (conn *MysqlORM) CloneWithCtx(ctx context.Context) *MysqlORM {
	return &MysqlORM{Ctx: ctx, Writer: conn.Writer, Reader: conn.Reader}
}

func (conn *MysqlORM) Exec(sql string, args ...interface{}) msql.Result {
	return conn.ExecCtx(conn.Ctx, sql, args...)
}
func (conn *MysqlORM) ExecCtx(ctx context.Context, sql string, args ...interface{}) msql.Result {
	logx.DebugPrint(sql)

	var result msql.Result
	var err error
	if conn.tx == nil {
		result, err = conn.Writer.ExecContext(ctx, sql, args...)
	} else {
		result, err = conn.tx.ExecContext(ctx, sql, args...)
	}
	errPanic(err)
	return result
}

func (conn *MysqlORM) QuerySql(sql string, args ...interface{}) *msql.Rows {
	return conn.QuerySqlCtx(conn.Ctx, sql, args...)
}

func (conn *MysqlORM) QuerySqlCtx(ctx context.Context, sql string, args ...interface{}) *msql.Rows {
	logx.DebugPrint(sql)

	var rows *msql.Rows
	var err error
	if conn.tx == nil {
		rows, err = conn.Reader.QueryContext(ctx, sql, args...)
	} else {
		rows, err = conn.tx.QueryContext(ctx, sql, args...)
	}
	errPanic(err)
	return rows
}

func (conn *MysqlORM) Trans() *MysqlORM {
	return conn.TransCtx(conn.Ctx)
}

func (conn *MysqlORM) TransCtx(ctx context.Context) *MysqlORM {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	errPanic(err)
	return &MysqlORM{Ctx: ctx, tx: tx}
}

func (conn *MysqlORM) TransFunc(fn func(newConn *MysqlORM)) {
	conn.TransFuncCtx(conn.Ctx, fn)
}

func (conn *MysqlORM) TransFuncCtx(ctx context.Context, fn func(newConn *MysqlORM)) {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	errPanic(err)

	nConn := MysqlORM{Ctx: ctx, tx: tx}
	defer nConn.TransEnd()
	fn(&nConn)
}

func (conn *MysqlORM) Commit() error {
	return conn.tx.Commit()
}

func (conn *MysqlORM) Rollback() error {
	return conn.tx.Rollback()
}

func (conn *MysqlORM) TransEnd() {
	var err error

	if pic := recover(); pic != nil {
		err = conn.Rollback()
	} else {
		err = conn.Commit()
	}
	if err != nil {
		logx.Errorf("Terrible mistake. trans end error: %s", err)
	}
}
