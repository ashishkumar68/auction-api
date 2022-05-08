package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/response"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/ashishkumar68/auction-api/utils"
	"github.com/morkid/paginate"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"time"
)

func (suite *ItemTestSuite) TestDontAllowAddItemAsAnonymousUser() {
	itemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)

	items := suite.repository.FindItemByName("ABC Washing Machine")
	assert.Empty(suite.T(), items)

	payload := `
{
    "name": "ABC Washing Machine",
    "description": "A washing machine",
    "category": 1,
    "brandName": "ABC",
    "marketValue": 20000,
	"lastBidDate": "2022-10-02T00:00:00Z"
}
`
	resp, err := http.Post(itemsRoute, suite.contentTypeJson, bytes.NewReader([]byte(payload)))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *ItemTestSuite) TestAllowAddItemAsLoggedInUser() {
	itemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)

	payload := `
{
    "name": "ABC Washing Machine",
    "description": "A washing machine",
    "category": 1,
    "brandName": "ABC",
    "marketValue": 20000,
	"lastBidDate": "2022-10-02T00:00:00Z"
}
`
	resp, err := client.MakeRequest(
		itemsRoute,
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		bytes.NewReader([]byte(payload)),
	)
	defer resp.Body.Close()

	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	var item models.Item
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &item)
	assert.Nil(suite.T(), err, "Could not read response from HTTP message")
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	items := suite.repository.FindItemByName("ABC Washing Machine")
	assert.Len(suite.T(), items, 1)

	actualItem := suite.repository.FindItemByUuid(item.Uuid)
	assert.NotNil(suite.T(), actualItem)
	assert.Contains(suite.T(), actualItem.Name, "ABC Washing Machine")
	assert.NotNil(suite.T(), actualItem.UserCreated)
	assert.Equal(suite.T(), "johnsmith24@abc.com", actualItem.UserCreated.Email)
}

func (suite *ItemTestSuite) TestListItemsAnonymously() {
	itemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01")
;
`)

	items := suite.repository.FindItemByName("ABC Item")
	assert.Len(suite.T(), items, 2)

	resp, err := client.MakeRequest(
		itemsRoute,
		"GET",
		map[string]string{},
		map[string]string{},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")

	var page paginate.Page
	responseBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read from response message body.")
	err = json.Unmarshal(responseBytes, &page)
	assert.Nil(suite.T(), err, "Could not parse items list from response message.")
	pageItems, err := json.Marshal(page.Items)
	assert.Nil(suite.T(), err, "Could not serialize payload to bytes.")
	err = json.Unmarshal(pageItems, &items)
	assert.Nil(suite.T(), err, "Could not deserialize items.")
	assert.Len(suite.T(), items, 2)
	assert.Contains(suite.T(), items[0].Name, "ABC Item")
	assert.Contains(suite.T(), items[1].Name, "ABC Item")
}

func (suite *ItemTestSuite) TestListItemsAsLoggedInUser() {
	itemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01")
;
`)

	items := suite.repository.FindItemByName("ABC Item")
	assert.Len(suite.T(), items, 2)

	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	resp, err := client.MakeRequest(
		itemsRoute,
		"GET",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")

	var page paginate.Page
	responseBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read from response message body.")
	err = json.Unmarshal(responseBytes, &page)
	assert.Nil(suite.T(), err, "Could not parse items list from response message.")
	pageItems, err := json.Marshal(page.Items)
	assert.Nil(suite.T(), err, "Could not serialize payload to bytes.")
	err = json.Unmarshal(pageItems, &items)
	assert.Nil(suite.T(), err, "Could not deserialize items.")
	assert.Len(suite.T(), items, 2)
	assert.Contains(suite.T(), items[0].Name, "ABC Item")
	assert.Contains(suite.T(), items[1].Name, "ABC Item")
}

