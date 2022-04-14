package gform

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/connx/gfrds"
	"github.com/qinchende/gofast/store/sqlx"
	"log"
	"time"
)

type (
	ConnCnf struct {
		ConnStr      string   `cnf:",NA"`
		ConnStrR     string   `cnf:",NA"`
		MaxOpen      int      `cnf:",def=100,range=[1:1000]"`
		MaxIdle      int      `cnf:",def=100,NA"`
		RedisCluster []string `cnf:",NA"`
	}
)

func OpenMysql(cf *ConnCnf) *sqlx.MysqlORM {
	mysqlOrm := sqlx.MysqlORM{Ctx: context.Background()}

	writer, err := sql.Open("mysql", cf.ConnStr)
	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}
	// See "Important settings" section.
	writer.SetConnMaxLifetime(time.Minute * 3)
	writer.SetMaxOpenConns(cf.MaxOpen)
	writer.SetMaxIdleConns(cf.MaxIdle)
	mysqlOrm.Writer = writer

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
	rds := cf.RedisCluster
	rdsNodes := make([]gfrds.GfRedis, len(rds))
	for i := 0; i < len(rds); i++ {
		rdsCnf := gfrds.ParseDsn(rds[i])
		rdsNodes[i] = *gfrds.NewGoRedis(rdsCnf)
	}
	mysqlOrm.SetRdsNodes(&rdsNodes)

	return &mysqlOrm
}
