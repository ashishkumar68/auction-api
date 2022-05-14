package m5

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Version0001 = &gormigrate.Migration{
	ID: "202205140001",
	Migrate: func(db *gorm.DB) error {
		err := db.Exec(`
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
`).Error
		if err != nil {
			return err
		}
		if err := db.AutoMigrate(&models.User{}); err != nil {
			return err
		}
		if err := db.AutoMigrate(&models.Item{}); err != nil {
			return err
		}
		if err := db.AutoMigrate(
			&models.Bid{}, &models.Reaction{}, &models.ItemComment{}, &models.ItemImage{},
		); err != nil {
			return err
		}

		return nil
	},
	Rollback: func(db *gorm.DB) error {
		if err := db.Migrator().DropTable(&models.ItemComment{}); err != nil {
			return err
		}
		if err := db.Migrator().DropTable(&models.ItemImage{}); err != nil {
			return err
		}
		if err := db.Migrator().DropTable(&models.Reaction{}); err != nil {
			return err
		}
		if err := db.Migrator().DropTable(&models.Bid{}); err != nil {
			return err
		}
		if err := db.Migrator().DropTable(&models.Item{}); err != nil {
			return err
		}
		if err := db.Migrator().DropTable(&models.User{}); err != nil {
			return err
		}

		return nil
	},
}