func (suite *ItemTestSuite) TestAllowItemBidsAsLoggedInUser() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	items := suite.repository.FindItemByName("ABC Item")
	assert.Len(suite.T(), items, 2)
	bidUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), bidUser)

	bidPayload := `
{
	"bidValue": 12
}
`
	token, err := services.GenerateNewJwtToken(bidUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/bid", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		bytes.NewReader([]byte(bidPayload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	item := suite.repository.FindItemById(itemId)
	newBid := suite.repository.FindBidByItem(item, bidUser)
	assert.NotNil(suite.T(), newBid)
	assert.False(suite.T(), newBid.IsZero())
	assert.Equal(suite.T(), newBid.Value, models.Value(12))
}

func (suite *ItemTestSuite) TestAllowEditItemByAuthor() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))
	assert.Equal(suite.T(), "ABC Item 1", item.Name)
	assert.Equal(suite.T(), "Item 1 Description", item.Description)
	assert.Equal(suite.T(), "ABC", item.BrandName)
	assert.Equal(suite.T(), models.ItemCategory(models.CategoryAppliancesInt), item.Category)
	assert.Equal(suite.T(), models.Value(20000), item.MarketValue)
	assert.True(suite.T(), item.LastBidDate.Equal(time.Date(2099, time.January, 1, 0, 0, 0, 0, time.UTC)))

	editItemPayload := `
{
	"name": "ABC Washing Machine UPDATED",
    "description": "A washing machine UPDATED",
    "category": 0,
    "brandName": "ABC UPDATED",
    "marketValue": 25000,
	"lastBidDate": "2025-10-02T00:00:00Z"
}
`
	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"PATCH",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		bytes.NewReader([]byte(editItemPayload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	item = suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.True(suite.T(), item.UserCreated.IsSameAs(suite.actionUser.BaseModel))
	assert.True(suite.T(), item.UserUpdated.IsSameAs(suite.actionUser.BaseModel))
	assert.Equal(suite.T(), "ABC Washing Machine UPDATED", item.Name)
	assert.Equal(suite.T(), "A washing machine UPDATED", item.Description)
	assert.Equal(suite.T(), "ABC UPDATED", item.BrandName)
	assert.Equal(suite.T(), models.ItemCategory(models.CategoryElectronicsInt), item.Category)
	assert.Equal(suite.T(), models.Value(25000), item.MarketValue)
	assert.True(suite.T(), item.LastBidDate.Equal(time.Date(2025, time.October, 2, 0, 0, 0, 0, time.UTC)))
}

func (suite *ItemTestSuite) TestDontAllowBidsForOffBidItem() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	suite.DB.Exec(`
UPDATE items SET off_bid = 1 WHERE id = 1;
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.True(suite.T(), item.IsOffBid())
	bidUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), bidUser)

	bidPayload := `
{
	"bidValue": 12
}
`
	token, err := services.GenerateNewJwtToken(bidUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/bid", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		bytes.NewReader([]byte(bidPayload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)
	var errResp response.HttpErrorResponse
	err = json.Unmarshal(respBytes, &errResp)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), services.ItemNotBidEligible.Error(), errResp.Error)
}

func (suite *ItemTestSuite) TestDontAllowItemOwnerBids() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	bidPayload := `
{
	"bidValue": 12
}
`
	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/bid", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		bytes.NewReader([]byte(bidPayload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)
	var errResp response.HttpErrorResponse
	err = json.Unmarshal(respBytes, &errResp)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), services.BidsNotAllowedByOwner.Error(), errResp.Error)
}

func (suite *ItemTestSuite) TestAllowMarkItemOffBidByAuthor() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.False(suite.T(), item.IsOffBid())
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))

	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/mark-off-bid", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"PUT",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	item = suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.True(suite.T(), item.IsOffBid())
}

func (suite *ItemTestSuite) TestDontAllowMarkItemOffBidByNonAuthor() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.False(suite.T(), item.IsOffBid())
	actionUser := suite.repository.FindUserById(6)
	assert.False(suite.T(), item.IsOwner(*actionUser))

	token, err := services.GenerateNewJwtToken(actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/mark-off-bid", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"PUT",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)

	item = suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.False(suite.T(), item.IsOffBid())
}

func (suite *ItemTestSuite) TestAddItemImagesByAuthor() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	assert.False(suite.T(), item.IsOffBid())
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))

	payload, contentType, err := client.MakeMultiPartWriterFromFiles(
		"images", suite.itemImageFile1, suite.itemImageFile2,
	)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": contentType},
		time.Second*10,
		payload,
	)
	defer resp.Body.Close()
	var itemImages []*models.ItemImage

	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)

	err = json.Unmarshal(respBytes, &itemImages)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	assert.Len(suite.T(), itemImages, 2)
	assert.NotNil(suite.T(), itemImages[0].ID)
	assert.NotEqual(suite.T(), "", itemImages[0].Name)

	itemImages = suite.repository.FindItemImages(item)
	assert.Len(suite.T(), itemImages, 2)
	assert.NotNil(suite.T(), itemImages[0].ID)
	assert.NotEqual(suite.T(), "", itemImages[0].Name)
}

func (suite *ItemTestSuite) TestAddItemImagesByRemoveExisting() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	assert.False(suite.T(), item.IsOffBid())
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))

	payload, contentType, err := client.MakeMultiPartWriterFromFiles(
		"images", suite.itemImageFile1,
	)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images?removeExisting=true", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": contentType},
		time.Second*10,
		payload,
	)
	var itemImages []*models.ItemImage

	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)
	err = json.Unmarshal(respBytes, &itemImages)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	assert.Len(suite.T(), itemImages, 1)
	assert.NotNil(suite.T(), itemImages[0].ID)
	assert.NotEqual(suite.T(), "", itemImages[0].Name)
	assert.Contains(suite.T(), itemImages[0].Name, "guitar_1")

	itemImages = suite.repository.FindItemImages(item)
	assert.Len(suite.T(), itemImages, 1)
	assert.NotNil(suite.T(), itemImages[0].ID)
	assert.NotEqual(suite.T(), "", itemImages[0].Name)
	assert.Contains(suite.T(), itemImages[0].Name, "guitar_1")
	resp.Body.Close()

	// upload file 2.
	payload, contentType, err = client.MakeMultiPartWriterFromFiles(
		"images", suite.itemImageFile2,
	)
	assert.Nil(suite.T(), err)

	resp, err = client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images?removeExisting=true", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": contentType},
		time.Second*10,
		payload,
	)
	defer resp.Body.Close()

	respBytes, err = io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)
	err = json.Unmarshal(respBytes, &itemImages)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	assert.Len(suite.T(), itemImages, 1)
	assert.NotNil(suite.T(), itemImages[0].ID)
	assert.NotEqual(suite.T(), "", itemImages[0].Name)
	assert.Contains(suite.T(), itemImages[0].Name, "guitar_2")

	itemImages = suite.repository.FindItemImages(item)
	assert.Len(suite.T(), itemImages, 1)
	assert.NotNil(suite.T(), itemImages[0].ID)
	assert.NotEqual(suite.T(), "", itemImages[0].Name)
	assert.Contains(suite.T(), itemImages[0].Name, "guitar_2")
}

func (suite *ItemTestSuite) TestRemoveItemImageByAuthor() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, "581b7c3c-3fa8-4642-801b-30f63111f621",'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	suite.DB.Exec(`
INSERT INTO item_images (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, path, item_id) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'guitar_2_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_2_abc.jpg", 1)
;
`)

	itemImg := suite.repository.FindItemImage(1, 1)
	assert.NotNil(suite.T(), itemImg)
	assert.NotNil(suite.T(), itemImg.Item)
	assert.False(suite.T(), itemImg.Item.IsZero())
	assert.False(suite.T(), itemImg.Item.IsOffBid())
	assert.True(suite.T(), itemImg.Item.IsOwner(*suite.actionUser))

	image1Bytes, err := io.ReadAll(suite.itemImageFile1)
	assert.Nil(suite.T(), err)
	dirPath := fmt.Sprintf("%s/items/%s/images", utils.GetGlobalUploadsDir(), itemImg.Item.Uuid)
	filePath := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), itemImg.Path)
	err = os.MkdirAll(dirPath, 0755)
	assert.Nil(suite.T(), err)

	err = os.WriteFile(filePath, image1Bytes, 0755)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemImg.ItemId, itemImg.Item.ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": suite.contentTypeJson},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	itemImg = suite.repository.FindItemImage(1, 1)
	assert.Nil(suite.T(), itemImg)
}

func (suite *ItemTestSuite) TestRemoveItemImagesByAuthor() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, "581b7c3c-3fa8-4642-801b-30f63111f621",'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	suite.DB.Exec(`
INSERT INTO item_images (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, path, item_id) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'guitar_2_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_2_abc.jpg", 1),
(2, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'guitar_1_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_1_abc.jpg", 1)
;
`)

	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))

	itemImages := suite.repository.FindItemImages(item)
	assert.NotNil(suite.T(), itemImages)
	assert.Len(suite.T(), itemImages, 2)

	image1Bytes, err := io.ReadAll(suite.itemImageFile1)
	assert.Nil(suite.T(), err)
	dirPath := fmt.Sprintf("%s/items/%s/images", utils.GetGlobalUploadsDir(), item.Uuid)
	err = os.MkdirAll(dirPath, 0755)
	assert.Nil(suite.T(), err)
	image2Bytes, err := io.ReadAll(suite.itemImageFile2)
	assert.Nil(suite.T(), err)

	file1Path := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), itemImages[0].Path)
	err = os.WriteFile(file1Path, image1Bytes, 0755)
	assert.Nil(suite.T(), err)
	file2Path := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), itemImages[1].Path)
	err = os.WriteFile(file2Path, image2Bytes, 0755)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": suite.contentTypeJson},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	itemImages = suite.repository.FindItemImages(item)
	assert.Nil(suite.T(), itemImages)
	assert.Len(suite.T(), itemImages, 0)
}

func (suite *ItemTestSuite) TestDontAllowRemoveItemImagesByNonAuthor() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, "581b7c3c-3fa8-4642-801b-30f63111f621",'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	suite.DB.Exec(`
INSERT INTO item_images (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, path, item_id) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'guitar_2_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_2_abc.jpg", 1),
(2, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'guitar_1_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_1_abc.jpg", 1)
;
`)

	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	actionUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), actionUser)
	assert.False(suite.T(), item.IsOwner(*actionUser))

	itemImages := suite.repository.FindItemImages(item)
	assert.NotNil(suite.T(), itemImages)
	assert.Len(suite.T(), itemImages, 2)

	dirPath := fmt.Sprintf("%s/items/%s/images", utils.GetGlobalUploadsDir(), item.Uuid)
	err := os.MkdirAll(dirPath, 0755)
	assert.Nil(suite.T(), err)

	image1Bytes, err := io.ReadAll(suite.itemImageFile1)
	assert.Nil(suite.T(), err)
	image2Bytes, err := io.ReadAll(suite.itemImageFile2)
	assert.Nil(suite.T(), err)

	file1Path := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), itemImages[0].Path)
	err = os.WriteFile(file1Path, image1Bytes, 0755)
	assert.Nil(suite.T(), err)
	file2Path := fmt.Sprintf("%s/%s", utils.GetGlobalUploadsDir(), itemImages[1].Path)
	err = os.WriteFile(file2Path, image2Bytes, 0755)
	assert.Nil(suite.T(), err)

	token, err := services.GenerateNewJwtToken(actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	resp.Body.Close()
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)
	assert.Nil(suite.T(), err)

	// trying to delete single item.
	itemImages = suite.repository.FindItemImages(item)
	assert.Len(suite.T(), itemImages, 2)

	resp, err = client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID, itemImages[0].ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)
	assert.Nil(suite.T(), err)

	// verify that items still exist.
	itemImages = suite.repository.FindItemImages(item)
	assert.Len(suite.T(), itemImages, 2)
}

func (suite *ItemTestSuite) TestDontAllowAddMoreThanMaxItemImages() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	assert.False(suite.T(), item.IsOffBid())
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))

	uploadFiles := []*os.File{
		suite.itemImageFile1, suite.itemImageFile2, suite.itemImageFile3, suite.itemImageFile4,
		suite.itemImageFile5, suite.itemImageFile6,
	}
	payload, contentType, err := client.MakeMultiPartWriterFromFiles("images", uploadFiles...)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": contentType},
		time.Second*10,
		payload,
	)
	_, err = io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	itemImages := suite.repository.FindItemImages(item)
	assert.Nil(suite.T(), itemImages)
}

func (suite *ItemTestSuite) TestAllowAddMaxItemImages() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	assert.False(suite.T(), item.IsOffBid())
	assert.True(suite.T(), item.IsOwner(*suite.actionUser))

	uploadFiles := []*os.File{
		suite.itemImageFile1, suite.itemImageFile2, suite.itemImageFile3, suite.itemImageFile4,
		suite.itemImageFile5,
	}

	payload, contentType, err := client.MakeMultiPartWriterFromFiles("images", uploadFiles...)
	assert.Nil(suite.T(), err)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken, "Content-Type": contentType},
		time.Second*10,
		payload,
	)

	defer resp.Body.Close()

	var itemImages []*models.ItemImage
	respBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	err = json.Unmarshal(respBytes, &itemImages)
	assert.Nil(suite.T(), err)

	countThumbnail := 0
	var thumbnail *models.ItemImage
	item = suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item.ItemImages)
	assert.Len(suite.T(), item.ItemImages, models.MaxImagesPerItem)
	for _, itemImg := range item.ItemImages {
		if itemImg.IsThumbnail {
			countThumbnail += 1
			thumbnail = itemImg
		}
	}

	assert.Equal(suite.T(), 1, countThumbnail)
	// fetch file.
	resp, err = client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/images/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, item.ID, thumbnail.ID),
		"GET",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		payload,
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(suite.T(), "application/octet-stream", resp.Header.Get("Content-Type"))
	assert.Equal(suite.T(), fmt.Sprintf("attachment; filename=%s", thumbnail.Name), resp.Header.Get("Content-Disposition"))
	assert.Equal(suite.T(), "binary", resp.Header.Get("Content-Transfer-Encoding"))
	assert.Equal(suite.T(), "no-cache", resp.Header.Get("Cache-Control"))
}

func (suite *ItemTestSuite) TestFetchItemBids() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) VALUES 
(6, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1),
(7, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe26@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
;
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)

	suite.DB.Exec(`
INSERT INTO bids (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, item_id, value) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,6,6,NULL, 1, 21000),
(2, uuid_v4(),'2022-04-07 06:46:03.528','2022-04-07 06:46:03.528',NULL,1,7,7,NULL, 1, 21100)
;
`)
	user6 := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), user6)
	user7 := suite.repository.FindUserById(7)
	assert.NotNil(suite.T(), user7)

	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	assert.False(suite.T(), item.IsOffBid())
	assert.False(suite.T(), item.IsOwner(*user6))
	assert.False(suite.T(), item.IsOwner(*user7))
	var bids []*models.Bid
	suite.repository.FindBidsByItem(item).Find(&bids)
	assert.NotNil(suite.T(), bids)
	assert.Len(suite.T(), bids, 2)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s/%d/bids", suite.itemsRoute, item.ID),
		"GET",
		map[string]string{},
		map[string]string{},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")

	var page paginate.Page
	responseBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read from response message body.")
	err = json.Unmarshal(responseBytes, &page)
	assert.Nil(suite.T(), err, "Could not parse items list from response message.")
	pageItems, err := json.Marshal(page.Items)
	assert.Nil(suite.T(), err, "Could not serialize payload to bytes.")
	err = json.Unmarshal(pageItems, &bids)
	assert.Nil(suite.T(), err, "Could not deserialize items.")
	assert.Len(suite.T(), bids, 2)
	assert.Equal(suite.T(), bids[0].Value, models.Value(21100))
	assert.Equal(suite.T(), bids[0].UserCreatedBy, user7.ID)
	assert.Equal(suite.T(), bids[1].Value, models.Value(21000))
	assert.Equal(suite.T(), bids[1].UserCreatedBy, user6.ID)
}
