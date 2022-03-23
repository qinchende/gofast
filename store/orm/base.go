package orm

import "database/sql"

type ApplyOrmStruct interface {
	TableName() string
	BeforeSave()
	AfterInsert(sql.Result)
	//AfterUpdate(sql.Result)
}
