package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
)

var db *gorm.DB
var ProdDBConfig *gorm.Config
var once sync.Once

func init() {
	once.Do(func() {
		ProdDBConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	})
}

func NewConnection(config *gorm.Config) *gorm.DB {
	if db != nil {
		return db
	}
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), config)
	if err != nil {
		log.Fatalln("could not create new database connection", err)
	}
	db = db.Set("gorm:table_options", " ENGINE=InnoDB")

	return db
}

