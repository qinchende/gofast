package orm

import "database/sql"

const (
	dbDefPrimaryKeyName = "ID"            // 默认主键的字段名
	dbDefUpdatedKeyName = "UpdatedAt"     // 默认主键的字段名
	dbConfigTag         = "dbc"           // 数据库字段配置tag头
	dbPrimaryKeyFlag    = "primary_field" // 数据库主键tag头中配置值
	dbUpdatedKeyFlag    = "updated_field" // 更新时间

	dbColumnNameTag  = "dbf" // 数据库字段名称，对应的tag
	dbColumnNameTag2 = "pms" // 数据库字段名称，次优先级
)

type ApplyOrmStruct interface {
	TableName() string
	BeforeSave()
	AfterInsert(sql.Result)
	//AfterQuery(*sql.Rows)
}
