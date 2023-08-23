package config

import (
	"fmt"
	"gorm.io/driver/postgres"
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
		/*	DSN := os.Getenv("MYSQL_USER") + ":@" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_ADDRESS") +
			":" + os.Getenv("MYSQL_PORT") + ")/" +
			os.Getenv("MYSQL_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"*/

		DSN := "host=" + os.Getenv("POSTGRES_SERVER") + " user=" + os.Getenv("POSTGRES_USER") +
			" password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DATABASE") +
			" port=" + os.Getenv("POSTGRES_PORT") + " sslmode=disable TimeZone=" + os.Getenv("APP_TIMEZONE")

		log.Println(DSN)
		DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})

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
