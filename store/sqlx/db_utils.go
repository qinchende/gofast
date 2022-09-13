package sqlx

import (
	"database/sql"
	"errors"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

var (
	//ErrNotMatchDestination  = errors.New("not matching destination to scan")
	//ErrNotReadableValue     = errors.New("value not addressable or interfaceable")
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
	sm := orm.Schema(dest)
	return scanSqlRowsOne(dest, sqlRows, sm, nil)
}

func ScanRows(dest any, sqlRows *sql.Rows) int64 {
	return scanSqlRowsSlice(dest, sqlRows, nil)
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
			key := sm.CacheLineKey(conn.Attrs.DbName, keyVal)
			rds := (*conn.rdsNodes)[0]
			_, _ = rds.Del(key)
			_, _ = rds.SetEX(key+"_del", "1", sm.ExpireDuration())
		}
	}

	return ct
}

func queryByIdWithCache(conn *OrmDB, dest any, id any) int64 {
	sm := orm.Schema(dest)

	key := sm.CacheLineKey(conn.Attrs.DbName, id)
	rds := (*conn.rdsNodes)[0]
	cacheStr, err := rds.Get(key)
	if err == nil && cacheStr != "" {
		if err = loadRecordFromGsonString(dest, cacheStr, sm); err == nil {
			return 1
		}
	}

	sqlRows := conn.QuerySql(selectSqlForID(sm), id)
	defer CloseSqlRows(sqlRows)

	var gro gsonResultOne
	ct := scanSqlRowsOne(dest, sqlRows, sm, &gro)
	if ct > 0 {
		keyDel := key + "_del"
		if cacheStr, _ := rds.Get(keyDel); cacheStr == "1" {
			_, _ = rds.Del(keyDel)
		} else if jsonValBytes, err := jsonx.Marshal(gro.Row); err == nil {
			_, _ = rds.Set(key, jsonValBytes, sm.ExpireDuration())
		}
	}
	return ct
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanSqlRowsOne(dest any, sqlRows *sql.Rows, sm *orm.ModelSchema, gro *gsonResultOne) int64 {
	if !sqlRows.Next() {
		if err := sqlRows.Err(); err != nil {
			ErrLog(err)
		} else {
			ErrLog(sql.ErrNoRows)
		}
		return 0
	}

	rte := reflect.TypeOf(dest).Elem()
	rve := reflect.Indirect(reflect.ValueOf(dest))
	switch rte.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		if rve.CanSet() {
			ErrPanic(sqlRows.Scan(rve.Addr().Interface()))
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
				fieldsAddr[cIdx] = sm.AddrByIndex(&rve, idx)
			} else {
				// 这个值会被丢弃
				fieldsAddr[cIdx] = new(any)
			}
		}
		err := sqlRows.Scan(fieldsAddr...)
		ErrPanic(err)

		// 返回行记录的值
		if gro != nil {
			gro.hasValue = true
			gro.Cls = sm.Columns()
			gro.Row = make([]any, len(gro.Cls))

			for idx, column := range gro.Cls {
				gro.Row[idx] = sm.ValueByIndex(&rve, smColumns[column])
			}
		}
	default:
		ErrPanic(ErrUnsupportedValueType)
	}
	return 1
}

// 解析查询到的数据记录
// TODO: 如果 dest 不是某个 struct，而是一个值类型的 slice 又如何处理呢？
func scanSqlRowsSlice(dest any, sqlRows *sql.Rows, gr *gsonResult) int64 {
	sm, sliceType, recordType, isPtr, isKV := checkDestType(dest)

	dbColumns, _ := sqlRows.Columns()
	smColumns := sm.ColumnsKV()

	clsLen := len(dbColumns)
	valuesAddr := make([]any, clsLen)
	var tpRecords []reflect.Value

	// 一般来说，我们的分页大小在25左右，即使要扩容，扩容一次到50也差不多了
	if gr != nil {
		gr.Cls = dbColumns
		gr.Rows = make([][]any, 0, 25)
		if gr.onlyGson != true {
			tpRecords = make([]reflect.Value, 0, 25)
		}
	}

	// 接受者如果是KV类型，相当于解析成了JSON格式，而不是具体类型的对象
	if isKV {
		clsType, _ := sqlRows.ColumnTypes()
		for i := 0; i < clsLen; i++ {
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

			if gr != nil {
				values := make([]any, len(valuesAddr))
				copy(values, valuesAddr)
				gr.Rows = append(gr.Rows, values)

				if gr.onlyGson == true {
					continue
				}
			}

			// 每条记录就是一个类JSON的 KV 对象
			record := make(map[string]any, clsLen)
			for i := 0; i < clsLen; i++ {
				record[dbColumns[i]] = reflect.ValueOf(valuesAddr[i]).Elem().Interface()
			}
			tpRecords = append(tpRecords, reflect.ValueOf(record))
		}
	} else {
		clsPos := make([]int8, clsLen)
		for i := 0; i < clsLen; i++ {
			idx, ok := smColumns[dbColumns[i]]
			if ok {
				clsPos[i] = idx
			} else {
				clsPos[i] = -1
				valuesAddr[i] = new(any)
			}
		}

		for sqlRows.Next() {
			recordPtr := reflect.New(recordType)
			recordVal := reflect.Indirect(recordPtr)

			for i := 0; i < clsLen; i++ {
				if clsPos[i] >= 0 {
					valuesAddr[i] = sm.AddrByIndex(&recordVal, clsPos[i])
				}
			}

			err := sqlRows.Scan(valuesAddr...)
			ErrPanic(err)

			if gr != nil {
				values := make([]any, len(valuesAddr))
				copy(values, valuesAddr)
				gr.Rows = append(gr.Rows, values)

				if gr.onlyGson == true {
					continue
				}
			}

			if isPtr {
				tpRecords = append(tpRecords, recordPtr)
			} else {
				tpRecords = append(tpRecords, recordVal)
			}
		}
	}

	if gr != nil {
		gr.Ct = int64(len(gr.Rows))

		if gr.onlyGson == true {
			return gr.Ct
		}
	}

	records := reflect.MakeSlice(sliceType, 0, len(tpRecords))
	records = reflect.Append(records, tpRecords...)
	reflect.ValueOf(dest).Elem().Set(records)
	return int64(len(tpRecords))
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Utils
func checkDestType(dest any) (*orm.ModelSchema, reflect.Type, reflect.Type, bool, bool) {
	dTyp := reflect.TypeOf(dest)
	if dTyp.Kind() != reflect.Ptr {
		panic("dest must be pointer.")
	}
	sliceType := dTyp.Elem()
	if sliceType.Kind() != reflect.Slice {
		panic("dest must be slice.")
	}
	sm := orm.SchemaOfType(dTyp)

	isPtr := false
	isKV := false
	recordType := sliceType.Elem()
	// 推荐: dest 传入的 slice 类型为指针类型，这样将来就不涉及变量值拷贝了。
	if recordType.Kind() == reflect.Ptr {
		isPtr = true
		recordType = recordType.Elem()
	} else if recordType.Name() == "KV" || recordType.Name() == "cst.KV" || recordType.Name() == "fst.KV" {
		isKV = true
	}

	return sm, sliceType, recordType, isPtr, isKV
}
