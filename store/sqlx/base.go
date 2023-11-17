package sqlx

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/store/dts"
	"time"
)

const (
	slowThreshold      = time.Millisecond * 500 // 执行超过500ms的语句需要优化分析，我们先打印出慢日志
	cacheDelFlagSuffix = "_del_mark"
)

const (
	CacheMem   uint8 = iota // 默认0：用本地内存缓存，无法实现分布式一致性。但支持存取对象，性能好
	CacheRedis              // 1：强大的redis缓存，支持分布式。需要序列化和反序列化，开销比内存型大
)

var (
	sharedAnyValue = new(dts.SqlSkip)
)

// 天然支持读写分离，只需要数据库连接配置文件，分别传入读写库的连接地址
type OrmDB struct {
	Attrs    *DBAttrs
	Ctx      context.Context
	Reader   *sql.DB          // 只读连接（从库）
	Writer   *sql.DB          // 只写连接（主库）
	tx       *sql.Tx          // 读写皆可（主库）单独用于处理事务的连接
	rdsNodes *[]redis.GfRedis // redis集群用来做缓存的
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

// 功能强大的控制参数，这个对象比较大，可以考虑用sync.Pool缓存
type SelectPet struct {
	List any // 解析得到的目标数据，必填项

	// 用于构造SQL语句的参数
	Sql      string // 自定义完整的SQL语句，注意(Sql和SqlCount是成对出现的)
	SqlCount string // 分页场景下自定义查询总数的SQL语句（如果传"false"字符串，将不会查询总数）
	Table    string // 表名
	Columns  string // 自定义查询字段
	Where    string // where
	Args     []any  // SQL语句参数，防注入
	OrderBy  string // order by
	orderByT string // order by inner temp
	GroupBy  string // group by
	groupByT string // group by inner temp
	PageSize uint32 // 分页大小
	Page     uint32 // 当前页
	Offset   uint32 // 查询偏移量
	Limit    uint32 // 查询限量

	// Gson结果相关
	GsonStr  string // 可以指定同时返回 GsonRows 数据
	GsonNeed bool   // 是否需要返回 GsonStr 数据，此时GsonVal的值可能是字符串，也可能是GsonRows对象？
	GsonOnly bool   // 只需要 GsonStr 数据，不用解析到 List

	// 缓存控制和其它标记字段
	isReady      bool   // 是否已经初始化
	CacheType    uint8  // 缓存类型
	CacheExpireS uint32 // 缓存过期时间秒（设置值大于0，意味着需要缓存，否则不需要缓存）
	cacheKey     string // 缓存的 key value（hash 值）
}
