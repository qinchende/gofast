package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/skill/stringx"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"strings"
)

func insertSql(mss *orm.ModelSchema) string {
	return mss.InsertSQL(func(ms *orm.ModelSchema) string {
		cls := ms.Columns()
		clsLen := len(cls)

		sBuf := strings.Builder{}
		sBuf.Grow(256)
		bVal := make([]byte, (clsLen-1)*2-1)

		priIdx := ms.PrimaryIndex()
		ct := 0
		for i := 1; i < clsLen; i++ {
			if ct > 0 {
				sBuf.WriteByte(',')
				bVal[ct] = ','
				ct++
			}
			// 写第一个字段值
			if priIdx == int8(i) {
				sBuf.WriteString(cls[0])
			} else {
				sBuf.WriteString(cls[i])
			}

			bVal[ct] = '?'
			ct++
		}
		return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", ms.TableName(), sBuf.String(), stringx.BytesToString(bVal))
	})
}

func deleteSql(mss *orm.ModelSchema) string {
	return mss.DeleteSQL(func(ms *orm.ModelSchema) string {
		return fmt.Sprintf("DELETE FROM %s WHERE %s=?;", ms.TableName(), ms.Columns()[ms.PrimaryIndex()])
	})
}

func updateSql(mss *orm.ModelSchema) string {
	return mss.UpdateSQL(func(ms *orm.ModelSchema) string {
		cls := ms.Columns()
		clsLen := len(cls) - 1
		sBuf := strings.Builder{}
		sBuf.Grow(256)

		priIdx := ms.PrimaryIndex()
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
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s=?;", ms.TableName(), sBuf.String(), cls[priIdx])
	})
}

// 更新特定字段
func updateSqlByFields(ms *orm.ModelSchema, rVal *reflect.Value, fNames ...string) (string, []any) {
	if len(fNames) == 1 {
		fNames = strings.Split(fNames[0], ",")
	}

	tgLen := len(fNames)
	if tgLen <= 0 {
		panic("sqlx: UpdateByFields args [fNames] is empty")
	}

	flsKV := ms.FieldsKV()
	cls := ms.Columns()
	sBuf := strings.Builder{}
	tValues := make([]any, tgLen+2)

	for i := 0; i < tgLen; i++ {
		idx, ok := flsKV[fNames[i]]
		if !ok {
			fst.GFPanicErr(fmt.Errorf("Field %s not exist.", fNames[i]))
		}

		// 更新字符串
		if i > 0 {
			sBuf.WriteByte(',')
		}
		sBuf.WriteString(cls[idx])
		sBuf.WriteString("=?")

		// 值
		tValues[i] = ms.ValueByIndex(rVal, idx)
	}

	// 更新字段
	upIdx := ms.UpdatedIndex()
	priIdx := ms.PrimaryIndex()
	if upIdx >= 0 {
		sBuf.WriteByte(',')
		sBuf.WriteString(cls[upIdx])
		sBuf.WriteString("=?")
		tValues[tgLen] = ms.ValueByIndex(rVal, upIdx)
		tValues[tgLen+1] = ms.ValueByIndex(rVal, priIdx)
	} else {
		tValues[tgLen] = ms.ValueByIndex(rVal, priIdx)
		tValues = tValues[:tgLen+1]
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s=?;", ms.TableName(), sBuf.String(), cls[priIdx]), tValues
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 查询 select * from

func selectSqlForID(mss *orm.ModelSchema) string {
	return mss.SelectSQL(func(ms *orm.ModelSchema) string {
		return fmt.Sprintf("SELECT * FROM %s WHERE %s=? LIMIT 1;", ms.TableName(), ms.Columns()[ms.PrimaryIndex()])
	})
}

func selectSqlForOne(mss *orm.ModelSchema, fields string, where string) string {
	if fields == "" {
		fields = "*"
	}
	if where == "" {
		where = "1=1"
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s LIMIT 1;", fields, mss.TableName(), where)
}

func selectSqlForSome(mss *orm.ModelSchema, fields string, where string) string {
	if fields == "" {
		fields = "*"
	}
	if where == "" {
		where = "1=1"
	}
	if strings.Index(where, "limit") < 0 {
		where += " LIMIT 10000" // 最多1万条记录
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s;", fields, mss.TableName(), where)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func checkPet(mss *orm.ModelSchema, pet *SelectPet) {
	if pet.Table == "" {
		pet.Table = mss.TableName()
	}
	if pet.Columns == "" {
		pet.Columns = "*"
	}
	if pet.Limit <= 0 {
		pet.Limit = 10000
	}
	if pet.Offset < 0 {
		pet.Offset = 0
	}
	if pet.Where == "" {
		pet.Where = "1=1"
	}
	pet.orderByT = ""
	if pet.OrderBy != "" {
		pet.orderByT = "ORDER BY " + pet.OrderBy
	}

	if pet.Page <= 0 {
		pet.Page = 1
	}
	if pet.PageSize <= 0 {
		pet.PageSize = 10
	}
}

func selectSqlForPet(mss *orm.ModelSchema, pet *SelectPet) string {
	checkPet(mss, pet)
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s %s LIMIT %d OFFSET %d;", pet.Columns, pet.Table, pet.Where, pet.orderByT, pet.Limit, pet.Offset)
}

func selectCountSqlForPet(mss *orm.ModelSchema, pet *SelectPet) string {
	return fmt.Sprintf("SELECT COUNT(*) AS COUNT FROM %s WHERE %s;", pet.Table, pet.Where)
}

func selectPagingSqlForPet(mss *orm.ModelSchema, pet *SelectPet) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s %s LIMIT %d OFFSET %d;", pet.Columns, pet.Table, pet.Where, pet.orderByT, pet.PageSize, (pet.Page-1)*pet.PageSize)
}
