package repositories

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/stretchr/testify/assert"
)

func (suite *RepositoryTestSuite) TestBidRepository_Save() {
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	existingBid := suite.repository.FindBidByItem(item, suite.actionUser)

	assert.Nil(suite.T(), existingBid)

	newBid := models.NewBidFromValues(item, 10, suite.actionUser)
	err := suite.repository.SaveBid(newBid)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.repository.FindBidByItem(item, suite.actionUser))
}
