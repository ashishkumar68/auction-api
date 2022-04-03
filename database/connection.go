package database

import (
	"context"
	"github.com/ashishkumar68/auction-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewConnectionWithContext(ctx context.Context) *gorm.DB {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB_DSN")), config.DBConfig)
	if err != nil {
		log.Fatalln("could not create new database connection", err)
	}
	db = db.Set("gorm:table_options", " ENGINE=InnoDB").WithContext(ctx)

	return db
}

