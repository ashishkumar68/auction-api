package migrations

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations/yr2022"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func DropAndCreateTables() {
	database.InitialiseDatabase()
	db := database.GetDBHandle()

	db.Exec("SET foreign_key_checks = 0;")
	db.Exec("DROP FUNCTION IF EXISTS uuid_v4;")
	db.Exec(`
CREATE FUNCTION uuid_v4()
    RETURNS CHAR(36)
BEGIN
    -- 1th and 2nd hex blocks are made of 6 random bytes
    SET @h1 = HEX(RANDOM_BYTES(4));
    SET @h2 = HEX(RANDOM_BYTES(2));

    -- 3th block will start with a 4 indicating the version, remaining is random
    SET @h3 = SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3);

    -- 4th block first nibble can only be 8, 9 A or B, remaining is random
    SET @h4 = CONCAT(HEX(FLOOR(ASCII(RANDOM_BYTES(1)) / 64)+8),
                SUBSTR(HEX(RANDOM_BYTES(2)), 2, 3));

    -- 5th block is made of 6 random bytes
    SET @h5 = HEX(RANDOM_BYTES(6));

    -- Build the complete UUID
    RETURN LOWER(CONCAT(
        @h1, '-', @h2, '-4', @h3, '-', @h4, '-', @h5
    ));
END
`)

	migrateTypes := []any{
		&models.User{}, &models.Item{}, &models.Bid{}, &models.Reaction{}, &models.ItemComment{}, &models.ItemImage{},
	}

	if err := db.Migrator().DropTable(migrateTypes...); err != nil {
		log.Fatalln(fmt.Sprintf("could not drop tables due to:"), err)
	}
	if err := db.Migrator().CreateTable(migrateTypes...); err != nil {
		log.Fatalln(fmt.Sprintf("could not create tables due to:"), err)
	}
	db.Exec("SET foreign_key_checks = 1;")

	log.Println("Migrations have been run successfully.")
}

func ForceTruncateAllTables(db *gorm.DB) {
	if config.AppEnvTest != os.Getenv("APP_ENV") {
		log.Fatalln("Force truncate is only allowed in test environment.")
	}
	for retry := 1; retry <= 3; retry++ {
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := db.Exec(`SET foreign_key_checks = 0;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`TRUNCATE TABLE users;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`TRUNCATE TABLE items;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`TRUNCATE TABLE bids;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`TRUNCATE TABLE reactions;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`TRUNCATE TABLE item_comments;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`TRUNCATE TABLE item_images;`).Error; err != nil {
				return err
			}
			if err := db.Exec(`SET foreign_key_checks = 1;`).Error; err != nil {
				return err
			}
			return nil
		})

		if txErr == nil {
			break
		}
	}
}

func RunMigrations() {
	database.InitialiseDatabase()
	db := database.GetDBHandle().WithContext(context.TODO())
	db.Logger = logger.Default.LogMode(logger.Info)

	m := gormigrate.New(db, gormigrate.DefaultOptions, getAllMigrationsList())
	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migrations were executed successfully.")
}

func getAllMigrationsList() []*gormigrate.Migration {
	migrations := make([]*gormigrate.Migration, 0)
	migrations = append(migrations, yr2022.GetMigrationsList()...)

	return migrations
}
