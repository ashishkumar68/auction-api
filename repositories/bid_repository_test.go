package repositories

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/stretchr/testify/assert"
)

func (suite *RepositoryTestSuite) TestBidRepository_Save() {
	user := suite.userRepo.Find(1)
	item := suite.itemRepo.Find(1)
	existingBid := suite.bidRepo.FindByItem(item, user)

	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), item)
	assert.True(suite.T(), existingBid.IsZero())

	newBid := models.NewBidFromValues(item, 10, user)
	err := suite.bidRepo.Save(newBid)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.bidRepo.FindByItem(item, user))
}