package actions

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
)

func (suite *IndexTestSuite) TestIndexRoute() {
	indexRoute := "/"

	log.Println(fmt.Sprintf("%s://%s:%s%s", suite.protocol, suite.host, suite.port, indexRoute))
	resp, err := http.Get(fmt.Sprintf("%s://%s:%s%s", suite.protocol, suite.host, suite.port, indexRoute))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read response from HTTP message")
	assert.Equal(suite.T(), "{}", string(body))
}

func (suite *IndexTestSuite) TestCurrentVersionRoute() {
	indexRoute := "/version"

	log.Println(fmt.Sprintf("%s://%s:%s%s", suite.protocol, suite.host, suite.port, indexRoute))
	resp, err := http.Get(fmt.Sprintf("%s://%s:%s%s", suite.protocol, suite.host, suite.port, indexRoute))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read response from HTTP message")
	assert.Equal(suite.T(), fmt.Sprintf(`{"version":"%s"}`, os.Getenv("VERSION")), string(body))
}
