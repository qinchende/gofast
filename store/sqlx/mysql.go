package sqlx

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

type MysqlORM struct {
	Client *sql.DB
	Ctx    context.Context
}

func (conn *MysqlORM) Insert(obj orm.ApplyOrmStruct) sql.Result {
	obj.BeforeSave() // 设置值
	schema, values := orm.SchemaValues(obj)

	priIdx := schema.PrimaryIndex()
	if priIdx > 0 {
		values[priIdx] = values[0]
	}

	ret := conn.Exec(insertSql(schema), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	return ret
}

func (conn *MysqlORM) Delete(obj orm.ApplyOrmStruct) sql.Result {
	schema := orm.Schema(obj)
	val := schema.PrimaryValue(obj)
	return conn.Exec(deleteSql(schema), val)
}

func (conn *MysqlORM) Update(obj orm.ApplyOrmStruct) sql.Result {
	obj.BeforeSave()
	schema, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := schema.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	return conn.Exec(updateSql(schema), values...)
}

// 通过给定的结构体字段更新数据
func (conn *MysqlORM) UpdateColumns(obj orm.ApplyOrmStruct, fields ...string) sql.Result {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	schema := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByFields(schema, &rVal, fields)
	return conn.Exec(upSQL, tValues...)
}

//// 不推荐这种方式: 1. 可能参数是值传递， 2. 反射取字段名称，3. 更新值不是传入参数的值，有歧义
//func (conn *MysqlORM) UpdateFields(obj orm.ApplyOrmStruct, fields ...interface{}) sql.Result {
//	fLen := len(fields)
//	names := make([]string, fLen)
//	for i := 0; i < fLen; i++ {
//		va := reflect.Indirect(reflect.ValueOf(fields[i]))
//		names[i] = va.Type().Name()
//	}
//
//	return conn.UpdateByNames(obj, names...)
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *MysqlORM) QueryID(obj orm.ApplyOrmStruct, id interface{}) {
	schema := orm.Schema(obj)
	rows := conn.QueryRaw(selectSqlByID(schema), id)
	defer rows.Close()

	smColumns := schema.ColumnsKV()
	dbColumns, _ := rows.Columns()
	fieldsAddr := make([]interface{}, len(dbColumns))

	rVal := reflect.Indirect(reflect.ValueOf(obj))
	// 每一个db-column都应该有对应的变量接收值
	for cIdx, column := range dbColumns {
		idx, ok := smColumns[column]
		if ok {
			fieldsAddr[cIdx] = schema.AddrByIndex(&rVal, idx)
		} else {
			// 这个值会被丢弃
			fieldsAddr[cIdx] = new(interface{})
		}
	}

	if rows.Next() {
		err := rows.Scan(fieldsAddr...)
		if err != nil {
			panic(err)
		}
	}
}

func (conn *MysqlORM) QueryWhere(obj orm.ApplyOrmStruct, condition string, values ...interface{}) []interface{} {
	schema := orm.Schema(obj)
	rows := conn.QueryRaw(selectSqlByCondition(schema, condition), values...)
	defer rows.Close()

	smColumns := schema.ColumnsKV()
	dbColumns, _ := rows.Columns()
	fieldsAddr := make([]interface{}, len(dbColumns))
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	rTpe := rVal.Type()

	rets := make([]interface{}, 0)
	for rows.Next() {
		newObj := reflect.Indirect(reflect.New(rTpe))

		// 每一个db-column都应该有对应的变量接收值
		for cIdx, column := range dbColumns {
			idx, ok := smColumns[column]
			if ok {
				fieldsAddr[cIdx] = schema.AddrByIndex(&newObj, idx)
			} else if fieldsAddr[cIdx] == nil {
				// 这个值会被丢弃
				fieldsAddr[cIdx] = new(interface{})
			}
		}

		err := rows.Scan(fieldsAddr...)
		if err != nil {
			panic(err)
		}

		rets = append(rets, newObj.Interface())
	}

	return rets
}

func (conn *MysqlORM) QueryFields(obj orm.ApplyOrmStruct, fields string, condition string, values ...interface{}) {
	schema := orm.Schema(obj)
	rows := conn.QueryRaw(selectSqlByFields(schema, fields, condition), values)
	defer rows.Close()

	//rVal := reflect.Indirect(reflect.ValueOf(obj))
	//rTpe := rVal.Type()
	//
	//newObj := reflect.New(rTpe)
}

func (conn *MysqlORM) QueryDemo(typ reflect.Type) {

}

//func (conn *MysqlORM) QuerySql(sql string, args ...interface{}) {
//
//}
