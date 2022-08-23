package sqlx

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"time"
)

func (conn *OrmDB) Prepare(sqlStr string, readonly bool) *StmtConn {
	return conn.PrepareCtx(conn.Ctx, sqlStr, readonly)
}

func (conn *OrmDB) PrepareCtx(ctx context.Context, sqlStr string, readonly bool) *StmtConn {
	logx.Debug(sqlStr)

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
	errPanic(err)
	return &StmtConn{ctx: ctx, stmt: stmt, sqlStr: sqlStr, readonly: readonly}
}

func (conn *StmtConn) Close() error {
	return conn.stmt.Close()
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *StmtConn) Exec(pms ...any) int64 {
	return conn.ExecCtx(conn.ctx, pms...)
}

func (conn *StmtConn) ExecCtx(ctx context.Context, pms ...any) int64 {
	if conn.readonly {
		panic("StmtConn just readonly, can't exec sql.")
		return 0
	}

	startTime := timex.Now()
	ret, err := conn.stmt.ExecContext(ctx, pms...)
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] exec: slow-call - %s", dur/time.Millisecond, conn.sqlStr)
	}

	errPanic(err)
	ct, _ := ret.RowsAffected()
	return ct
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *StmtConn) QueryRow(dest any, pms ...any) int64 {
	return conn.QueryRowCtx(conn.ctx, dest, pms...)
}

func (conn *StmtConn) QueryRowCtx(ctx context.Context, dest any, pms ...any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	sm := orm.SchemaOfType(dstVal.Type())

	startTime := timex.Now()
	sqlRows, err := conn.stmt.QueryContext(ctx, pms...)
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] query: slow-call - %s", dur/time.Millisecond, conn.sqlStr)
	}

	errPanic(err)
	defer sqlRows.Close()
	return parseQueryRow(&dstVal, sqlRows, sm)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *StmtConn) QueryRows(dest any, pms ...any) int64 {
	return conn.QueryRowsCtx(conn.ctx, dest, pms...)
}

func (conn *StmtConn) QueryRowsCtx(ctx context.Context, dest any, pms ...any) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)
	sm := orm.SchemaOfType(dItemType)

	startTime := timex.Now()
	sqlRows, err := conn.stmt.QueryContext(ctx, pms...)
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] query: slow-call - %s", dur/time.Millisecond, conn.sqlStr)
	}

	errPanic(err)
	defer sqlRows.Close()
	return parseQueryRows(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}
