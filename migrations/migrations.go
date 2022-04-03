package migrations

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/models"
	"log"
)

func DropAndCreateTables() {
	db := database.NewConnectionWithContext(context.TODO())

	db.Exec("SET foreign_key_checks = 0;")
	if err := db.Migrator().DropTable(
		&models.User{},
		&models.Item{}); err != nil {
		log.Fatalln(fmt.Sprintf("could not drop tables due to:"), err)
	}
	if err := db.Migrator().CreateTable(
		&models.User{},
		&models.Item{}); err != nil {
		log.Fatalln(fmt.Sprintf("could not create tables due to:"), err)
	}
	db.Exec("SET foreign_key_checks = 1;")

	log.Println("Migrations have been run successfully.")
}
