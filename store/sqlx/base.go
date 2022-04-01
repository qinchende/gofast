package sqlx

import (
	"context"
	"database/sql"
)

// 天然支持读写分离，只需要数据库连接配置文件，分别传入读写库的连接地址
type MysqlORM struct {
	Reader *sql.DB // 只读连接（从库）
	Writer *sql.DB // 只写连接（主库）
	Client *sql.DB // 读写皆可（主库）
	Ctx    context.Context
}

type SelectPet struct {
	Sql     string
	Table   string
	Columns string
	Offset  int64
	Limit   int64
	Where   string
	Prams   []interface{}
}
