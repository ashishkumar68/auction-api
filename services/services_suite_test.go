package services

import (
	"context"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	DB *gorm.DB

	repository  *repositories.Repository
	userService UserService
	itemService ItemService
	actionUser  *models.User
}

// SetupSuite runs before suite
func (suite *ServiceTestSuite) SetupSuite() {
	config.LoadDBConfig()
	database.InitialiseDatabase()
}

// SetupTest runs before each test.
func (suite *ServiceTestSuite) SetupTest() {

	suite.DB = database.GetDBHandle().WithContext(context.TODO())
	suite.userService = NewUserService(suite.DB)
	suite.itemService = NewItemService(suite.DB)

	migrations.ForceTruncateAllTables(suite.DB)
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (10, uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmith25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES 
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,10,10,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,10,10,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-02")
;
`)
	suite.repository = repositories.NewRepository(suite.DB)
	suite.actionUser = suite.repository.FindUserById(10)
	assert.NotNil(suite.T(), suite.actionUser)
}

func (suite *ServiceTestSuite) TearDownTest() {
	migrations.ForceTruncateAllTables(suite.DB)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
