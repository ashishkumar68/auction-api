package repositories

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/stretchr/testify/assert"
)

// go: github.com/brianvoe/gofakeit@v3.19.1+incompatible: invalid version: module contains a go.mod file, so module path must match major version ("github.com/brianvoe/gofakeit/v3")
func (suite *RepositoryTestSuite) TestBidRepository_Save() {
	item := suite.repository.FindItemById(1)
	user := suite.repository.FindUserById(1)
	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), item)
	existingBid := suite.repository.FindBidByItem(item, user)

	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), item)
	assert.Nil(suite.T(), existingBid)

	newBid := models.NewBidFromValues(item, 10, user)
	err := suite.repository.SaveBid(newBid)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.repository.FindBidByItem(item, user))
}
