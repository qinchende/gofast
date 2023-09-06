package sqlx

import (
	"database/sql"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
)

func CloseSqlRows(rows *sql.Rows) {
	panicIfErr(rows.Close())
}

func ScanRow(dest any, sqlRows *sql.Rows) int64 {
	return scanSqlRowsOne(dest, sqlRows, nil, nil)
}

func ScanRows(dest any, sqlRows *sql.Rows) int64 {
	return scanSqlRowsSlice(dest, sqlRows, nil)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func parseSqlResult(ret sql.Result, keyVal any, conn *OrmDB, ts *orm.TableSchema) int64 {
	ct, err := ret.RowsAffected()
	panicIfErr(err)

	// 判断是否要删除缓存，删除缓存的逻辑要特殊处理，
	// TODO：删除Key要有策略，比如删除之后加一个删除标记，后面设置缓存策略先查询这个标记，如果有标记就删除标记但本次不设置缓存
	if ct > 0 && ts.CacheAll() {
		// 目前只支持第一个redis实例作缓存
		if conn.rdsNodes != nil {
			key := ts.CacheLineKey(conn.Attrs.DbName, keyVal)
			rds := (*conn.rdsNodes)[0]
			_, _ = rds.Del(key)
			_, _ = rds.SetEX(key+"_del", "1", ts.ExpireDuration())
		}
	}

	return ct
}

func queryByPrimaryWithCache(conn *OrmDB, obj any, id any) int64 {
	ts := orm.Schema(obj)

	key := ts.CacheLineKey(conn.Attrs.DbName, id)
	rds := (*conn.rdsNodes)[0]
	cacheStr, err := rds.Get(key)
	if err == nil && cacheStr != "" {
		if err = bindFromGsonValueString(obj, lang.STB(cacheStr), ts); err == nil {
			return 1
		}
	}

	sqlRows := conn.QuerySql(selectSqlForPrimary(ts), id)
	defer CloseSqlRows(sqlRows)

	var gro gsonResultOne
	ct := scanSqlRowsOne(obj, sqlRows, ts, &gro)
	if ct > 0 {
		keyDel := key + "_del"
		if cacheStr, _ = rds.Get(keyDel); cacheStr == "1" {
			_, _ = rds.Del(keyDel)
		} else if jsonValBytes, err := jsonx.Marshal(gro.Row); err == nil {
			_, _ = rds.Set(key, jsonValBytes, ts.ExpireDuration())
		}
	}
	return ct
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func scanSqlRowsOne(obj any, sqlRows *sql.Rows, ts *orm.TableSchema, gro *gsonResultOne) int64 {
	if !sqlRows.Next() {
		panicIfErr(sqlRows.Err())
		return 0
	}

	dstTyp := reflect.TypeOf(obj).Elem()
	dstVal := reflect.Indirect(reflect.ValueOf(obj))

	// 1. 基础值类型只取第一行第一列值。2. 结构体类型只取第一行数据
	switch dstTyp.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		if dstVal.CanSet() {
			panicIfErr(sqlRows.Scan(dstVal.Addr().Interface()))
		} else {
			cst.PanicString("Passed in variable is not settable.")
		}
	case reflect.Struct:
		if ts == nil {
			ts = orm.Schema(obj)
		}
		dbColumns, _ := sqlRows.Columns()
		fieldsAddr := make([]any, len(dbColumns))

		// 每一个db-column都应该有对应的变量接收值
		for cIdx, column := range dbColumns {
			idx := ts.ColumnIndex(column)
			if idx >= 0 {
				fieldsAddr[cIdx] = ts.AddrByIndex(&dstVal, int8(idx))
			} else {
				fieldsAddr[cIdx] = new(any) // 这个值会被丢弃
			}
		}
		panicIfErr(sqlRows.Scan(fieldsAddr...))

		// 如果需要，返回行记录的值
		if gro != nil {
			gro.hasValue = true
			gro.Cls = ts.Columns()
			gro.Row = make([]any, len(gro.Cls))

			for idx, column := range gro.Cls {
				gro.GsonRow.SetByIndex(idx, ts.ValueByIndex(&dstVal, int8(ts.ColumnIndex(column))))
			}
		}
	default:
		cst.PanicString("Unsupported unmarshal type.")
	}
	return 1
}

// 解析查询到的数据记录
// TODO: 如果 dest 不是某个 struct，而是一个值类型的 slice 又如何处理呢？
func scanSqlRowsSlice(objs any, sqlRows *sql.Rows, gr *gsonResult) int64 {
	ts, sliceType, recordType, isPtr, isKV := checkDestType(objs)

	dbColumns, _ := sqlRows.Columns()
	//msColumns := ts.ColumnsKV()

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
			panicIfErr(sqlRows.Scan(valuesAddr...))

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
			clsPos[i] = int8(ts.ColumnIndex(dbColumns[i]))
			//clsPos[i] = idx
			if clsPos[i] < 0 {
				valuesAddr[i] = new(any)
			}
		}

		for sqlRows.Next() {
			recordPtr := reflect.New(recordType)
			recordVal := reflect.Indirect(recordPtr)

			for i := 0; i < clsLen; i++ {
				if clsPos[i] >= 0 {
					valuesAddr[i] = ts.AddrByIndex(&recordVal, clsPos[i])
				}
			}

			panicIfErr(sqlRows.Scan(valuesAddr...))

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
	reflect.ValueOf(objs).Elem().Set(records)
	return int64(len(tpRecords))
}
