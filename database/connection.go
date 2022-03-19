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
var writeMutex sync.Mutex

func NewConnection() *gorm.DB {
	if db != nil {
		return db
	}
	// to allow only one goroutine
	writeMutex.Lock()
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("could not create new database connection", err)
	}
	writeMutex.Unlock()
	db = db.Set("gorm:table_options", " ENGINE=InnoDB")

	return db
}

