package sqlx

import (
	"context"
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/logx"
)

func (conn *OrmDB) SetRdsNodes(nodes *[]gfrds.GfRedis) {
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

func errPanic(err error) {
	if err != nil {
		logx.Error(err.Error())
		panic(err)
	}
}

func errLog(err error) {
	if err != nil {
		logx.Error(err.Error())
	}
}
