package orm

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/connx/redis"
	"github.com/qinchende/gofast/store/sqlx"
	"log"
	"strings"
	"time"
)

type (
	MysqlConfig struct {
		ConnStr    string   `v:"must"`                    // 主库连接
		ConnStrR   string   `v:"must=false"`              // 从库连接
		MaxOpen    int      `v:"def=100,range=[1:1000]"`  // 最大连接数
		MaxIdle    int      `v:"def=100"`                 // 最大空闲连接数
		RedisNodes []string `v:"must=false,len=[10:300]"` // 做缓存的Redis节点
	}
)

func OpenMysql(cf *MysqlConfig) *sqlx.OrmDB {
	ormDB := sqlx.OrmDB{Attrs: &sqlx.DBAttrs{DriverName: "mysql"}, Ctx: context.Background()}
	ormDB.Cmd = &sqlx.MysqlBuilder{}

	// DBName ->
	dbConfig, _ := mysql.ParseDSN(cf.ConnStr)
	if dbConfig != nil {
		// 必须统一数据库名称，全部转换成小写
		// 将来表缓存的时候需要用到这里的DBName
		ormDB.Attrs.DbName = strings.ToLower(dbConfig.DBName)
	}

	// 主库连接
	writer, err := sql.Open(ormDB.Attrs.DriverName, cf.ConnStr)
	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}
	// See "Important settings" section.
	writer.SetConnMaxLifetime(time.Minute * 3)
	writer.SetMaxOpenConns(cf.MaxOpen)
	writer.SetMaxIdleConns(cf.MaxIdle)
	ormDB.Writer = writer

	// 从库连接
	// 如果配置文件配置了只读数据库，应用于读写分离
	if cf.ConnStrR != "" {
		reader, err := sql.Open(ormDB.Attrs.DriverName, cf.ConnStrR)
		if err != nil {
			log.Fatalf("Conn %s err: %s", cf.ConnStrR, err)
		}
		// See "Important settings" section.
		reader.SetConnMaxLifetime(time.Minute * 3)
		reader.SetMaxOpenConns(cf.MaxOpen)
		reader.SetMaxIdleConns(cf.MaxIdle)
		ormDB.Reader = reader
	} else {
		ormDB.Reader = ormDB.Writer
	}

	// redis cache
	rds := cf.RedisNodes
	rdsNodes := make([]redis.GfRedis, len(rds))
	for i := 0; i < len(rds); i++ {
		rdsCnf := redis.ParseDsn(rds[i])
		rdsNodes[i] = *redis.NewGoRedis(rdsCnf)
	}
	ormDB.SetRdsNodes(&rdsNodes)

	return &ormDB
}
