package actions

import (
	"context"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type IndexTestSuite struct {
	suite.Suite
	DB       *gorm.DB
	protocol string
	host     string
	port     string

	repository *repositories.Repository
}

// SetupSuite runs before suite
func (suite *IndexTestSuite) SetupSuite() {
	config.LoadDBConfig()
	database.InitialiseDatabase()

	suite.protocol = "http"
	suite.host = os.Getenv("HOST")
	suite.port = os.Getenv("PORT")
}

// SetupTest runs before each test.
func (suite *IndexTestSuite) SetupTest() {
	suite.DB = database.GetDBHandle().WithContext(context.TODO())
	migrations.ForceTruncateAllTables(suite.DB)
}

func (suite *IndexTestSuite) TearDownTest() {
	migrations.ForceTruncateAllTables(suite.DB)
}

func TestIndexTestSuite(t *testing.T) {
	suite.Run(t, new(IndexTestSuite))
}
