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
	obj.BeforeSave() // 设置值
	schema, values := orm.SchemaValues(obj)
	insSQL := insertSql(schema)
	logx.DebugPrint(insSQL)

	ret := conn.Exec(insSQL, values[1:]...)
	obj.AfterInsert(ret) // 反写值
	return ret
}

func (conn *MysqlORM) Update(obj orm.ApplyOrmStruct) sql.Result {
	obj.BeforeSave()
	schema, values := orm.SchemaValues(obj)

	upSQL := updateSql(schema)
	logx.DebugPrint(upSQL)

	// primaryKey value as the last value
	idVal := values[0]
	copy(values, values[1:])
	values[len(values)-1] = idVal

	return conn.Exec(upSQL, values...)
}

//// 更新特定字段
//func (conn *MysqlORM) UpdateColumns(obj orm.ApplyOrmStruct, fields ...interface{}) sql.Result {
//	obj.BeforeSave()
//	schema, values := orm.SchemaValues(obj)
//
//	upSQL, tValues := updateColumnsSql(schema, values, fields)
//	logx.DebugPrint(upSQL)
//
//	return conn.Exec(upSQL, tValues...)
//}

func (conn *MysqlORM) UpdateByNames(obj orm.ApplyOrmStruct, names ...string) sql.Result {
	obj.BeforeSave()
	schema, values := orm.SchemaValues(obj)

	upSQL, tValues := updateColumnsSql(schema, values, names)
	logx.DebugPrint(upSQL)

	return conn.Exec(upSQL, tValues...)
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
