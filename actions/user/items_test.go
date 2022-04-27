package user

import (
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/morkid/paginate"
	"github.com/stretchr/testify/assert"
	"io"
	"time"
)

func (suite *UserTestSuite) TestItListsItemsAddedByUser() {
	suite.DB.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (2, uuid_v4(), NOW(), NOW(), "John", "Doe", "johndoe25@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
;
`)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,1,1,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01"),
(3, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,2,2,NULL,'ABC Item 3','Item 3 Description','1','ABC','24000', "2099-02-01")
;
`)

	user2 := suite.repository.FindUserById(2)
	assert.NotNil(suite.T(), user2)
	// item 1
	item1 := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item1)
	assert.True(suite.T(), item1.IsOwner(*suite.actionUser))
	assert.False(suite.T(), item1.IsOwner(*user2))
	// item 2
	item2 := suite.repository.FindItemById(2)
	assert.NotNil(suite.T(), item2)
	assert.True(suite.T(), item2.IsOwner(*suite.actionUser))
	assert.False(suite.T(), item2.IsOwner(*user2))
	// item 3
	item3 := suite.repository.FindItemById(3)
	assert.NotNil(suite.T(), item3)
	assert.False(suite.T(), item3.IsOwner(*suite.actionUser))
	assert.True(suite.T(), item3.IsOwner(*user2))

	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s/items", suite.userRoute),
		"GET",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")

	var items []models.Item
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
	assert.True(suite.T(), items[0].IsOwner(*suite.actionUser))
	assert.True(suite.T(), items[1].IsOwner(*suite.actionUser))
	assert.False(suite.T(), items[0].IsOwner(*user2))
	assert.False(suite.T(), items[1].IsOwner(*user2))
	resp.Body.Close()

	// user 2 items
	token, err = services.GenerateNewJwtToken(user2, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)
	resp, err = client.MakeRequest(
		fmt.Sprintf("%s/items", suite.userRoute),
		"GET",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")

	responseBytes, err = io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read from response message body.")
	err = json.Unmarshal(responseBytes, &page)
	assert.Nil(suite.T(), err, "Could not parse items list from response message.")
	pageItems, err = json.Marshal(page.Items)
	assert.Nil(suite.T(), err, "Could not serialize payload to bytes.")
	err = json.Unmarshal(pageItems, &items)
	assert.Nil(suite.T(), err, "Could not deserialize items.")
	assert.Len(suite.T(), items, 1)
	assert.True(suite.T(), items[0].IsOwner(*user2))
	assert.False(suite.T(), items[0].IsOwner(*suite.actionUser))
}
