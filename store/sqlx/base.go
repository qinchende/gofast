package sqlx

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/connx/gfrds"
)

const (
	timeFormat     = "2006-01-02 15:04:05"
	timeFormatMini = "01-02 15:04:05"
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
	ctx      context.Context // 传递上下文
	stmt     *sql.Stmt       // 标准库Stmt对象
	sqlStr   string          // 预执行SQL语句
	readonly bool            // 是否连接只读库
}

const (
	CacheMem   uint8 = iota // 默认0：用本地内存缓存，无法实现分布式一致性。但支持存取对象，性能好
	CacheRedis              // 1：强大的redis缓存，支持分布式。需要序列化和反序列化，开销比内存型大
)

// 功能强大的
type SelectPet struct {
	Target   any        // 解析的目标对象数组
	Sql      string     // 自定义完整的SQL语句，注意(Sql和SqlCount是成对出现的)
	SqlCount string     // 分页场景下自定义查询总数的SQL语句（如果传"false"字符串，将不会查询总数）
	Table    string     // 表名
	Columns  string     // 自定义查询字段
	Where    string     // where
	OrderBy  string     // order by
	orderByT string     // order by inner temp
	GroupBy  string     // group by
	groupByT string     // group by inner temp
	Args     []any      // SQL语句参数，防注入
	PageSize uint32     // 分页大小
	Page     uint32     // 当前页
	Offset   uint32     // 查询偏移量
	Limit    uint32     // 查询限量
	Cache    *PetCache  // 缓存设置参数
	Result   *PetResult // 扩展返回数据的形式
	isReady  bool       // 是否已经初始化
}

type PetCache struct {
	sqlHash   string
	ExpireS   uint32 // 过期时间（秒）
	CacheType uint8  // 缓存类型
}

type PetResult struct {
	Target  any
	GsonStr bool // 不解析Target对象，直接返回原始值类型
}

type SP SelectPet
type PC PetCache
