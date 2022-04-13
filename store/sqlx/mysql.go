package sqlx

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

func (conn *MysqlORM) Insert(obj orm.ApplyOrmStruct) int64 {
	obj.BeforeSave() // 设置值
	sm, values := orm.SchemaValues(obj)

	priIdx := sm.PrimaryIndex()
	if priIdx > 0 {
		values[priIdx] = values[0]
	}

	ret := conn.Exec(insertSql(sm), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	return parseResult(ret)
}

func (conn *MysqlORM) Delete(obj interface{}) int64 {
	sm := orm.Schema(obj)
	val := sm.PrimaryValue(obj)
	ret := conn.Exec(deleteSql(sm), val)
	return parseResult(ret)
}

func (conn *MysqlORM) Update(obj orm.ApplyOrmStruct) int64 {
	obj.BeforeSave()
	sm, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := sm.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.Exec(updateSql(sm), values...)
	return parseResult(ret)
}

// 通过给定的结构体字段更新数据
func (conn *MysqlORM) UpdateColumns(obj orm.ApplyOrmStruct, columns ...string) int64 {
	rVal := reflect.Indirect(reflect.ValueOf(obj))
	sm := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByColumns(sm, &rVal, columns)
	ret := conn.Exec(upSQL, tValues...)
	return parseResult(ret)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *MysqlORM) QueryID(dest interface{}, id interface{}) int64 {
	rVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(rVal.Type())
	sqlRows := conn.QuerySql(selectSqlByID(sm), id)
	defer sqlRows.Close()

	return parseQueryRow(&rVal, sqlRows, sm)
}

func (conn *MysqlORM) QueryRow(dest interface{}, where string, pms ...interface{}) int64 {
	rVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(rVal.Type())
	sqlRows := conn.QuerySql(selectSqlByOne(sm, where), pms...)
	defer sqlRows.Close()

	return parseQueryRow(&rVal, sqlRows, sm)
}

func parseQueryRow(rVal *reflect.Value, sqlRows *sql.Rows, sm *orm.ModelSchema) int64 {
	if !sqlRows.Next() {
		return 0
	}

	dbColumns, _ := sqlRows.Columns()
	smColumns := sm.ColumnsKV()
	fieldsAddr := make([]interface{}, len(dbColumns))

	// 每一个db-column都应该有对应的变量接收值
	for cIdx, column := range dbColumns {
		idx, ok := smColumns[column]
		if ok {
			fieldsAddr[cIdx] = sm.AddrByIndex(rVal, idx)
		} else {
			// 这个值会被丢弃
			fieldsAddr[cIdx] = new(interface{})
		}
	}
	err := sqlRows.Scan(fieldsAddr...)
	errPanic(err)
	return 1
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *MysqlORM) QueryRows(dest interface{}, where string, pms ...interface{}) int64 {
	return conn.QueryRows2(dest, "*", where, pms...)
}

func (conn *MysqlORM) QueryRows2(dest interface{}, fields string, where string, pms ...interface{}) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)

	sm := orm.SchemaOfType(dItemType)
	sqlRows := conn.QuerySql(selectSqlByWhere(sm, fields, where), pms...)
	defer sqlRows.Close()

	return parseQueryRows(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}

// 高级查询，可以自定义更多参数
func (conn *MysqlORM) QueryPet(dest interface{}, pet *SelectPet) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)

	sm := orm.SchemaOfType(dItemType)
	if pet.Sql == "" {
		pet.Sql = selectSqlByPet(sm, pet)
	}
	sqlRows := conn.QuerySql(pet.Sql, pet.Prams...)
	defer sqlRows.Close()

	return parseQueryRows(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 带缓存版本
func (conn *MysqlORM) QueryIDCC(dest interface{}, id interface{}) int64 {
	rVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(rVal.Type())
	sqlRows := conn.QuerySql(selectSqlByID(sm), id)
	defer sqlRows.Close()

	return parseQueryRow(&rVal, sqlRows, sm)
}

func (conn *MysqlORM) QueryPetCC(dest interface{}, pet *SelectPetCC) int64 {
	return 0
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 解析查询到的数据记录
func parseQueryRows(dest interface{}, sqlRows *sql.Rows, sm *orm.ModelSchema,
	dSliceTyp reflect.Type, dItemType reflect.Type, isPtr bool, isKV bool) int64 {
	dbColumns, _ := sqlRows.Columns()
	smColumns := sm.ColumnsKV()

	dbClsLen := len(dbColumns)
	valuesAddr := make([]interface{}, dbClsLen)
	tpItems := make([]reflect.Value, 0, 25)
	if isKV {
		// TODO：可以通过 sqlRows.ColumnsType() 进一步确定字段的类型
		clsType, _ := sqlRows.ColumnTypes()
		for i := 0; i < dbClsLen; i++ {
			typ := clsType[i].ScanType()
			if typ.String() == "sql.RawBytes" {
				valuesAddr[i] = new(string)
			} else {
				valuesAddr[i] = new(interface{})
			}
		}

		for sqlRows.Next() {
			err := sqlRows.Scan(valuesAddr...)
			errPanic(err)

			obj := make(map[string]interface{}, dbClsLen)
			for i := 0; i < dbClsLen; i++ {
				obj[dbColumns[i]] = reflect.ValueOf(valuesAddr[i]).Elem().Interface()
			}
			tpItems = append(tpItems, reflect.ValueOf(obj))
		}
	} else {
		dbClsIndex := make([]int8, dbClsLen)
		for i := 0; i < dbClsLen; i++ {
			idx, ok := smColumns[dbColumns[i]]
			if ok {
				dbClsIndex[i] = idx
			} else {
				dbClsIndex[i] = -1
				valuesAddr[i] = new(interface{})
			}
		}

		for sqlRows.Next() {
			itemPtr := reflect.New(dItemType)
			itemVal := reflect.Indirect(itemPtr)

			//// 每一个db-column都应该有对应的变量接收值
			//for cIdx, column := range dbColumns {
			//	// TODO：这里可以优化，不用每次map查找，而是只查一次，然后缓存index关系
			//	idx, ok := smColumns[column]
			//	if ok {
			//		valuesAddr[cIdx] = sm.AddrByIndex(&itemVal, idx)
			//	} else if valuesAddr[cIdx] == nil {
			//		valuesAddr[cIdx] = new(interface{}) // 这个值会被丢弃
			//	}
			//}
			for i := 0; i < dbClsLen; i++ {
				if dbClsIndex[i] >= 0 {
					valuesAddr[i] = sm.AddrByIndex(&itemVal, dbClsIndex[i])
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
	}

	records := reflect.MakeSlice(dSliceTyp, 0, len(tpItems))
	records = reflect.Append(records, tpItems...)
	reflect.ValueOf(dest).Elem().Set(records)
	return int64(len(tpItems))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Utils
func checkDestType(dest interface{}) (reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		panic("dest must be pointer.")
	}
	dSliceTyp := dTyp.Elem()
	if dSliceTyp.Kind() != reflect.Slice {
		panic("dest must be slice.")
	}

	isPtr := false
	isKV := false
	dItemType := dSliceTyp.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if dItemType.Kind() == reflect.Ptr {
		isPtr = true
		dItemType = dItemType.Elem()
	} else if dItemType.String() == "fst.KV" {
		isKV = true
	}

	return dSliceTyp, dItemType, isPtr, isKV
}

func parseResult(ret sql.Result) int64 {
	ct, err := ret.RowsAffected()
	errLog(err)
	return ct
}
