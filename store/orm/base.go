package orm

import (
	"database/sql"
	"github.com/qinchende/gofast/store/dts"
)

const (
	dbDefAutoIncKeyName = "ID"            // 默认数据库自增键字段名
	dbDefPrimaryKeyName = "ID"            // 默认主键的字段名
	dbDefUpdatedKeyName = "UpdatedAt"     // 默认主键的字段名
	dbConfigTag         = "dbc"           // 数据库字段配置tag头
	dbAutoIncKeyFlag    = "auto_field"    // 数据库自增字段tag标记
	dbPrimaryKeyFlag    = "primary_field" // 数据库主键tag标记
	dbUpdatedKeyFlag    = "updated_field" // 更新时间
)

var (
	kvTableSchema = new(TableSchema) // 虚拟 cst.KV 类型 对应的 TableSchema 变量
)

type OrmStruct interface {
	GfAttrs(parent OrmStruct) *TableAttrs
	BeforeSave()
	AfterInsert(sql.Result)
}

// 表结构体Schema, 限制表最多127列（用int8计数）
type TableSchema struct {
	SS dts.StructSchema

	tAttrs TableAttrs // 实体类型的相关控制属性

	autoIndex    int8 // 自增字段原始索引位置
	primaryIndex int8 // 主键字段原始索引位置
	updatedIndex int8 // 更新字段原始索引位置，没有则为-1

	insertSQL string // 全字段insert（将来会建立通用缓存中心，这里暂时这样用）
	updateSQL string // 全字段update
	deleteSQL string // delete
	selectSQL string // select
}

// GoFast ORM Model need some attributes.
type TableAttrs struct {
	TableName string // 数据库表名称
	CacheAll  bool   // 是否缓存所有记录
	ExpireS   uint32 // 过期时间（秒）默认7天

	// 内部状态标记
	columnsHash uint64 // 本结构体的哈希值
	cacheKeyFmt string // 行记录缓存的Key前缀
}
