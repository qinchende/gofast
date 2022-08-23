package sqlx

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/connx/gfrds"
)

// 天然支持读写分离，只需要数据库连接配置文件，分别传入读写库的连接地址
type OrmDB struct {
	Attrs    *DBAttrs
	Ctx      context.Context
	Reader   *sql.DB          // 只读连接（从库）
	Writer   *sql.DB          // 只写连接（主库）
	tx       *sql.Tx          // 读写皆可（主库）单独用于处理事务的连接
	rdsNodes *[]gfrds.GfRedis // redis集群用来做缓存的
}

type DBAttrs struct {
	DriverName string // 数据库类型名称
	DbName     string // 数据库名
}

type StmtConn struct {
	ctx      context.Context
	stmt     *sql.Stmt
	sqlStr   string
	readonly bool
}

const (
	CacheMem   uint8 = iota // 默认0：用本地内存缓存，无法实现分布式一致性。但支持存取对象，性能好
	CacheRedis              // 1：强大的redis缓存，支持分布式。需要序列化和反序列化，开销比内存型大
)

type SelectPet struct {
	Sql     string
	Table   string
	Columns string
	Offset  int64
	Limit   int64
	Where   string
	Prams   []any
}

type SelectPetCache struct {
	SelectPet
	ExpireS   uint32 // 过期时间（秒）
	CacheType uint8  // 缓存类型
}

type SP SelectPet
type SPC SelectPetCache
