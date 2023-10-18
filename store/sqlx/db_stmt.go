package sqlx

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

func (conn *OrmDB) Prepare(sqlStr string, readonly bool) *StmtConn {
	return conn.PrepareCtx(conn.Ctx, sqlStr, readonly)
}

func (conn *OrmDB) PrepareCtx(ctx context.Context, sqlStr string, readonly bool) *StmtConn {
	var stmt *sql.Stmt
	var err error

	if conn.tx == nil {
		if readonly == true {
			stmt, err = conn.Reader.PrepareContext(ctx, sqlStr)
		} else {
			stmt, err = conn.Writer.PrepareContext(ctx, sqlStr)
		}
	} else {
		stmt, err = conn.tx.PrepareContext(ctx, sqlStr)
	}

	panicIfSqlErr(err)
	return &StmtConn{ctx: ctx, stmt: stmt, sqlStr: sqlStr, readonly: readonly}
}

func (conn *StmtConn) Close() {
	panicIfSqlErr(conn.stmt.Close())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *StmtConn) Exec(args ...any) int64 {
	return conn.ExecCtx(conn.ctx, args...)
}

func (conn *StmtConn) ExecCtx(ctx context.Context, args ...any) int64 {
	if conn.readonly {
		cst.PanicString("StmtConn just readonly, can't exec sql.")
		return 0
	}

	args = formatArgs(args)
	if logx.ShowDebug() {
		logx.Debug(realSql(conn.sqlStr, args...))
	}
	startTime := timex.Now()
	ret, err := conn.stmt.ExecContext(ctx, args...)
	dur := timex.NowDiff(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] slow-call - %s", dur/time.Millisecond, realSql(conn.sqlStr, args...))
	}

	panicIfSqlErr(err)
	ct, _ := ret.RowsAffected()
	return ct
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *StmtConn) QueryRow(dest any, args ...any) int64 {
	return conn.QueryRowCtx(conn.ctx, dest, args...)
}

func (conn *StmtConn) QueryRowCtx(ctx context.Context, dest any, args ...any) int64 {
	sqlRows, err := conn.queryContext(ctx, args...)
	defer CloseSqlRows(sqlRows)
	panicIfSqlErr(err)
	return scanSqlRowsOne(dest, sqlRows, nil)
}

func (conn *StmtConn) QueryRows(dest any, args ...any) int64 {
	return conn.QueryRowsCtx(conn.ctx, dest, args...)
}

func (conn *StmtConn) QueryRowsCtx(ctx context.Context, dest any, args ...any) int64 {
	sqlRows, err := conn.queryContext(ctx, args...)
	defer CloseSqlRows(sqlRows)
	panicIfSqlErr(err)
	return scanSqlRowsSlice(dest, sqlRows, nil)
}

func (conn *StmtConn) queryContext(ctx context.Context, args ...any) (sqlRows *sql.Rows, err error) {
	args = formatArgs(args)
	if logx.ShowDebug() {
		logx.Debug(realSql(conn.sqlStr, args...))
	}
	startTime := timex.Now()
	sqlRows, err = conn.stmt.QueryContext(ctx, args...)
	dur := timex.NowDiff(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] slow-call - %s", dur/time.Millisecond, realSql(conn.sqlStr, args...))
	}
	return
}
