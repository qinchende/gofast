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
		ConnStr string `cnf:",NA"`
		MaxOpen int    `cnf:",def=100,range=[1:1000]"`
		MaxIdle int    `cnf:",def=100,NA"`
	}
)

func OpenMysql(cf *ConnConfig) *sqlx.MysqlORM {
	mysqlX := sqlx.MysqlORM{Ctx: context.Background()}

	db, err := sql.Open("mysql", cf.ConnStr)
	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(cf.MaxOpen)
	db.SetMaxIdleConns(cf.MaxIdle)

	//mysqlX.Client = db
	mysqlX.Writer = db
	mysqlX.Reader = db
	return &mysqlX
}
