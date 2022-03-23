package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/skill/stringx"
	"github.com/qinchende/gofast/store/orm"
	"strings"
)

func insertSql(ms *orm.ModelSchema) string {
	cls := ms.Columns()

	vls := make([]byte, (len(cls)-1)*2-1)
	for i := 0; i < len(cls)-2; i++ {
		vls[2*i] = '?'
		vls[2*i+1] = ','
	}
	vls[len(vls)-1] = '?'

	return fmt.Sprintf("INSERT INTO %s (%s) values (%s);", ms.TableName(), strings.Join(cls[1:], ","), stringx.BytesToString(vls))
}

func updateSql(ms *orm.ModelSchema) string {
	cls := ms.Columns()

	sBuf := strings.Builder{}
	for i := 1; i < len(cls); i++ {
		if i > 1 {
			sBuf.WriteString(", ")
		}
		sBuf.WriteString(cls[i])
		sBuf.WriteString("=?")
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", ms.TableName(), sBuf.String(), cls[0])
}

// 更新特定字段
func updateColumnsSql(ms *orm.ModelSchema, values []interface{}, fields []string) (string, []interface{}) {
	cls := ms.Columns()
	fls := ms.Fields()

	tValues := make([]interface{}, len(fields)+2)
	count := 0
	upIdx := ms.UpdatedIndex()

	sBuf := strings.Builder{}
	// 当前schema有更新时间的字段
	if upIdx >= 0 {
		sBuf.WriteString(cls[upIdx])
		sBuf.WriteString("=?")

		tValues[count] = values[upIdx]
		count++
	}

	var clIndex int8
	for i := 0; i < len(fields); i++ {
		clIndex = fls[fields[i]]
		// 发现更新的字段是主键，跳过去
		if clIndex == 0 {
			continue
		}
		if sBuf.Len() > 0 {
			sBuf.WriteString(", ")
		}
		sBuf.WriteString(cls[clIndex])
		sBuf.WriteString("=?")

		tValues[count] = values[clIndex]
		count++
	}

	// 加入主键ID值
	tValues[count] = values[0]
	count++

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?;", ms.TableName(), sBuf.String(), cls[0]), tValues[:count]
}
