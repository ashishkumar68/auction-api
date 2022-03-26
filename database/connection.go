package database

import (
	"github.com/ashishkumar68/auction-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func NewConnection() *gorm.DB {
	if db != nil {
		return db
	}
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), config.DBConfig)
	if err != nil {
		log.Fatalln("could not create new database connection", err)
	}
	db = db.Set("gorm:table_options", " ENGINE=InnoDB")

	return db
}

