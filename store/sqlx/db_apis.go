package sqlx

import (
	"github.com/qinchende/gofast/cst"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/mapx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"time"
)

func (conn *OrmDB) Insert(obj orm.OrmStruct) int64 {
	obj.BeforeSave() // 设置值
	sm, values := orm.SchemaValues(obj)

	priIdx := sm.PrimaryIndex()
	if priIdx > 0 {
		values[priIdx] = values[0]
	}

	ret := conn.ExecSql(insertSql(sm), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	ct, err := ret.RowsAffected()
	ErrLog(err)
	return ct
}

func (conn *OrmDB) Delete(obj any) int64 {
	sm := orm.Schema(obj)
	val := sm.PrimaryValue(obj)
	ret := conn.ExecSql(deleteSql(sm), val)
	return parseSqlResult(ret, val, conn, sm)
}

func (conn *OrmDB) Update(obj orm.OrmStruct) int64 {
	obj.BeforeSave()
	sm, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := sm.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.ExecSql(updateSql(sm), values...)
	return parseSqlResult(ret, tVal, conn, sm)
}

// 通过给定的结构体字段更新数据
func (conn *OrmDB) UpdateFields(obj orm.OrmStruct, fNames ...string) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(obj))
	sm := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByFields(sm, &dstVal, fNames...)
	ret := conn.ExecSql(upSQL, tValues...)
	return parseSqlResult(ret, tValues[len(tValues)-1], conn, sm)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对应ID值的一行记录
func (conn *OrmDB) QueryID(dest any, id any) int64 {
	sm := orm.Schema(dest)
	sqlRows := conn.QuerySql(selectSqlForID(sm), id)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(dest, sqlRows, sm)
}

// 对应ID值的一行记录，支持行记录缓存
func (conn *OrmDB) QueryIDCache(dest any, id any) int64 {
	return queryByIdWithCache(conn, dest, id)
}

// 查询一行记录，查询条件自定义
func (conn *OrmDB) QueryRow(dest any, where string, args ...any) int64 {
	return conn.QueryRow2(dest, "*", where, args...)
}

func (conn *OrmDB) QueryRow2(dest any, fields string, where string, args ...any) int64 {
	sm := orm.Schema(dest)
	sqlRows := conn.QuerySql(selectSqlForOne(sm, fields, where), args...)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(dest, sqlRows, sm)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 查询多行记录
func (conn *OrmDB) QueryRows(dest any, where string, args ...any) int64 {
	return conn.QueryRows2(dest, "*", where, args...)
}

func (conn *OrmDB) QueryRows2(dest any, fields string, where string, args ...any) int64 {
	sm := orm.Schema(dest)
	sqlRows := conn.QuerySql(selectSqlForSome(sm, fields, where), args...)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsSlice(dest, sqlRows, sm)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 高级查询，可以自定义更多参数
func (conn *OrmDB) QueryPet(pet *SelectPet) int64 {
	sm := orm.Schema(pet.Target)

	// 生成Sql语句
	sql := pet.Sql
	if sql == "" {
		sql = selectSqlForPet(sm, checkPet(sm, pet))
	}

	// 需要走缓存版本
	if pet.PetCache != nil && pet.PetCache.ExpireS > 0 {
		pet.Args = formatArgs(pet.Args)
		pet.sqlHash = sm.CacheSqlKey(realSql(sql, pet.Args...))

		// TODO：获取缓存的值
		key := pet.sqlHash
		rds := (*conn.rdsNodes)[0]
		cValStr, err := rds.Get(key)
		if err == nil && cValStr != "" {
			dest := cst.KV{"count": 0, "records": cst.KV{}}
			if err = jsonx.UnmarshalFromString(&dest, cValStr); err == nil {
				err := mapx.ApplySliceOfConfig(pet.Target, dest["records"])
				ErrPanic(err)
				ct, _ := lang.ToInt64(dest["count"])
				return ct
			}
		}

		// TODO: 执行SQL查询并设置缓存
		sqlRows := conn.QuerySql(sql, pet.Args...)
		defer CloseSqlRows(sqlRows)
		ct := scanSqlRowsSlice(pet.Target, sqlRows, sm)

		if ct > 0 {
			if jsonValBytes, err := jsonx.Marshal(cst.KV{"count": ct, "records": pet.Target}); err == nil {
				_, _ = rds.Set(key, jsonValBytes, time.Duration(pet.ExpireS)*time.Second)
			}
		}
		return ct
	}

	sqlRows := conn.QuerySql(sql, pet.Args...)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsSlice(pet.Target, sqlRows, sm)
}

func (conn *OrmDB) QueryPetPaging(pet *SelectPet) (int64, int64) {
	sm := orm.Schema(pet.Target)

	sqlCount := pet.SqlCount
	if sqlCount == "" {
		sqlCount = selectCountSqlForPet(sm, checkPet(sm, pet))
	}
	sql := pet.Sql
	if sql == "" {
		sql = selectPagingSqlForPet(sm, pet)
	}

	// 此条件下一共多少条
	sqlRows1 := conn.QuerySql(sqlCount, pet.Args...)
	defer CloseSqlRows(sqlRows1)
	// 此条件下的分页记录
	sqlRows2 := conn.QuerySql(sql, pet.Args...)
	defer CloseSqlRows(sqlRows2)

	var total int64
	scanSqlRowsOne(&total, sqlRows1, sm)

	return scanSqlRowsSlice(pet.Target, sqlRows2, sm), total
}

func (conn *OrmDB) DeletePetCache(pet *SelectPet) (err error) {
	sm := orm.Schema(pet.Target)
	// 生成Sql语句
	sql := pet.Sql
	if sql == "" {
		sql = selectSqlForPet(sm, checkPet(sm, pet))
	}

	pet.Args = formatArgs(pet.Args)
	pet.sqlHash = sm.CacheSqlKey(realSql(sql, pet.Args...))

	key := pet.sqlHash
	rds := (*conn.rdsNodes)[0]
	_, err = rds.Del(key)
	return
}
