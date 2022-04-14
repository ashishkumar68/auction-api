package database

import (
	"github.com/ashishkumar68/auction-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var dbObj *gorm.DB

func InitialiseDatabase() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), config.DBConfig)
	if err != nil {
		log.Fatalln("could not create new database connection", err)
	}
	dbObj = db.Set("gorm:table_options", " ENGINE=InnoDB")
}

func GetDBHandle() *gorm.DB {
	return dbObj
}

//func NewConnectionWithContext(ctx context.Context) *gorm.DB {
//	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), config.DBConfig)
//	if err != nil {
//		log.Fatalln("could not create new database connection", err)
//	}
//	db = db.Set("gorm:table_options", " ENGINE=InnoDB").WithContext(ctx)
//
//	return db
//}
