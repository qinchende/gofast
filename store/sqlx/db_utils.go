package sqlx

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
)

var (
	ErrNotMatchDestination  = errors.New("not matching destination to scan")
	ErrNotReadableValue     = errors.New("value not addressable or interfaceable")
	ErrNotSettable          = errors.New("passed in variable is not settable")
	ErrUnsupportedValueType = errors.New("unsupported unmarshal type")
)

func ErrPanic(err error) {
	if err != nil {
		logx.Stack(err.Error())
		panic(err)
	}
}

func ErrLog(err error) {
	if err != nil {
		logx.Stack(err.Error())
	}
}

func CloseSqlRows(rows *sql.Rows) {
	ErrLog(rows.Close())
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func ScanRow(dest any, sqlRows *sql.Rows) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	sm := orm.SchemaOfType(dstVal.Type())
	return scanSqlRowsOne(&dstVal, sqlRows, sm)
}

func ScanRows(dest any, sqlRows *sql.Rows) int64 {
	dSliceTyp, dItemType, isPtr, isKV := checkDestType(dest)
	sm := orm.SchemaOfType(dItemType)
	return scanSqlRowsSlice(dest, sqlRows, sm, dSliceTyp, dItemType, isPtr, isKV)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func parseSqlResult(ret sql.Result, keyVal any, conn *OrmDB, sm *orm.ModelSchema) int64 {
	ct, err := ret.RowsAffected()
	ErrLog(err)

	// 判断是否要删除缓存，删除缓存的逻辑要特殊处理，
	// TODO：删除Key要有策略，比如删除之后加一个删除标记，后面设置缓存策略先查询这个标记，如果有标记就删除标记但本次不设置缓存
	if ct > 0 && sm.CacheAll() {
		// 目前只支持第一个redis实例作缓存
		if conn.rdsNodes != nil {
			key := fmt.Sprintf(sm.CachePreFix(), conn.Attrs.DbName, keyVal)
			rds := (*conn.rdsNodes)[0]
			_, _ = rds.Del(key)
			_, _ = rds.SetEX(key+"_del", "1", sm.ExpireDuration())
		}
	}

	return ct
}

func queryByIdWithCache(conn *OrmDB, dest any, id any) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(dest))
	sm := orm.SchemaOfType(dstVal.Type())

	// TODO：获取缓存的值
	key := fmt.Sprintf(sm.CachePreFix(), conn.Attrs.DbName, id)
	rds := (*conn.rdsNodes)[0]
	cValStr, err := rds.Get(key)
	if err == nil && cValStr != "" {
		if err = jsonx.UnmarshalFromString(dest, cValStr); err == nil {
			return 1
		}
	}

	// TODO: 执行SQL查询并设置缓存
	sqlRows := conn.QuerySql(selectSqlForID(sm), id)
	defer CloseSqlRows(sqlRows)
	ct := scanSqlRowsOne(&dstVal, sqlRows, sm)
	if ct > 0 {
		keyDel := key + "_del"
		if cValStr, _ := rds.Get(keyDel); cValStr == "1" {
			_, _ = rds.Del(keyDel)
		} else if jsonValBytes, err := jsonx.Marshal(dest); err == nil {
			_, _ = rds.Set(key, jsonValBytes, sm.ExpireDuration())
		}
	}
	return ct
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanSqlRowsOne(dstVal *reflect.Value, sqlRows *sql.Rows, sm *orm.ModelSchema) int64 {
	if !sqlRows.Next() {
		if err := sqlRows.Err(); err != nil {
			ErrLog(err)
		} else {
			ErrLog(sql.ErrNoRows)
		}
		return 0
	}

	rte := dstVal.Type()
	switch rte.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		if dstVal.CanSet() {
			err := sqlRows.Scan(dstVal.Interface())
			ErrPanic(err)
		} else {
			ErrPanic(ErrNotSettable)
		}
	case reflect.Struct:
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
		ErrPanic(err)
	default:
		ErrPanic(ErrUnsupportedValueType)
	}
	return 1
}

// 解析查询到的数据记录
func scanSqlRowsSlice(dest any, sqlRows *sql.Rows, sm *orm.ModelSchema, dSliceTyp reflect.Type, dItemType reflect.Type, isPtr bool, isKV bool) int64 {
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
			ErrPanic(err)

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
			ErrPanic(err)

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

func realSql(sqlStr string, args ...any) string {
	return fmt.Sprintf(strings.ReplaceAll(sqlStr, "?", "%#v"), args...)
}
