package sqlx

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/store/orm"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (conn *MysqlORM) QueryID(obj orm.ApplyOrmStruct, id interface{}) {
	schema := orm.Schema(obj)
	sql := selectSqlByID(schema)

	rows := conn.QuerySql(sql, id)
	err := rows.Scan(obj)
	if err != nil {
		logx.Info(err)
	}
	logx.Info(obj)
}

func (conn *MysqlORM) QueryWhere(obj orm.ApplyOrmStruct) {

}
