package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	err error
	DB  *gorm.DB
)

func InitDB() *gorm.DB {
	if os.Getenv("DB_DRIVER") == "mysql" {
		DSN := os.Getenv("DB_USERNAME") + ":@webscrapper_db##@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
		DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})

		if err != nil {
			panic("connectionString error")
		}
		return DB
	}

	fmt.Println("DB_DRIVER not supported")
	return nil
}

func GetDBInstance() *gorm.DB {
	return DB
}
