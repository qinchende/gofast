package sqlx

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

// 执行超过500ms的语句需要优化分析，我们先打印出慢日志
const slowThreshold = time.Millisecond * 500

func (conn *OrmDB) Prepare(q string) *sql.Stmt {
	return conn.PrepareCtx(conn.Ctx, q)
}

func (conn *OrmDB) PrepareCtx(ctx context.Context, q string) *sql.Stmt {
	var stmt *sql.Stmt
	var err error
	if conn.tx == nil {
		stmt, err = conn.Writer.PrepareContext(ctx, q)
	} else {
		stmt, err = conn.tx.PrepareContext(ctx, q)
	}
	errPanic(err)
	return stmt
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *OrmDB) Exec(sql string, args ...any) sql.Result {
	return conn.ExecCtx(conn.Ctx, sql, args...)
}

func (conn *OrmDB) ExecCtx(ctx context.Context, sqlStr string, args ...any) sql.Result {
	logx.Debug(sqlStr)

	var result sql.Result
	var err error
	startTime := timex.Now()
	if conn.tx == nil {
		result, err = conn.Writer.ExecContext(ctx, sqlStr, args...)
	} else {
		result, err = conn.tx.ExecContext(ctx, sqlStr, args...)
	}
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] exec: slow-call - %s", dur/time.Millisecond, sqlStr)
	}
	errPanic(err)
	return result
}

func (conn *OrmDB) QuerySql(sql string, args ...any) *sql.Rows {
	return conn.QuerySqlCtx(conn.Ctx, sql, args...)
}

func (conn *OrmDB) QuerySqlCtx(ctx context.Context, sqlStr string, args ...any) *sql.Rows {
	logx.Debug(sqlStr)

	var rows *sql.Rows
	var err error
	startTime := timex.Now()
	if conn.tx == nil {
		rows, err = conn.Reader.QueryContext(ctx, sqlStr, args...)
	} else {
		rows, err = conn.tx.QueryContext(ctx, sqlStr, args...)
	}
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] query: slow-call - %s", dur/time.Millisecond, sqlStr)
	}
	errPanic(err)
	return rows
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *OrmDB) TransBegin() *OrmDB {
	return conn.TransCtx(conn.Ctx)
}

func (conn *OrmDB) TransCtx(ctx context.Context) *OrmDB {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	errPanic(err)
	return &OrmDB{Attrs: conn.Attrs, Ctx: ctx, tx: tx, rdsNodes: conn.rdsNodes}
}

func (conn *OrmDB) TransFunc(fn func(newConn *OrmDB)) {
	conn.TransFuncCtx(conn.Ctx, fn)
}

func (conn *OrmDB) TransFuncCtx(ctx context.Context, fn func(newConn *OrmDB)) {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	errPanic(err)

	nConn := OrmDB{Attrs: conn.Attrs, Ctx: ctx, tx: tx, rdsNodes: conn.rdsNodes}
	defer nConn.TransEnd()
	fn(&nConn)
}

func (conn *OrmDB) Commit() error {
	return conn.tx.Commit()
}

func (conn *OrmDB) Rollback() error {
	return conn.tx.Rollback()
}

func (conn *OrmDB) TransEnd() {
	var err error

	if pic := recover(); pic != nil {
		err = conn.Rollback()
	} else {
		err = conn.Commit()
	}
	if err != nil {
		logx.ErrorF("Terrible mistake. TransEnd error: %s", err)
	}
}
