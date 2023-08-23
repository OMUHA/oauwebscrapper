package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	err error
	DB  *gorm.DB
)

func InitDB() *gorm.DB {
	if os.Getenv("DB_DRIVER") == "mysql" {
		DSN := os.Getenv("MYSQL_USER") + ":@webscrapper_db@tcp(" + os.Getenv("MYSQL_ADDRESS") +
			":" + os.Getenv("MYSQL_PORT") + ")/" +
			os.Getenv("MYSQL_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
		log.Println(DSN)
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
