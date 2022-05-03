package item

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type ItemTestSuite struct {
	suite.Suite
	DB                   *gorm.DB
	protocol             string
	host                 string
	port                 string
	indexRoute           string
	apiBaseRoute         string
	itemsRoute           string
	contentTypeJson      string
	contentTypeMultipart string
	loggedInToken        string
	baseFSItemsPath      string
	itemImageFile1       *os.File
	itemImageFile2       *os.File

	repository *repositories.Repository
	actionUser *models.User
}

// SetupSuite runs before suite
func (suite *ItemTestSuite) SetupSuite() {
	config.LoadDBConfig()
	database.InitialiseDatabase()

	suite.protocol = "http"
	suite.host = os.Getenv("HOST")
	suite.port = os.Getenv("PORT")
	suite.indexRoute = "/"
	suite.apiBaseRoute = "/api"
	suite.contentTypeJson = "application/json"
	suite.contentTypeMultipart = "multipart/form-data"
	suite.baseFSItemsPath = fmt.Sprintf("%s/%s", os.Getenv("FILE_UPLOADS_DIR"), models.BaseFSItemsPrefix)
	suite.itemsRoute = fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)

	itemImageFile1, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_1.jpeg", os.Getenv("PROJECTDIR")))
	assert.Nilf(suite.T(), err, "could not load test file")
	suite.itemImageFile1 = itemImageFile1
	itemImageFile2, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_2.jpg", os.Getenv("PROJECTDIR")))
	assert.Nilf(suite.T(), err, "could not load test file")
	suite.itemImageFile2 = itemImageFile2
}

// SetupTest runs before each test.
func (suite *ItemTestSuite) SetupTest() {
	suite.DB = database.GetDBHandle().WithContext(context.TODO())
	migrations.ForceTruncateAllTables(suite.DB)

	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (5, uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmith24@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.repository = repositories.NewRepository(suite.DB)
	suite.actionUser = suite.repository.FindUserById(5)
	assert.NotNil(suite.T(), suite.actionUser)

	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), "", token)
	suite.loggedInToken = token
	// clean up items file system.
	err = os.RemoveAll(suite.baseFSItemsPath)
	assert.Nil(suite.T(), err, fmt.Sprintf("could not clear items base path: %s", suite.baseFSItemsPath))
}

func (suite *ItemTestSuite) TearDownTest() {
	migrations.ForceTruncateAllTables(suite.DB)
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}
