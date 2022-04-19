package services

import (
	"context"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/stretchr/testify/assert"
	"time"
)

func (suite *ServiceTestSuite) TestPlaceBidOnItem() {
	item := suite.repository.FindItemById(1)
	user := suite.repository.FindUserById(2)
	assert.NotNil(suite.T(), item)
	assert.NotNil(suite.T(), user)

	placeBidForm := forms.PlaceNewItemBidForm{
		ItemId:        item.ID,
		BidUserId:     user.ID,
		BidValue:      12,
		AuditableForm: forms.AuditableForm{ActionUser: suite.actionUser},
	}
	newBid, err := suite.itemService.PlaceItemBid(context.TODO(), placeBidForm)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), newBid)

	existingBid := suite.repository.FindBidByItem(item, user)
	assert.NotNil(suite.T(), existingBid)
	assert.True(suite.T(), existingBid.ID == newBid.ID)
	assert.True(suite.T(), existingBid.Value == newBid.Value)
}

func (suite *ServiceTestSuite) TestItDoesntPlaceBidOnItemAfterLastBidDate() {
	suite.DB.Exec(`
UPDATE items SET last_bid_date = "2022-02-01" WHERE id = 2;
`)
	item := suite.repository.FindItemById(2)
	assert.NotNil(suite.T(), item)

	assert.True(suite.T(), item.LastBidDate.Format(time.RFC3339) == "2022-02-01T00:00:00Z")
	assert.True(suite.T(), time.Now().After(item.LastBidDate))

	placeBidForm := forms.PlaceNewItemBidForm{
		ItemId:        item.ID,
		BidUserId:     suite.actionUser.ID,
		BidValue:      14,
		AuditableForm: forms.AuditableForm{ActionUser: suite.actionUser},
	}
	//itemService := NewItemService(suite.DB)
	newBid, err := suite.itemService.PlaceItemBid(context.TODO(), placeBidForm)
	assert.Nil(suite.T(), newBid)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), ItemNotBidEligible, err)
}
