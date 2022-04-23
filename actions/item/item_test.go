package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/morkid/paginate"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
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
	assert.Nil(suite.T(), err, "Could not generate new token for create item test.")

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
		[]byte(payload),
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
	assert.Nil(suite.T(), err, "Could not generate new token for create item test.")
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
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)
	items := suite.repository.FindItemByName("ABC Item")
	assert.Len(suite.T(), items, 2)

	bidPayload := `
{
	"bidValue": 12
}
`
	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err, "Could not generate new token for create item test.")
	itemId := uint(1)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/bid", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		[]byte(bidPayload),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	item := suite.repository.FindItemById(itemId)
	newBid := suite.repository.FindBidByItem(item, suite.actionUser)
	assert.NotNil(suite.T(), newBid)
	assert.False(suite.T(), newBid.IsZero())
	assert.Equal(suite.T(), newBid.Value, models.Value(12))
}
