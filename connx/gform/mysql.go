package gform

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/store/sqlx"
	"log"
	"strings"
	"time"
)

type (
	ConnCnf struct {
		ConnStr    string   `v:"required"`
		ConnStrR   string   `v:"required=false"`
		MaxOpen    int      `v:"def=100,range=[1:1000]"`
		MaxIdle    int      `v:"def=100"`
		RedisNodes []string `v:"required=false,len=[10:300]"`
	}
)

func OpenMysql(cf *ConnCnf) *sqlx.MysqlORM {
	mysqlOrm := sqlx.MysqlORM{Attrs: &sqlx.MysqlORMAttrs{}, Ctx: context.Background()}

	// DBName ->
	dbAttrs, _ := mysql.ParseDSN(cf.ConnStr)
	if dbAttrs != nil {
		// 必须统一数据库名称，全部转换成小写
		// 将来表缓存的时候需要用到这里的DBName
		mysqlOrm.Attrs.DBName = strings.ToLower(dbAttrs.DBName)
	}

	// 主库连接
	writer, err := sql.Open("mysql", cf.ConnStr)
	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}
	// See "Important settings" section.
	writer.SetConnMaxLifetime(time.Minute * 3)
	writer.SetMaxOpenConns(cf.MaxOpen)
	writer.SetMaxIdleConns(cf.MaxIdle)
	mysqlOrm.Writer = writer

	// 从库连接
	// 如果配置文件配置了只读数据库，应用于读写分离
	if cf.ConnStrR != "" {
		reader, err := sql.Open("mysql", cf.ConnStrR)
		if err != nil {
			log.Fatalf("Conn %s err: %s", cf.ConnStrR, err)
		}
		// See "Important settings" section.
		reader.SetConnMaxLifetime(time.Minute * 3)
		reader.SetMaxOpenConns(cf.MaxOpen)
		reader.SetMaxIdleConns(cf.MaxIdle)
		mysqlOrm.Reader = reader
	} else {
		mysqlOrm.Reader = mysqlOrm.Writer
	}

	// redis cache
	rds := cf.RedisNodes
	rdsNodes := make([]gfrds.GfRedis, len(rds))
	for i := 0; i < len(rds); i++ {
		rdsCnf := gfrds.ParseDsn(rds[i])
		rdsNodes[i] = *gfrds.NewGoRedis(rdsCnf)
	}
	mysqlOrm.SetRdsNodes(&rdsNodes)

	return &mysqlOrm
}
