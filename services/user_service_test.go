package services

import (
	"context"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/stretchr/testify/assert"
)

func (suite *ServiceTestSuite) TestPlaceBidOnItem() {
	item := suite.repository.FindItemById(1)
	user := suite.repository.FindUserById(1)
	assert.NotNil(suite.T(), item)
	assert.NotNil(suite.T(), user)

	placeBidForm := forms.PlaceNewItemBidForm{
		ItemId:        item.ID,
		BidUserId:     user.ID,
		BidValue:      12,
		AuditableForm: forms.AuditableForm{ActionUser: suite.actionUser},
	}
	//itemService := NewItemService(suite.DB)
	newBid, err := suite.itemService.PlaceItemBid(context.TODO(), placeBidForm)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), newBid)

	existingBid := suite.repository.FindBidByItem(item, user)
	assert.NotNil(suite.T(), existingBid)
	assert.True(suite.T(), existingBid.ID == newBid.ID)
	assert.True(suite.T(), existingBid.Value == newBid.Value)
}
