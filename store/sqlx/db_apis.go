package sqlx

import (
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
	"time"
)

func (conn *OrmDB) Insert(obj orm.OrmStruct) int64 {
	obj.BeforeSave() // 设置值
	ms, values := orm.SchemaValues(obj)

	autoIdx := ms.AutoIndex()
	if autoIdx > 0 {
		values[autoIdx] = values[0]
	}

	ret := conn.ExecSql(insertSql(ms), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	ct, err := ret.RowsAffected()
	panicIfErr(err)
	return ct
}

func (conn *OrmDB) Delete(obj any) int64 {
	ms := orm.Schema(obj)
	val := ms.PrimaryValue(obj)
	ret := conn.ExecSql(deleteSql(ms), val)
	return parseSqlResult(ret, val, conn, ms)
}

func (conn *OrmDB) Update(obj orm.OrmStruct) int64 {
	obj.BeforeSave()
	ms, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := ms.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.ExecSql(updateSql(ms), values...)
	return parseSqlResult(ret, tVal, conn, ms)
}

// 通过给定的结构体字段更新数据
func (conn *OrmDB) UpdateFields(obj orm.OrmStruct, fNames ...string) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(obj))
	ms := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := updateSqlByFields(ms, &dstVal, fNames...)
	ret := conn.ExecSql(upSQL, tValues...)
	return parseSqlResult(ret, tValues[len(tValues)-1], conn, ms)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对应ID值的一行记录
func (conn *OrmDB) QueryPrimary(dest any, id any) int64 {
	ms := orm.Schema(dest)
	sqlRows := conn.QuerySql(selectSqlForPrimary(ms), id)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(dest, sqlRows, ms, nil)
}

// 对应ID值的一行记录，支持行记录缓存
func (conn *OrmDB) QueryPrimaryCache(dest any, id any) int64 {
	return queryByPrimaryWithCache(conn, dest, id)
}

// 查询一行记录，查询条件自定义
func (conn *OrmDB) QueryRow(dest any, where string, args ...any) int64 {
	return conn.QueryRow2(dest, "*", where, args...)
}

func (conn *OrmDB) QueryRow2(dest any, fields string, where string, args ...any) int64 {
	ms := orm.Schema(dest)
	sqlRows := conn.QuerySql(selectSqlForOne(ms, fields, where), args...)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(dest, sqlRows, ms, nil)
}

// 自定义SQL语句查询，得到一条记录。或者只取第一条记录的第一个字段值
func (conn *OrmDB) QueryRowSql(dest any, sql string, args ...any) int64 {
	sqlRows := conn.QuerySql(sql, args...)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(dest, sqlRows, nil, nil)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 查询多行记录
func (conn *OrmDB) QueryRows(dest any, where string, args ...any) int64 {
	return conn.QueryRows2(dest, "*", where, args...)
}

func (conn *OrmDB) QueryRows2(dest any, fields string, where string, args ...any) int64 {
	ms := orm.Schema(dest)
	sqlRows := conn.QuerySql(selectSqlForSome(ms, fields, where), args...)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsSlice(dest, sqlRows, nil)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 高级查询，可以自定义更多参数
func (conn *OrmDB) QueryPet(pet *SelectPet) int64 {
	ms := orm.Schema(pet.Target)

	fillPet(ms, pet)
	// 生成Sql语句
	sql := pet.Sql
	if sql == "" {
		sql = selectSqlForPet(ms, pet)
	}

	ct, _ := conn.innerQueryPet(sql, "", pet, ms)
	return ct
}

// 返回 count , total
func (conn *OrmDB) innerQueryPet(sql, sqlCount string, pet *SelectPet, ms *orm.ModelSchema) (int64, int64) {
	withCache := pet.Cache != nil && pet.Cache.ExpireS > 0
	gsonStr := pet.Result != nil && pet.Result.GsonStr == true

	var gr *gsonResult
	rds := (*conn.rdsNodes)[0]
	// 1. 需要走缓存版本
	if withCache {
		pet.Args = formatArgs(pet.Args)
		pet.Cache.sqlHash = ms.CacheSqlKey(realSql(sql, pet.Args...))

		cacheStr, err := rds.Get(pet.Cache.sqlHash)
		if err == nil && cacheStr != "" {
			if gsonStr {
				pet.Result.Target = cacheStr
				return 1, 0
			}

			gr = new(gsonResult)
			panicIfErr(loadRecordsFromGsonString(pet.Target, cacheStr, gr))
			return gr.Ct, gr.Tt
		}
		gr = new(gsonResult)
	}
	if gsonStr {
		if gr == nil {
			gr = new(gsonResult)
		}
		gr.onlyGson = true
	}

	// 2. 执行SQL查询并设置缓存
	var tt int64
	if sqlCount != "" {
		// 此条件下一共多少条
		sqlRows1 := conn.QuerySql(sqlCount, pet.Args...)
		defer CloseSqlRows(sqlRows1)
		scanSqlRowsOne(&tt, sqlRows1, ms, nil)

		if tt <= 0 {
			return 0, 0
		}
		if gr != nil {
			gr.Tt = tt
		}
	}

	sqlRows := conn.QuerySql(sql, pet.Args...)
	defer CloseSqlRows(sqlRows)
	ct := scanSqlRowsSlice(pet.Target, sqlRows, gr)

	var err error
	if gsonStr {
		ret, err := jsonx.Marshal(gr.Gson)
		panicIfErr(err)
		pet.Result.Target = lang.BytesToString(ret)
	}
	if ct > 0 && withCache {
		cacheStr := new(any)
		if gsonStr {
			*cacheStr = pet.Result.Target
		} else {
			*cacheStr, err = jsonx.Marshal(gr.Gson)
			panicIfErr(err)
		}
		_, _ = rds.Set(pet.Cache.sqlHash, *cacheStr, time.Duration(pet.Cache.ExpireS)*time.Second)
	}
	return ct, tt
}

func (conn *OrmDB) QueryPetPaging(pet *SelectPet) (int64, int64) {
	ms := orm.Schema(pet.Target)

	fillPet(ms, pet)
	sqlCt := pet.SqlCount
	if sqlCt == "" {
		sqlCt = selectCountSqlForPet(ms, pet)
	} else if strings.ToLower(sqlCt) == "false" {
		sqlCt = ""
	}
	sql := pet.Sql
	if sql == "" {
		sql = selectPagingSqlForPet(ms, pet)
	}

	return conn.innerQueryPet(sql, sqlCt, pet, ms)
}

func (conn *OrmDB) DeletePetCache(pet *SelectPet) (err error) {
	ms := orm.Schema(pet.Target)

	fillPet(ms, pet)
	// 生成Sql语句
	sql := pet.Sql
	if sql == "" {
		sql = selectSqlForPet(ms, pet)
	}

	pet.Args = formatArgs(pet.Args)
	pet.Cache.sqlHash = ms.CacheSqlKey(realSql(sql, pet.Args...))

	key := pet.Cache.sqlHash
	rds := (*conn.rdsNodes)[0]
	_, err = rds.Del(key)
	return
}
