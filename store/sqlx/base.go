// Copyright 2023 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sqlx

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/store/dts"
	"github.com/qinchende/gofast/store/orm"
	"reflect"
	"sync"
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
	TypeSqlRawBytes reflect.Type
	sharedAnyValue  = new(dts.SqlSkip)
)

func init() {
	TypeSqlRawBytes = reflect.TypeOf(sql.RawBytes{})
}

// 不同数据库，定义不同的SQL语句执行器
type CmdBuilder interface {
	// 增删改
	Insert(*orm.TableSchema) string
	Delete(*orm.TableSchema) string
	Update(*orm.TableSchema) string
	UpdateColumns(ts *orm.TableSchema, rVal *reflect.Value, columns ...string) (string, []any)

	// 查询
	SelectPrimary(ts *orm.TableSchema) string
	SelectRow(ts *orm.TableSchema, columns string, where string) string
	SelectRows(ts *orm.TableSchema, columns string, where string) string
	SelectByPet(*SelectPet) string
	SelectCountByPet(*SelectPet) string
	SelectPagingByPet(*SelectPet) string

	// 初始化Pet
	InitPet(*SelectPet, *orm.TableSchema)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 天然支持读写分离，只需要数据库连接配置文件，分别传入读写库的连接地址
type OrmDB struct {
	Cmd      CmdBuilder       // 数据库执行语句生成
	Attrs    *DBAttrs         // 数据库属性
	Ctx      context.Context  // 上下文，方便取消执行
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
	ArgsMem  [5]any // 最多5个变量的内存空间，用此参数可以避免一次切片内存分配
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
	GsonStr  string // 已经是序列化之后的字符串
	GsonNeed bool   // 是否需要返回 GsonStr 数据，此时GsonVal的值可能是字符串，也可能是GsonRows对象？
	GsonOnly bool   // 只需要 GsonStr 数据，不用解析到 List

	// 缓存控制和其它标记字段
	isReady      bool   // 是否已经初始化
	CacheType    uint8  // 缓存类型
	CacheExpireS uint32 // 缓存过期时间秒（设置值大于0，意味着需要缓存，否则不需要缓存）
	cacheKey     string // 缓存的 key value（hash 值）
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// SelectPet pool
var (
	_SelectPetDefValue SelectPet
	selectPetPool      = sync.Pool{New: func() any { return &SelectPet{} }}
)

// 取缓存对象，并初始化成默认值
func GetSelectPet() *SelectPet {
	pet := selectPetPool.Get().(*SelectPet)
	*pet = _SelectPetDefValue
	return pet
}

// 直接取缓存对象并返回，该对象数据没有被初始化，使用需谨慎
func GetSelectPetRaw() *SelectPet {
	return selectPetPool.Get().(*SelectPet)
}

// 回收对象
func PutSelectPet(pet *SelectPet) {
	selectPetPool.Put(pet)
}
