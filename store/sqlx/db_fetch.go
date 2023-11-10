package sqlx

import (
	"database/sql"
	"fmt"
	"github.com/qinchende/gofast/core/rt"
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/gson"
	"github.com/qinchende/gofast/store/jde"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"time"
	"unsafe"
)

func CloseSqlRows(rows *sql.Rows) {
	panicIfSqlErr(rows.Close())
}

func ScanRow(obj any, sqlRows *sql.Rows) int64 {
	return scanSqlRowsOne(obj, sqlRows, nil)
}

func ScanRows(objs any, sqlRows *sql.Rows) int64 {
	return scanSqlRowsList(objs, sqlRows)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 在修改数据库记录的情况下（delete | update），处理返回结果，同时改变缓存状态
func parseSqlResult(conn *OrmDB, ret sql.Result, keyVal any, ts *orm.TableSchema) int64 {
	ct, err := ret.RowsAffected()
	panicIfSqlErr(err)

	// 判断是否要删除缓存，删除缓存的逻辑要特殊处理，
	// 删除Key要有策略，比如删除之后加一个删除标记，后面设置缓存策略先查询这个标记，如果有标记就删除标记但本次不设置缓存
	if ct > 0 && ts.CacheAll() {
		// 目前只支持第一个redis实例作缓存
		if conn.rdsNodes != nil {
			key := ts.CacheLineKey(conn.Attrs.DbName, keyVal)
			rds := (*conn.rdsNodes)[0]
			// TODO：NOTE: 下面两句必须保证是原子的，才能尽可能避免BUG
			_, _ = rds.Del(key)
			_, _ = rds.SetEX(key+cacheDelFlagSuffix, "1", ts.ExpireDuration())
		}
	}
	return ct
}

// 通过表的主键查询到一条记录。并对单条记录缓存。
// 缓存的数据仅仅为 GsonRow 的 values，而不需要记录 fields ，因为默认都是 按model的字段顺序来记录。
func queryByPrimaryWithCache(conn *OrmDB, obj any, id any) int64 {
	ts := orm.Schema(obj)
	if ts.CacheAll() == false {
		return queryByPrimary(conn, obj, id, ts)
	}

	key := ts.CacheLineKey(conn.Attrs.DbName, id)
	rds := (*conn.rdsNodes)[0]
	cacheStr, err := rds.Get(key)
	if err == nil && cacheStr != "" {
		if err = jde.DecodeGsonRowFromValueString(obj, cacheStr); err == nil {
			return 1
		}
		// Note: 缓存解析失败啥也不管，将再次查询解析并缓存，此时会覆盖旧缓存数据
	}

	ct := queryByPrimary(conn, obj, id, ts)
	if ct > 0 {
		// 先查询缓存删除标记，如果存在标记本次不设置缓存，而且删除标记
		keyDel := key + cacheDelFlagSuffix
		if cacheStr, _ = rds.Get(keyDel); cacheStr == "1" {
			_, _ = rds.Del(keyDel)
		} else if jsonBytes, err2 := jde.EncodeToOnlyGsonRowValuesBytes(obj); err2 == nil {
			_, _ = rds.Set(key, jsonBytes, ts.ExpireDuration())
		}
	}
	return ct
}

// 通过主键查询表记录，同时绑定到对象
func queryByPrimary(conn *OrmDB, obj any, id any, ts *orm.TableSchema) int64 {
	sqlRows := conn.QuerySql(selectSqlOfPrimary(ts), id)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(obj, sqlRows, ts)
}

// 返回 count , total
func queryByPet(conn *OrmDB, sql, sqlCount string, pet *SelectPet, ts *orm.TableSchema) (ct int64, tt int64) {
	// 1. 需要缓存
	if pet.CacheExpireS > 0 {
		rds := (*conn.rdsNodes)[0]
		pet.Args = formatArgs(pet.Args)
		pet.cacheKey = ts.CacheSqlKey(realSql(sql, pet.Args...))

		// 如果有缓存，直接取缓存并解析
		if cacheStr, err := rds.Get(pet.cacheKey); err == nil && cacheStr != "" {
			// A. 直接返回GSON字符串
			if pet.GsonNeed {
				pet.GsonVal = cacheStr
				if pet.GsonOnly {
					return 1, 0 // 不做验证，直接返回缓存中的字符串
				}
			}

			// B. GSON字符串解析成对象
			ret := jde.DecodeGsonRowsFromString(pet.Dest, cacheStr)
			panicIfSqlErr(ret.Err)
			return ret.Ct, ret.Tt
		}
	}

	// 2. 执行SQL查询，必要时设置缓存
	// 先查total, 此条件下一共多少条
	if sqlCount != "" {
		sqlRowsTt := conn.QuerySql(sqlCount, pet.Args...)
		defer CloseSqlRows(sqlRowsTt)
		scanSqlRowsOne(&tt, sqlRowsTt, ts)
		if tt <= 0 {
			return 0, 0
		}
	}

	// 需要 GsonRows 对象
	var grs *gson.GsonRows
	if pet.GsonNeed || pet.CacheExpireS > 0 {
		grs = new(gson.GsonRows)
	}

	sqlRows := conn.QuerySql(sql, pet.Args...)
	defer CloseSqlRows(sqlRows)
	ct = scanSqlRowsListSuper(pet.Dest, sqlRows, grs, pet.GsonOnly)

	var err error
	if pet.GsonNeed {
		pet.GsonVal, err = jde.EncodeToString(grs)
		panicIfSqlErr(err)
	}
	if ct > 0 && pet.CacheExpireS > 0 {
		cacheStr := new(any)
		if pet.GsonNeed {
			*cacheStr = pet.Dest
		} else {
			*cacheStr, err = jsonx.Marshal(grs)
			panicIfSqlErr(err)
		}
		rds := (*conn.rdsNodes)[0]
		_, _ = rds.Set(pet.cacheKey, *cacheStr, time.Duration(pet.CacheExpireS)*time.Second)
	}
	return ct, tt
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 解析单条记录
// Return:
// 1. int64 返回解析到的记录数。只可能是 0 或者 1
func scanSqlRowsOne(obj any, sqlRows *sql.Rows, ts *orm.TableSchema) int64 {
	if !sqlRows.Next() {
		panicIfSqlErr(sqlRows.Err())
	}

	dstVal := reflect.ValueOf(obj)
	if dstVal.Kind() != reflect.Pointer {
		cst.PanicString("Dest must be pointer value.")
	}
	dstVal = reflect.Indirect(dstVal)
	if !dstVal.IsValid() {
		cst.PanicString("Invalid value")
	}
	dstType := dstVal.Type()
	dstKind := dstType.Kind()
	destPtr := (*rt.AFace)(unsafe.Pointer(&obj)).DataPtr

	// NOTE: 当前绑定对象支持三种情况
	// 1. 目标值是结构体类型，只取第一行数据
	if dstKind == reflect.Struct && dstType.String() != "time.Time" {
		if ts == nil {
			ts = orm.SchemaByType(dstType)
		}

		dbColumns, _ := sqlRows.Columns()         // 数据库执行返回字段
		scanValues := make([]any, len(dbColumns)) // 目标值地址

		// Note: 每一个db-column都应该有对应的变量接收值
		for cIdx := range dbColumns {
			fIdx := ts.ColumnIndex(dbColumns[cIdx]) // 查询 db-column 对应struct中字段的索引
			if fIdx >= 0 {
				scanner := ts.SS.FieldsAttr[fIdx].SqlValue
				if scanner != nil {
					scanValues[cIdx] = scanner(destPtr)
				} else {
					scanValues[cIdx] = ts.AddrByIndex(&dstVal, int8(fIdx))
				}
			} else {
				scanValues[cIdx] = sharedAnyValue // 这个值会被丢弃，所以用一个共享的占位变量即可
			}
		}
		panicIfSqlErr(sqlRows.Scan(scanValues...))
		return 1
	}

	//// older +++ modify by cd.net on 20231103
	//if dstKind == reflect.Struct {
	//	if ts == nil {
	//		ts = orm.SchemaByType(dstType)
	//	}
	//
	//	dbColumns, _ := sqlRows.Columns()         // 查询返回的结果字段
	//	scanValues := make([]any, len(dbColumns)) // 目标值地址
	//
	//	// Note: 每一个db-column都应该有对应的变量接收值
	//	for cIdx := range dbColumns {
	//		fIdx := ts.ColumnIndex(dbColumns[cIdx]) // 查询 db-column 对应struct中字段的索引
	//		if fIdx >= 0 {
	//			scanValues[cIdx] = ts.AddrByIndex(&dstVal, int8(fIdx))
	//		} else {
	//			scanValues[cIdx] = sharedAnyValue // 这个值会被丢弃，所以用一个共享的占位变量即可
	//		}
	//	}
	//	panicIfSqlErr(sqlRows.Scan(scanValues...))
	//	return 1
	//}

	// 2. 如果是 KV 类型呢，即目标值只返回 hash 即可
	if dstKind == reflect.Map {
		// 只支持这两种map变量
		typeStr := dstType.String()
		if typeStr != "cst.KV" {
			cst.PanicString(fmt.Sprintf("Unsupported map type of %s.", typeStr))
		}

		// todo: do kv bind
		dbColumns, _ := sqlRows.Columns()         // 数据库执行返回字段
		scanValues := make([]any, len(dbColumns)) // 目标值地址
		kvs := make(cst.KV, len(dbColumns))

		// 根据结果类型做适当的转换
		clsTypes, _ := sqlRows.ColumnTypes()
		for i := range dbColumns {
			typ := clsTypes[i].ScanType()
			if typ.String() == "sql.RawBytes" {
				scanValues[i] = new(string)
			} else {
				scanValues[i] = &scanValues[i]
			}
			kvs[dbColumns[i]] = scanValues[i]
		}
		panicIfSqlErr(sqlRows.Scan(scanValues...))
		dstVal.Set(reflect.ValueOf(kvs))
		return 1
	}

	// 3. 目标值是基础值类型，只取第一行第一列值
	switch dstKind {
	case reflect.Int:
		panicIfSqlErr(sqlRows.Scan(dts.IntValue(destPtr)))
	case reflect.Int8:
		panicIfSqlErr(sqlRows.Scan(dts.Int8Value(destPtr)))
	case reflect.Int16:
		panicIfSqlErr(sqlRows.Scan(dts.Int16Value(destPtr)))
	case reflect.Int32:
		panicIfSqlErr(sqlRows.Scan(dts.Int32Value(destPtr)))
	case reflect.Int64:
		panicIfSqlErr(sqlRows.Scan(dts.Int64Value(destPtr)))

	case reflect.Uint:
		panicIfSqlErr(sqlRows.Scan(dts.UintValue(destPtr)))
	case reflect.Uint8:
		panicIfSqlErr(sqlRows.Scan(dts.Uint8Value(destPtr)))
	case reflect.Uint16:
		panicIfSqlErr(sqlRows.Scan(dts.Uint16Value(destPtr)))
	case reflect.Uint32:
		panicIfSqlErr(sqlRows.Scan(dts.Uint32Value(destPtr)))
	case reflect.Uint64:
		panicIfSqlErr(sqlRows.Scan(dts.Uint64Value(destPtr)))

	case reflect.Float32:
		panicIfSqlErr(sqlRows.Scan(dts.Float32Value(destPtr)))
	case reflect.Float64:
		panicIfSqlErr(sqlRows.Scan(dts.Float64Value(destPtr)))

	case reflect.Bool:
		panicIfSqlErr(sqlRows.Scan(dts.BoolValue(destPtr)))
	case reflect.String:
		panicIfSqlErr(sqlRows.Scan(dts.StringValue(destPtr)))
	case reflect.Struct:
		// 此时只可能是time.Time -> dstType.String() == "time.Time"
		panicIfSqlErr(sqlRows.Scan(dts.TimeValue(destPtr)))
	case reflect.Interface:
		panicIfSqlErr(sqlRows.Scan(dts.AnyValue(destPtr)))

	//case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64,
	//	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	//	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//	if dstVal.CanSet() {
	//		panicIfSqlErr(sqlRows.Scan(dstVal.Addr().Interface()))
	//	} else {
	//		cst.PanicString("Variable can't settable.")
	//	}
	default:
		cst.PanicString("Unsupported unmarshal type.")
	}
	return 1
}

func scanSqlRowsList(objs any, sqlRows *sql.Rows) int64 {
	return scanSqlRowsListSuper(objs, sqlRows, nil, false)
}

// 解析多条记录
// TODO: 如果目标值类型 不是某个 struct，而是一个值类型的 list 又如何处理呢？
// Return:
// 1. int64 返回解析到的记录数 >= 0
func scanSqlRowsListSuper(objs any, sqlRows *sql.Rows, grs *gson.GsonRows, dropBind bool) int64 {
	ts, sliceType, recordType, isPtr, isKV := checkDestType(objs)

	// msColumns := ts.ColumnsKV()
	var tpRows []reflect.Value
	dbColumns, _ := sqlRows.Columns()
	clsLen := len(dbColumns)
	scanValues := make([]any, clsLen)

	// 一般来说，我们的分页大小在25左右，即使要扩容，扩容一次到50也差不多了
	if grs != nil {
		grs.Cls = dbColumns
		grs.Rows = make([][]any, 0, 25)
		if dropBind == false {
			tpRows = make([]reflect.Value, 0, 25)
		}
	}

	// A. 接受者如果是KV类型，相当于解析成了JSON格式，而不是具体类型的对象
	if isKV {
		clsType, _ := sqlRows.ColumnTypes()
		for i := 0; i < clsLen; i++ {
			typ := clsType[i].ScanType()
			// 查询结果绝大部分都是sql.RawBytes，直接解析成 string 类型即可
			if typ.String() == "sql.RawBytes" {
				scanValues[i] = new(string)
			} else {
				scanValues[i] = new(any)
			}
		}

		for sqlRows.Next() {
			panicIfSqlErr(sqlRows.Scan(scanValues...))
			if grs != nil {
				rowValues := make([]any, len(scanValues))
				copy(rowValues, scanValues)
				grs.Rows = append(grs.Rows, rowValues)
			}
			if dropBind == true {
				continue
			}

			// 每条记录就是一个类JSON的 KV 对象
			kv := make(map[string]any, clsLen)
			for i := 0; i < clsLen; i++ {
				kv[dbColumns[i]] = reflect.ValueOf(scanValues[i]).Elem().Interface()
			}
			tpRows = append(tpRows, reflect.ValueOf(kv))
		}
	} else {
		// B. 要么就是某个struct类型
		clsPos := make([]int8, clsLen)
		for i := 0; i < clsLen; i++ {
			clsPos[i] = int8(ts.ColumnIndex(dbColumns[i]))
			if clsPos[i] < 0 {
				scanValues[i] = sharedAnyValue // 这一列的值会被丢弃，所以用一个共享的占位变量即可
			}
		}

		for sqlRows.Next() {
			recordPtr := reflect.New(recordType)
			recordVal := reflect.Indirect(recordPtr)

			for i := 0; i < clsLen; i++ {
				if clsPos[i] >= 0 {
					scanValues[i] = ts.AddrByIndex(&recordVal, clsPos[i])
				}
			}

			panicIfSqlErr(sqlRows.Scan(scanValues...))
			if grs != nil {
				rowValues := make([]any, len(scanValues))
				copy(rowValues, scanValues)
				grs.Rows = append(grs.Rows, rowValues)
			}
			if dropBind == true {
				continue
			}

			if isPtr {
				tpRows = append(tpRows, recordPtr)
			} else {
				tpRows = append(tpRows, recordVal)
			}
		}
	}

	if grs != nil {
		grs.Ct = int64(len(grs.Rows))
	}
	if dropBind == true {
		return grs.Ct
	}

	records := reflect.MakeSlice(sliceType, 0, len(tpRows))
	records = reflect.Append(records, tpRows...)
	reflect.ValueOf(objs).Elem().Set(records)
	return int64(len(tpRows))
}
