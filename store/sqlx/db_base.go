package sqlx

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

func (conn *OrmDB) SetRdsNodes(nodes *[]redis.GfRedis) {
	if len(*nodes) > 0 {
		conn.rdsNodes = nodes
	} else {
		conn.rdsNodes = nil
	}
}

func (conn *OrmDB) CloneWithCtx(ctx context.Context) *OrmDB {
	newConn := *conn
	newConn.Ctx = ctx
	return &newConn
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *OrmDB) ExecSql(sqlStr string, args ...any) sql.Result {
	return conn.ExecSqlCtx(conn.Ctx, sqlStr, args...)
}

func (conn *OrmDB) ExecSqlCtx(ctx context.Context, sqlStr string, args ...any) sql.Result {
	args = formatArgs(args)
	if logx.ShowDebug() {
		logx.Debug(realSql(sqlStr, args...))
	}

	var result sql.Result
	var err error
	startTime := timex.Now()
	if conn.tx == nil {
		result, err = conn.Writer.ExecContext(ctx, sqlStr, args...)
	} else {
		result, err = conn.tx.ExecContext(ctx, sqlStr, args...)
	}
	dur := timex.NowDiff(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] exec: slow-call - %s", dur/time.Millisecond, realSql(sqlStr, args...))
	}
	panicIfSqlErr(err)
	return result
}

func (conn *OrmDB) QuerySql(sqlStr string, args ...any) *sql.Rows {
	return conn.QuerySqlCtx(conn.Ctx, sqlStr, args...)
}

func (conn *OrmDB) QuerySqlCtx(ctx context.Context, sqlStr string, args ...any) *sql.Rows {
	args = formatArgs(args)
	if logx.ShowDebug() {
		logx.Debug(realSql(sqlStr, args...))
	}

	var rows *sql.Rows
	var err error
	startTime := timex.Now()
	if conn.tx == nil {
		rows, err = conn.Reader.QueryContext(ctx, sqlStr, args...)
	} else {
		rows, err = conn.tx.QueryContext(ctx, sqlStr, args...)
	}
	dur := timex.NowDiff(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] query: slow-call - %s", dur/time.Millisecond, realSql(sqlStr, args...))
	}
	panicIfSqlErr(err)
	return rows
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *OrmDB) TransBegin() *OrmDB {
	return conn.TransCtx(conn.Ctx)
}

func (conn *OrmDB) TransCtx(ctx context.Context) *OrmDB {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	panicIfSqlErr(err)
	return &OrmDB{Attrs: conn.Attrs, Ctx: ctx, tx: tx, rdsNodes: conn.rdsNodes}
}

func (conn *OrmDB) TransFunc(fn func(newConn *OrmDB)) {
	conn.TransFuncCtx(conn.Ctx, fn)
}

func (conn *OrmDB) TransFuncCtx(ctx context.Context, fn func(newConn *OrmDB)) {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	panicIfSqlErr(err)

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
	if pic := recover(); pic != nil {
		panicIfSqlErr(conn.Rollback())
		cst.Panic(pic)
	} else {
		panicIfSqlErr(conn.Commit())
	}
	// 出现了非常严重的错误，可能没有提交或回滚成功
}
