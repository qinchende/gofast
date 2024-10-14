package clickh

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

var serverAddr = []string{"YourIP:9000"}

// Stand Model
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func connStandSqlModel() *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: serverAddr,
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		//TLS: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		//Debug:                true,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{{Name: "my-stand-model-demo", Version: "0.1"}},
		},
	})
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	return conn

	// or DSN URL
	// ++++++++++++++++++++++++++++++++++++++++++++
	//dsnUrl := "clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60"
	//conn, err := sql.Open("clickhouse", dsnUrl)
	//if err != nil {
	//	panic(err)
	//}
	//return conn
}

func standModelDemo() {
	db := connStandSqlModel()
	defer db.Close()

	// 创建表格
	createTable := "CREATE TABLE test (id Int32, value String) ENGINE = Memory"
	if _, err := db.Exec(createTable); err != nil {
		panic(err)
	}

	// 插入数据
	insertData := "INSERT INTO test (id, value) VALUES(1, 'Hello'), (2, 'World')"
	if _, err := db.Exec(insertData); err != nil {
		panic(err)
	}

	// 查询数据
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id int
	var value string
	for rows.Next() {
		if err = rows.Scan(&id, &value); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Value: %s\n", id, value)
	}

	// 清理表
	if _, err = db.Exec("DROP TABLE test"); err != nil {
		panic(err)
	}
}

// Native Model
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func connNativeModel() (driver.Conn, error) {
	var ctx = context.Background()
	var conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: serverAddr,
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{{Name: "my-native-model-demo", Version: "0.1"}},
		},
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		//TLS: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
	})
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx); err != nil {
		if exc, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exc.Code, exc.Message, exc.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func nativeModelDemo() {
	db, err := connNativeModel()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()

	// 创建表格
	createTable := "CREATE TABLE test (id Int32, value String) ENGINE = Memory"
	if err = db.Exec(ctx, createTable); err != nil {
		panic(err)
	}

	// 插入数据
	insertData := "INSERT INTO test (id, value) VALUES(1, 'Hello'), (2, 'World')"
	if err = db.Exec(ctx, insertData); err != nil {
		panic(err)
	}

	// 查询数据
	rows, err := db.Query(ctx, "SELECT * FROM test")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id int32 // must be int32
	var value string
	for rows.Next() {
		if err = rows.Scan(&id, &value); err != nil {
			panic(err)
		}
		fmt.Printf("ID: %d, Value: %s\n", id, value)
	}

	// 清理表
	if err = db.Exec(ctx, "DROP TABLE test"); err != nil {
		panic(err)
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func RunExample() {
	//standModelDemo()
	nativeModelDemo()
}
