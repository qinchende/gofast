// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sqlx

import (
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
)

func (conn *OrmDB) Insert(obj orm.OrmStruct) int64 {
	obj.BeforeSave() // 设置值
	ts, values := orm.SchemaValues(obj)

	autoIdx := ts.AutoIndex()
	if autoIdx > 0 {
		values[autoIdx] = values[0]
	}

	ret := conn.ExecSql(conn.InsertSql(ts), values[1:]...)
	obj.AfterInsert(ret) // 反写值，比如主键ID
	ct, err := ret.RowsAffected()
	panicIfSqlErr(err)
	return ct
}

func (conn *OrmDB) Delete(obj any) int64 {
	ts := orm.Schema(obj)
	val := ts.PrimaryValue(obj)
	ret := conn.ExecSql(conn.DeleteSql(ts), val)
	return parseSqlResult(conn, ret, val, ts)
}

func (conn *OrmDB) Update(obj orm.OrmStruct) int64 {
	obj.BeforeSave()
	ts, values := orm.SchemaValues(obj)

	fLen := len(values)
	priIdx := ts.PrimaryIndex()
	tVal := values[priIdx]
	values[priIdx] = values[fLen-1]
	values[fLen-1] = tVal

	ret := conn.ExecSql(conn.UpdateSql(ts), values...)
	return parseSqlResult(conn, ret, tVal, ts)
}

// 通过给定的结构体字段更新数据
func (conn *OrmDB) UpdateFields(obj orm.OrmStruct, fNames ...string) int64 {
	dstVal := reflect.Indirect(reflect.ValueOf(obj))
	ts := orm.Schema(obj)

	obj.BeforeSave()
	upSQL, tValues := conn.UpdateSqlByFields(ts, &dstVal, fNames...)
	ret := conn.ExecSql(upSQL, tValues...)
	return parseSqlResult(conn, ret, tValues[len(tValues)-1], ts)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 对应ID值的一行记录
func (conn *OrmDB) QueryPrimary(obj any, id any) int64 {
	ts := orm.Schema(obj)
	return queryByPrimary(conn, obj, id, ts)
}

// 对应 PrimaryKey（一般是ID）值的一行记录，支持行记录缓存
func (conn *OrmDB) QueryPrimaryCache(obj any, id any) int64 {
	return queryByPrimaryWithCache(conn, obj, id)
}

// 查询一行记录，查询条件自定义
func (conn *OrmDB) QueryRow(obj any, where string, args ...any) int64 {
	return conn.QueryRow2(obj, "*", where, args...)
}

func (conn *OrmDB) QueryRow2(obj any, fields string, where string, args ...any) int64 {
	ts := orm.Schema(obj)
	sqlRows := conn.QuerySql(conn.SelectSqlOfOne(ts, fields, where), args...)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(obj, sqlRows, ts)
}

// 自定义SQL语句查询，得到一条记录。或者只取第一条记录的第一个字段值
func (conn *OrmDB) QuerySqlRow(obj any, sql string, args ...any) int64 {
	sqlRows := conn.QuerySql(sql, args...)
	defer CloseSqlRows(sqlRows)
	return scanSqlRowsOne(obj, sqlRows, nil)
}

// 执行类似 select count(*) from table where xxx ，得到执行结果第一行第一个值，值必须可转成 int64
func (conn *OrmDB) QuerySqlInt64(sql string, args ...any) (ct int64) {
	sqlRows := conn.QuerySql(sql, args...)
	defer CloseSqlRows(sqlRows)
	scanSqlRowsOne(&ct, sqlRows, nil) // 不成功就会抛异常
	return
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 查询多行记录
func (conn *OrmDB) QueryRows(objs any, where string, args ...any) int64 {
	return conn.QueryRows2(objs, "*", where, args...)
}

func (conn *OrmDB) QueryRows2(objs any, fields string, where string, args ...any) int64 {
	ts := orm.Schema(objs)
	sqlRows := conn.QuerySql(conn.SelectSqlOfSome(ts, fields, where), args...)
	defer CloseSqlRows(sqlRows)

	return scanSqlRowsList(objs, sqlRows)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 高级查询，可以自定义更多参数
func (conn *OrmDB) QueryPet(pet *SelectPet) int64 {
	ts := orm.SchemaNil(pet.List)
	sql := pet.Sql
	if sql == "" {
		conn.ReadyForSql(ts, pet)
		sql = conn.SelectSqlByPet(pet)
	}

	ct, _ := queryByPet(conn, sql, "", pet, ts)
	return ct
}

// 分页版本，更方便用于数据查询管理
func (conn *OrmDB) QueryPetPaging(pet *SelectPet) (int64, int64) {
	ts := orm.Schema(pet.List)
	sql := pet.Sql
	if sql == "" {
		conn.ReadyForSql(ts, pet)
		sql = conn.SelectPagingSqlByPet(pet)
	}

	sqlCt := pet.SqlCount
	if sqlCt == "" {
		conn.ReadyForSql(ts, pet)
		sqlCt = conn.SelectCountSqlByPet(pet)
	} else if strings.ToLower(sqlCt) == "false" { // 不查total，用于无级分页
		sqlCt = ""
	}

	return queryByPet(conn, sql, sqlCt, pet, ts)
}

func (conn *OrmDB) DeletePetCache(pet *SelectPet) (err error) {
	ts := orm.Schema(pet.List)

	// 生成Sql语句
	sql := pet.Sql
	if sql == "" {
		conn.ReadyForSql(ts, pet)
		sql = conn.SelectSqlByPet(pet)
	}

	pet.Args = formatArgs(pet.Args)
	pet.cacheKey = ts.CacheSqlKey(conn.Attrs.DbName, realSql(sql, pet.Args...))

	key := pet.cacheKey
	rds := (*conn.rdsNodes)[0]
	_, err = rds.Del(key)
	return
}
