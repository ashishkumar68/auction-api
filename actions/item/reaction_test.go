package item

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"time"
)

func (suite *ItemTestSuite) TestAddReactionToItems() {
	baseItemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01");
`)

	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)

	reaction := suite.repository.FindReactionByItemAndUser(item, suite.actionUser)
	assert.Nil(suite.T(), reaction)
	addNewItemReactionRoute := fmt.Sprintf("%s/%d/reaction", baseItemsRoute, item.ID)
	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)

	payload := `{"reactionType": 0}`
	resp, err := client.MakeRequest(
		addNewItemReactionRoute,
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token, "Content-Type": suite.contentTypeJson},
		time.Second*10,
		bytes.NewReader([]byte(payload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	responseBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read from API response.")
	var newReaction models.Reaction
	err = json.Unmarshal(responseBytes, &newReaction)
	assert.Nil(suite.T(), err, "Could not parse API response.")
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	newR := suite.repository.FindReactionByItemAndUser(item, suite.actionUser)
	assert.NotNil(suite.T(), newR)
	assert.Equal(suite.T(), uint8(models.ReactionTypeLike), newR.Type)
}

func (suite *ItemTestSuite) TestUpdateItemsReaction() {
	baseItemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01");
`)
	suite.DB.Exec(`
INSERT INTO reactions (uuid,created_at,updated_at,deleted_at,version,created_by,updated_by,deleted_by,item_id,type) VALUES 
(uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, NULL, 1, 0)
;
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	existingReaction := suite.repository.FindReactionByItemAndUser(item, suite.actionUser)
	assert.NotNil(suite.T(), existingReaction)
	assert.Equal(suite.T(), uint8(models.ReactionTypeLike), existingReaction.Type)

	addNewItemReactionRoute := fmt.Sprintf("%s/%d/reaction", baseItemsRoute, item.ID)

	token, err := services.GenerateNewJwtToken(suite.actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err)

	payload := `{"reactionType": 1}`
	resp, err := client.MakeRequest(
		addNewItemReactionRoute,
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token, "Content-Type": suite.contentTypeJson},
		time.Second*10,
		bytes.NewReader([]byte(payload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	responseBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "Could not read from API response.")
	var reaction models.Reaction
	err = json.Unmarshal(responseBytes, &reaction)
	assert.Nil(suite.T(), err, "Could not parse API response.")
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	newR := suite.repository.FindReactionByItemAndUser(item, suite.actionUser)
	assert.NotNil(suite.T(), newR)
	assert.Equal(suite.T(), uint8(models.ReactionTypeDislike), newR.Type)
}

func (suite *ItemTestSuite) TestRemoveReactionFromItem() {
	baseItemsRoute := fmt.Sprintf("%s://%s:%s%s/items", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01");
`)
	suite.DB.Exec(`
INSERT INTO reactions (uuid,created_at,updated_at,deleted_at,version,created_by,updated_by,deleted_by,item_id,type) VALUES 
(uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, NULL, 1, 0)
;
`)
	item := suite.repository.FindItemById(1)
	assert.NotNil(suite.T(), item)
	existingReaction := suite.repository.FindReactionByItemAndUser(item, suite.actionUser)
	assert.NotNil(suite.T(), existingReaction)
	assert.Equal(suite.T(), uint8(models.ReactionTypeLike), existingReaction.Type)

	itemReactionRoute := fmt.Sprintf("%s/%d/reaction", baseItemsRoute, item.ID)
	resp, err := client.MakeRequest(
		itemReactionRoute,
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)
	respBytes, _ := io.ReadAll(resp.Body)
	log.Println("RESPONSE: ============>", string(respBytes))

	existingReaction = suite.repository.FindReactionByItemAndUser(item, suite.actionUser)
	assert.Nil(suite.T(), existingReaction)
}
