package repositories

import (
	"context"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type RepositoryTestSuite struct {
	suite.Suite
	DB *gorm.DB

	repository *Repository
}

// SetupSuite runs before suite
func (suite *RepositoryTestSuite) SetupSuite() {
	config.LoadDBConfig()
	database.InitialiseDatabase()
}

// SetupTest runs before each test.
func (suite *RepositoryTestSuite) SetupTest() {
	suite.DB = database.GetDBHandle().WithContext(context.TODO())
	migrations.ForceTruncateAllTables(suite.DB)

	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (1, uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmith24@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES 
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,1,1,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01")
;
`)

	suite.repository = NewRepository(suite.DB)

	user := suite.repository.FindUserById(1)
	suite.DB.Set("actionUser", user)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	migrations.ForceTruncateAllTables(suite.DB)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
