package clickh

import (
	"context"
	"database/sql"
	"log"
)

type (
	ConnConfig struct {
		ConnStr string `cnf:",NA"`
	}
	ClickHouseX struct {
		Cli *sql.DB
		Ctx context.Context
	}
)

func NewClickH(cf *ConnConfig) *ClickHouseX {
	chX := ClickHouseX{Ctx: context.Background()}

	conn, err := sql.Open("clickhouse", cf.ConnStr)
	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}

	chX.Cli = conn
	return &chX
}
