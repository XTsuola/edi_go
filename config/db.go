package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB // 全局数据库对象

func InitDB() {
	// 配置
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "123456"
	dbname := "edi-mms"

	// PostgreSQL DSN 正确写法
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	fmt.Println("连接信息：", dsn)

	// ✅ 正确连接 PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 测试连接
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("获取数据库实例失败：", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Ping 失败：", err)
	}

	fmt.Println("✅ 成功连接 PostgreSQL 数据库：edi")
}
