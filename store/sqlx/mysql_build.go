package sqlx

import (
	"fmt"
	"github.com/qinchende/gofast/skill/bytesconv"
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

	return fmt.Sprintf("INSERT INTO %s (%s) values (%s);", ms.TableName(),
		strings.Join(cls[1:], ","), bytesconv.BytesToString(vls))
}
