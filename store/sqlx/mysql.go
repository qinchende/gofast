package sqlx

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"time"
)

func (conn *MysqlORM) Insert(obj orm.OrmStruct) int64 {
	obj.BeforeSave() // 设置值
	sm, values := orm.SchemaValues(obj)

	priIdx := sm.PrimaryIndex()
	if priIdx > 0 {
		values[priIdx] = values[0]
	}

	ret := conn.Exec(insertSql(sm), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	ct, err := ret.RowsAffected()
	errLog(err)
	return ct
}

func (conn *MysqlORM) Delete(obj any) int64 {
	sm := orm.Schema(obj)
	val := sm.PrimaryValue(obj)
	ret := conn.Exec(deleteSql(sm), val)
	return parseResult(ret, val, conn, sm)
}

func (conn *MysqlORM) Update(obj orm.OrmStruct) int64 {
	obj.BeforeSave()
	sm, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := sm.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.Exec(updateSql(sm), values...)
	return parseResult(ret, tVal, conn, sm)
}

// 通过给定的结构体字段更新数据
func (conn *MysqlORM) UpdateColumns(obj orm.OrmStruct, columns ...string) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(obj))
	sm := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByColumns(sm, &dstVal, columns)
	ret := conn.Exec(upSQL, tValues...)
	return parseResult(ret, tValues[len(tValues)-1], conn, sm)
}

func parseResult(ret sql.Result, keyVal any, conn *MysqlORM, sm *orm.ModelSchema) int64 {
	ct, err := ret.RowsAffected()
	errLog(err)

	// 判断是否要删除缓存
	if ct > 0 && sm.CacheAll() {
		// 目前只支持第一个redis实例作缓存
		if conn.rdsNodes != nil {
			key := fmt.Sprintf(sm.CachePreFix(), conn.Attrs.DBName, keyVal)
			_, _ = (*conn.rdsNodes)[0].Del(key)
		}
	}

	return ct
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *MysqlORM) QueryID(dest any, id any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(dstVal.Type())
	sqlRows := conn.QuerySql(selectSqlByID(sm), id)
	defer sqlRows.Close()

	return parseQueryRow(&dstVal, sqlRows, sm)
}

// TODO: 当前重点解决缓存的问题
func (conn *MysqlORM) QueryIDCC(dest any, id any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	sm := orm.SchemaOfType(dstVal.Type())

	key := fmt.Sprintf(sm.CachePreFix(), conn.Attrs.DBName, id)
	cValStr, err := (*conn.rdsNodes)[0].Get(key)
	if err == nil && cValStr != "" {
		if err = jsonx.UnmarshalFromString(dest, cValStr); err == nil {
			return 1
		}
	}

	// 执行SQL查询并设置缓存
	sqlRows := conn.QuerySql(selectSqlByID(sm), id)
	defer sqlRows.Close()
	ct := parseQueryRow(&dstVal, sqlRows, sm)
	if ct > 0 {
		if jsonValBytes, err := jsonx.Marshal(dest); err == nil {
			_, _ = (*conn.rdsNodes)[0].Set(key, jsonValBytes, time.Duration(sm.ExpireS())*time.Second)
			//logx.Info(str, err)
		}
	}
	return ct
}

func (conn *MysqlORM) QueryRow(dest any, where string, pms ...any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))

	sm := orm.SchemaOfType(dstVal.Type())
	sqlRows := conn.QuerySql(selectSqlByOne(sm, where), pms...)
	defer sqlRows.Close()

	return parseQueryRow(&dstVal, sqlRows, sm)
}

func parseQueryRow(dstVal *reflect.Value, sqlRows *sql.Rows, sm *orm.ModelSchema) int64 {
	if !sqlRows.Next() {
		return 0
	}

	dbColumns, _ := sqlRows.Columns()
	smColumns := sm.ColumnsKV()
	fieldsAddr := make([]any, len(dbColumns))

	// 每一个db-column都应该有对应的变量接收值
	for cIdx, column := range dbColumns {
		idx, ok := smColumns[column]
		if ok {
			fieldsAddr[cIdx] = sm.AddrByIndex(dstVal, idx)
		} else {
			// 这个值会被丢弃
			fieldsAddr[cIdx] = new(any)
		}
	}
	err := sqlRows.Scan(fieldsAddr...)
	errPanic(err)
	return 1
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (conn *MysqlORM) QueryRows(dest any, where string, pms ...any) int64 {
	return conn.QueryRows2(dest, "*", where, pms...)
}

func (conn *MysqlORM) QueryRows2(dest any, fields string, where string, pms ...any) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)

	sm := orm.SchemaOfType(dItemType)
	sqlRows := conn.QuerySql(selectSqlByWhere(sm, fields, where), pms...)
	defer sqlRows.Close()

	return parseQueryRows(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}

// 高级查询，可以自定义更多参数
func (conn *MysqlORM) QueryPet(dest any, pet *SelectPet) int64 {
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

//func (conn *MysqlORM) QueryPetCC(dest interface{}, pet *SelectPetCC) int64 {
//	return 0
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 解析查询到的数据记录
func parseQueryRows(dest any, sqlRows *sql.Rows, sm *orm.ModelSchema,
	dSliceTyp reflect.Type, dItemType reflect.Type, isPtr bool, isKV bool) int64 {
	dbColumns, _ := sqlRows.Columns()
	smColumns := sm.ColumnsKV()

	dbClsLen := len(dbColumns)
	valuesAddr := make([]any, dbClsLen)
	tpItems := make([]reflect.Value, 0, 25)
	// 接受者如果是KV类型，相当于解析成了JSON格式，而不是具体类型的对象
	if isKV {
		// TODO：可以通过 sqlRows.ColumnsType() 进一步确定字段的类型
		clsType, _ := sqlRows.ColumnTypes()
		for i := 0; i < dbClsLen; i++ {
			typ := clsType[i].ScanType()
			if typ.String() == "sql.RawBytes" {
				valuesAddr[i] = new(string)
			} else {
				valuesAddr[i] = new(any)
			}
		}

		for sqlRows.Next() {
			err := sqlRows.Scan(valuesAddr...)
			errPanic(err)

			obj := make(map[string]any, dbClsLen)
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
				valuesAddr[i] = new(any)
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
func checkDestType(dest any) (reflect.Type, reflect.Type, bool, bool) {
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
	} else if dItemType.Name() == "KV" {
		isKV = true
	}

	return dSliceTyp, dItemType, isPtr, isKV
}
