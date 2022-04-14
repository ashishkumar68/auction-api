package migrations

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/models"
	"log"
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

	if err := db.Migrator().DropTable(
		&models.User{},
		&models.Item{},
		&models.Bid{}); err != nil {
		log.Fatalln(fmt.Sprintf("could not drop tables due to:"), err)
	}
	if err := db.Migrator().CreateTable(
		&models.User{},
		&models.Item{},
		&models.Bid{}); err != nil {
		log.Fatalln(fmt.Sprintf("could not create tables due to:"), err)
	}
	db.Exec("SET foreign_key_checks = 1;")

	log.Println("Migrations have been run successfully.")
}
