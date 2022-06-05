package user

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

type UserTestSuite struct {
	suite.Suite
	DB              *gorm.DB
	protocol        string
	host            string
	port            string
	indexRoute      string
	apiBaseRoute    string
	contentTypeJson string
	userRoute       string
	loggedInToken   string

	repository *repositories.Repository
	actionUser *models.User
}

// SetupSuite runs before suite
func (suite *UserTestSuite) SetupSuite() {
	config.LoadDBConfig()
	database.InitialiseDatabase()

	suite.protocol = "http"
	suite.host = os.Getenv("HOST")
	suite.port = os.Getenv("PORT")
	suite.indexRoute = "/"
	suite.apiBaseRoute = "/api"
	suite.contentTypeJson = "application/json"
	suite.userRoute = fmt.Sprintf("%s://%s:%s%s/users", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
}

// SetupTest runs before each test.
func (suite *UserTestSuite) SetupTest() {
	suite.DB = database.GetDBHandle().WithContext(context.TODO())
	migrations.ForceTruncateAllTables(suite.DB)

	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (1, uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmith26@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.repository = repositories.NewRepository(suite.DB)
	suite.actionUser = suite.repository.FindUserById(1)
	assert.NotNil(suite.T(), suite.actionUser)
	// New token per test.
	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), "", token)
	suite.loggedInToken = token
}

func (suite *UserTestSuite) TearDownTest() {
	migrations.ForceTruncateAllTables(suite.DB)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
