package orm

type ApplyOrmStruct interface {
	TableName() string
	BeforeSave()
}
