package sqlx

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/store/orm"
)

type MysqlORM struct {
	Client *sql.DB
	Ctx    context.Context
}

func (conn *MysqlORM) Insert(obj orm.ApplyOrmStruct) sql.Result {
	obj.BeforeSave()
	schema, values := orm.SchemaValues(obj)
	insSQL := insertSql(schema)
	logx.DebugPrint(insSQL)
	return conn.Exec(insSQL, values[1:]...)
}

func (conn *MysqlORM) Update(obj orm.ApplyOrmStruct) sql.Result {
	return nil
}

func (conn *MysqlORM) UpdateColumns(obj orm.ApplyOrmStruct, fields ...string) sql.Result {
	obj.BeforeSave()
	schema, values := orm.SchemaValues(obj)
	logx.Info(schema)
	logx.Info(values)
	return nil
}

//func (conn *MysqlORM) Select(obj interface{}) sql.Result {
//
//}
//
//func (conn *MysqlORM) Delete(obj interface{}) sql.Result {
//
//}

//// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//func (conn *MysqlORM) Query() {
//
//}
//
//func (conn *MysqlORM) Exec() {
//
//}
