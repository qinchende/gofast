package sqlx

import (
	"fmt"
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
func updateSqlByFields(ms *orm.ModelSchema, rVal *reflect.Value, fields []string) (string, []interface{}) {
	tgLen := len(fields)
	if tgLen <= 0 {
		panic("UpdateByNames params [names] is empty")
	}

	fls := ms.Fields()
	cls := ms.Columns()
	sBuf := strings.Builder{}
	tValues := make([]interface{}, tgLen+2)

	for i := 0; i < tgLen; i++ {
		idx, ok := fls[fields[i]]
		if !ok {
			panic(fmt.Errorf("field %s not exist", fields[i]))
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

func selectSqlByID(mss *orm.ModelSchema) string {
	return mss.SelectSQL(func(ms *orm.ModelSchema) string {
		return fmt.Sprintf("SELECT * FROM %s WHERE %s=?;", ms.TableName(), ms.Columns()[ms.PrimaryIndex()])
	})
}
