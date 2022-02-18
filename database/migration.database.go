package database

import (
	"github.com/OMUHA/oauwebscrapper/app/models"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Book{})
}
