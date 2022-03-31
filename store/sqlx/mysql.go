package sqlx

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

func (conn *MysqlORM) Insert(obj orm.ApplyOrmStruct) int64 {
	obj.BeforeSave() // 设置值
	schema, values := orm.SchemaValues(obj)

	priIdx := schema.PrimaryIndex()
	if priIdx > 0 {
		values[priIdx] = values[0]
	}

	ret := conn.Exec(insertSql(schema), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	return parseResult(ret)
}

func (conn *MysqlORM) Delete(obj interface{}) int64 {
	schema := orm.Schema(obj)
	val := schema.PrimaryValue(obj)
	ret := conn.Exec(deleteSql(schema), val)
	return parseResult(ret)
}

func (conn *MysqlORM) Update(obj orm.ApplyOrmStruct) int64 {
	obj.BeforeSave()
	schema, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := schema.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.Exec(updateSql(schema), values...)
	return parseResult(ret)
}

// 通过给定的结构体字段更新数据
func (conn *MysqlORM) UpdateColumns(obj orm.ApplyOrmStruct, columns ...string) int64 {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	schema := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByColumns(schema, &rVal, columns)
	ret := conn.Exec(upSQL, tValues...)
	return parseResult(ret)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *MysqlORM) QueryID(dest interface{}, id interface{}) int64 {
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

func (conn *MysqlORM) QueryWhere(dest interface{}, where string, pms ...interface{}) int64 {
	return conn.QueryColumns(dest, "*", where, pms...)
}

func (conn *MysqlORM) QueryColumns(dest interface{}, fields string, where string, pms ...interface{}) int64 {
	dSliceTyp, dItemType, isPtr := checkQueryType(dest)

	schema := orm.SchemaOfType(dItemType)
	sqlRows := conn.QueryRaw(selectSqlByWhere(schema, fields, where), pms...)
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
	return int64(len(tpItems))
}

func (conn *MysqlORM) QueryPet(dest interface{}, pet *SelectPet) int64 {
	dSliceTyp, dItemType, isPtr := checkQueryType(dest)

	logx.Info(dSliceTyp)
	logx.Info(dItemType)
	logx.Info(isPtr)

	schema := orm.SchemaOfType(dItemType)
	if pet.Select == "" {
		pet.Select = selectSqlByPet(schema, pet)
	}
	sqlRows := conn.QueryRaw(pet.Select, pet.Prams...)
	defer sqlRows.Close()

	//dbColumns, _ := sqlRows.Columns()
	//smColumns := schema.ColumnsKV()

	return 0
}

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

func parseResult(ret sql.Result) int64 {
	ct, err := ret.RowsAffected()
	errLog(err)
	return ct
}

//func (conn *MysqlORM) QuerySql(sql string, args ...interface{}) {
//
//}
