package mysql

import (
	"context"
	"database/sql"
	"github.com/qinchende/gofast/store/sqlx"
	"log"
	"time"
)

type (
	ConnConfig struct {
		ConnStr  string `cnf:",NA"`
		ConnStrR string `cnf:",NA"`
		MaxOpen  int    `cnf:",def=100,range=[1:1000]"`
		MaxIdle  int    `cnf:",def=100,NA"`
	}
)

func OpenMysql(cf *ConnConfig) *sqlx.MysqlORM {
	mysqlX := sqlx.MysqlORM{Ctx: context.Background()}

	writer, err := sql.Open("mysql", cf.ConnStr)
	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}
	// See "Important settings" section.
	writer.SetConnMaxLifetime(time.Minute * 3)
	writer.SetMaxOpenConns(cf.MaxOpen)
	writer.SetMaxIdleConns(cf.MaxIdle)
	mysqlX.Writer = writer

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
		mysqlX.Reader = reader
	} else {
		mysqlX.Reader = mysqlX.Writer
	}

	return &mysqlX
}
