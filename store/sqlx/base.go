package sqlx

import (
	"context"
	"database/sql"
)

// 天然支持读写分离，只需要数据库连接配置文件，分别传入读写库的连接地址
type MysqlORM struct {
	Ctx    context.Context
	Reader *sql.DB // 只读连接（从库）
	Writer *sql.DB // 只写连接（主库）
	tx     *sql.Tx // 读写皆可（主库）单独用于处理事务的连接
}

//const (
//	ConnWriter uint8 = iota // 默认0：从读
//	ConnReader              // 1：主写
//)

type SelectPet struct {
	Sql     string
	Table   string
	Columns string
	Offset  int64
	Limit   int64
	Where   string
	Prams   []interface{}
}

const (
	CacheMem   uint8 = iota // 默认0：内存
	CacheRedis              // 1：redis
)

type SelectPetCC struct {
	CacheType uint8
	SelectPet
}
