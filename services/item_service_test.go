package services

import (
	"context"
	"github.com/ashishkumar68/auction-api/forms"
	"github.com/stretchr/testify/assert"
	"time"
)

func (suite *ServiceTestSuite) TestPlaceBidOnItem() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	bidUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), bidUser)
	assert.False(suite.T(), item.IsOwner(*bidUser))

	placeBidForm := forms.PlaceNewItemBidForm{
		ItemId:        item.ID,
		BidUserId:     bidUser.ID,
		BidValue:      12,
		AuditableForm: forms.AuditableForm{ActionUser: bidUser},
	}
	newBid, err := suite.itemService.PlaceItemBid(context.TODO(), placeBidForm)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), newBid)

	existingBid := suite.repository.FindBidByItem(item, bidUser)
	assert.NotNil(suite.T(), existingBid)
	assert.True(suite.T(), existingBid.ID == newBid.ID)
	assert.True(suite.T(), existingBid.Value == newBid.Value)
}

func (suite *ServiceTestSuite) TestItDoesntPlaceBidOnItemAfterLastBidDate() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
UPDATE items SET last_bid_date = "2022-02-01" WHERE id = 2;
`)
	item := suite.repository.FindItemById(2)
	assert.NotNil(suite.T(), item)
	bidUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), bidUser)
	assert.False(suite.T(), item.IsOwner(*bidUser))

	assert.True(suite.T(), item.LastBidDate.Format(time.RFC3339) == "2022-02-01T00:00:00Z")
	assert.True(suite.T(), time.Now().After(item.LastBidDate))

	placeBidForm := forms.PlaceNewItemBidForm{
		ItemId:        item.ID,
		BidUserId:     bidUser.ID,
		BidValue:      14,
		AuditableForm: forms.AuditableForm{ActionUser: bidUser},
	}
	newBid, err := suite.itemService.PlaceItemBid(context.TODO(), placeBidForm)
	assert.Nil(suite.T(), newBid)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), ItemNotBidEligible, err)
}
