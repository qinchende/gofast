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
func (conn *MysqlORM) QueryID(dest interface{}, id interface{}) int {
	rVal := reflect.Indirect(reflect.ValueOf(dest))

	schema := orm.SchemaOfType(rVal.Type())
	rows := conn.QueryRaw(selectSqlByID(schema), id)
	defer rows.Close()

	if !rows.Next() {
		return 0
	}

	dbColumns, _ := rows.Columns()
	smColumns := schema.ColumnsKV()
	fieldsAddr := make([]interface{}, len(dbColumns))

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
	err := rows.Scan(fieldsAddr...)
	errPanic(err)
	return 1
}

func (conn *MysqlORM) QueryWhere(dest interface{}, condition string, pms ...interface{}) {
	dSliceTyp, dItemType, isPtr := checkQueryType(dest)

	schema := orm.SchemaOfType(dItemType)
	sqlRows := conn.QueryRaw(selectSqlByCondition(schema, condition), pms...)
	defer sqlRows.Close()

	dbColumns, _ := sqlRows.Columns()
	smColumns := schema.ColumnsKV()

	valuesAddr := make([]interface{}, len(dbColumns))
	tpItems := make([]reflect.Value, 0, 25)
	for sqlRows.Next() {
		itemPtr := reflect.New(dItemType)
		itemVal := reflect.Indirect(itemPtr)

		// 每一个db-column都应该有对应的变量接收值
		for cIdx, column := range dbColumns {
			// TODO：这里可以优化，不用每次map查找，而是只查一次，然后缓存index关系
			idx, ok := smColumns[column]
			if ok {
				valuesAddr[cIdx] = schema.AddrByIndex(&itemVal, idx)
			} else if valuesAddr[cIdx] == nil {
				valuesAddr[cIdx] = new(interface{}) // 这个值会被丢弃
			}
		}

		err := sqlRows.Scan(valuesAddr...)
		errPanic(err)

		if isPtr {
			tpItems = append(tpItems, itemPtr)
		} else {
			tpItems = append(tpItems, itemVal)
		}
	}

	records := reflect.MakeSlice(dSliceTyp, 0, len(tpItems))
	records = reflect.Append(records, tpItems...)
	reflect.ValueOf(dest).Elem().Set(records)
}

//
//func (conn *MysqlORM) QueryFields(obj orm.ApplyOrmStruct, fields string, condition string, values ...interface{}) {
//	schema := orm.Schema(obj)
//	rows := conn.QueryRaw(selectSqlByFields(schema, fields, condition), values)
//	defer rows.Close()
//
//	//rVal := reflect.Indirect(reflect.ValueOf(obj))
//	//rTpe := rVal.Type()
//	//
//	//newObj := reflect.New(rTpe)
//}

//func (conn *MysqlORM) QueryDemo(typ reflect.Type) {
//
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Utils
func checkQueryType(dest interface{}) (reflect.Type, reflect.Type, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		panic("dest must be pointer.")
	}
	dSliceTyp := dTyp.Elem()
	if dSliceTyp.Kind() != reflect.Slice {
		panic("dest must be slice.")
	}

	isPtr := false
	dItemType := dSliceTyp.Elem()
	if dItemType.Kind() == reflect.Ptr {
		isPtr = true
		dItemType = dItemType.Elem()
	}

	return dSliceTyp, dItemType, isPtr
}

//func (conn *MysqlORM) QuerySql(sql string, args ...interface{}) {
//
//}
