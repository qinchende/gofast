package gormc

import (
	mysqlGF "github.com/qinchende/gofast/connx/mysql"
	mysqlGorm "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type (
	OrmX gorm.DB
)

func NewGormConn(cf *mysqlGF.ConnConfig) *gorm.DB {
	ormX, err := gorm.Open(mysqlGorm.New(mysqlGorm.Config{
		DSN:                       cf.ConnStr,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Conn %s err: %s", cf.ConnStr, err)
	}
	return ormX
}
