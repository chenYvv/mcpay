package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// DB 对象
var DB *gorm.DB

// MysqlDB 数据库
func MysqlDB(dsn string, maxOpenConnections int, maxIdleConnections int, maxLifeSeconds int) (*gorm.DB, error) {

	//var dbConfig gorm.Dialector
	//
	//dbConfig = mysql.New(mysql.Config{
	//	DSN: dsn,
	//})

	// 使用 gorm.Open 连接数据库
	var err error
	instance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 处理错误
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	//instance = instance.Debug()

	// 获取底层的 sqlDB
	SqlDB, err := instance.DB()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	if maxOpenConnections > 0 {
		// 设置最大连接数
		SqlDB.SetMaxOpenConns(maxOpenConnections)
		// 设置最大空闲连接数
		SqlDB.SetMaxIdleConns(maxIdleConnections)
		// 设置每个链接的过期时间
		SqlDB.SetConnMaxLifetime(time.Duration(maxLifeSeconds) * time.Second)
	}

	return instance, nil
}
