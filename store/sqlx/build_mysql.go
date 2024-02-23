// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/aid/lang"
	"github.com/qinchende/gofast/core/cst"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
)

type MysqlBuilder struct {
}

func (*MysqlBuilder) Insert(ts *orm.TableSchema) string {
	return ts.InsertSQL(func(ts *orm.TableSchema) string {
		cls := ts.Columns()
		clsLen := len(cls)

		sBuf := strings.Builder{}
		sBuf.Grow(256)
		bVal := make([]byte, (clsLen-1)*2-1)

		// insert 时 auto increment 字段不需要赋值。我们和需要赋值的字段调换位置
		autoIdx := ts.AutoIndex()
		ct := 0
		for i := 1; i < clsLen; i++ {
			if ct > 0 {
				sBuf.WriteByte(',')
				bVal[ct] = ','
				ct++
			}
			// 写第一个字段值
			if autoIdx == int8(i) {
				sBuf.WriteString(cls[0])
			} else {
				sBuf.WriteString(cls[i])
			}

			bVal[ct] = '?'
			ct++
		}
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", ts.TableName(), sBuf.String(), lang.BTS(bVal))
	})
}

func (*MysqlBuilder) Delete(ts *orm.TableSchema) string {
	return ts.DeleteSQL(func(ts *orm.TableSchema) string {
		return fmt.Sprintf("DELETE FROM %s WHERE %s=?;", ts.TableName(), ts.Columns()[ts.PrimaryIndex()])
	})
}

func (*MysqlBuilder) Update(ts *orm.TableSchema) string {
	return ts.UpdateSQL(func(ts *orm.TableSchema) string {
		cls := ts.Columns()
		clsLen := len(cls) - 1
		sBuf := strings.Builder{}
		sBuf.Grow(256)

		priIdx := ts.PrimaryIndex()
		for i := 0; i < clsLen; i++ {
			if i > 0 {
				sBuf.WriteByte(',')
			}

			if priIdx == int8(i) {
				sBuf.WriteString(cls[clsLen])
			} else {
				sBuf.WriteString(cls[i])
			}
			sBuf.WriteString("=?")
		}
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s=?;", ts.TableName(), sBuf.String(), cls[priIdx])
	})
}

// 更新特定字段
func (*MysqlBuilder) UpdateColumns(ts *orm.TableSchema, rVal *reflect.Value, columns ...string) (string, []any) {
	// 有可能要更新的字段是用逗号隔开的字符串
	if len(columns) == 1 {
		columns = strings.Split(columns[0], ",")
	}

	tgLen := len(columns)
	if tgLen <= 0 {
		cst.PanicString("sqlx: UpdateColumns args [columns] is empty")
	}

	cls := ts.Columns()
	sBuf := strings.Builder{}
	tValues := make([]any, tgLen+2)

	for i := 0; i < tgLen; i++ {
		idx := ts.ColumnIndex(columns[i])
		if idx < 0 {
			cst.PanicString(fmt.Sprintf("sqlx: Field %s not exist.", columns[i]))
		}

		// 更新字符串
		if i > 0 {
			sBuf.WriteByte(',')
		}
		sBuf.WriteString(cls[idx])
		sBuf.WriteString("=?")

		// 值
		tValues[i] = ts.ValueByIndex(rVal, int8(idx))
	}

	// 更新字段
	upIdx := ts.UpdatedIndex()
	priIdx := ts.PrimaryIndex()
	if upIdx >= 0 {
		sBuf.WriteByte(',')
		sBuf.WriteString(cls[upIdx])
		sBuf.WriteString("=?")
		tValues[tgLen] = ts.ValueByIndex(rVal, upIdx)
		tValues[tgLen+1] = ts.ValueByIndex(rVal, priIdx)
	} else {
		tValues[tgLen] = ts.ValueByIndex(rVal, priIdx)
		tValues = tValues[:tgLen+1]
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s=?;", ts.TableName(), sBuf.String(), cls[priIdx]), tValues
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 查询 select * from
func (*MysqlBuilder) SelectPrimary(ts *orm.TableSchema) string {
	return ts.SelectSQL(func(ts *orm.TableSchema) string {
		return fmt.Sprintf("SELECT * FROM %s WHERE %s=? LIMIT 1;", ts.TableName(), ts.Columns()[ts.PrimaryIndex()])
	})
}

func (*MysqlBuilder) SelectRow(ts *orm.TableSchema, columns string, where string) string {
	if columns == "" {
		columns = "*"
	}
	if where == "" {
		where = "1=1"
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1;", columns, ts.TableName(), where)
}

func (*MysqlBuilder) SelectRows(ts *orm.TableSchema, columns string, where string) string {
	if columns == "" {
		columns = "*"
	}
	if where == "" {
		where = "1=1"
	}
	if strings.Index(where, "limit") < 0 {
		where += " LIMIT 10000" // 最多1万条记录
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s;", columns, ts.TableName(), where)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (*MysqlBuilder) InitPet(pet *SelectPet, ts *orm.TableSchema) {
	if pet.isReady {
		return
	}

	if pet.Table == "" {
		pet.Table = ts.TableName()
	}
	if pet.Columns == "" {
		pet.Columns = "*"
	}
	if pet.Limit == 0 {
		pet.Limit = 100
	}
	//if pet.Offset == 0 {
	//	pet.Offset = 0
	//}
	if pet.Where == "" {
		pet.Where = "1=1"
	}
	pet.orderByT = ""
	if pet.OrderBy != "" {
		pet.orderByT = " ORDER BY " + pet.OrderBy
	}
	pet.groupByT = ""
	if pet.GroupBy != "" {
		pet.groupByT = " GROUP BY " + pet.GroupBy
	}
	if pet.Page == 0 {
		pet.Page = 1
	}
	if pet.PageSize == 0 {
		pet.PageSize = 100
	}

	pet.isReady = true
}

func (*MysqlBuilder) SelectByPet(pet *SelectPet) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s%s%s LIMIT %d OFFSET %d;", pet.Columns, pet.Table, pet.Where, pet.groupByT,
		pet.orderByT, pet.Limit, pet.Offset)
}

func (*MysqlBuilder) SelectCountByPet(pet *SelectPet) string {
	if pet.GroupBy == "" {
		return fmt.Sprintf("SELECT COUNT(*) AS COUNT FROM %s WHERE %s;", pet.Table, pet.Where)
	}
	return fmt.Sprintf("SELECT COUNT(DISTINCT(%s)) AS COUNT FROM %s WHERE %s;", pet.GroupBy, pet.Table, pet.Where)
}

func (*MysqlBuilder) SelectPagingByPet(pet *SelectPet) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s%s%s LIMIT %d OFFSET %d;", pet.Columns, pet.Table, pet.Where, pet.groupByT,
		pet.orderByT, pet.PageSize, (pet.Page-1)*pet.PageSize)
}
