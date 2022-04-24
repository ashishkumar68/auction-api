package item

import (
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"time"
)

func (suite *ItemTestSuite) TestAddItemComment() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)

	items := suite.repository.FindItemByName("ABC Item")
	assert.Len(suite.T(), items, 2)

	payload := `
{
	"comment": "this is a test item comment."
}
`
	itemId := uint(1)
	item := suite.repository.FindItemById(itemId)
	assert.NotNil(suite.T(), item)
	assert.Equal(suite.T(), 0, int(suite.repository.CountCommentsByItem(item)))

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		[]byte(payload),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	var itemComment models.ItemComment
	itemCommentBytes, err := io.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, "could not read from add item comment API response")
	err = json.Unmarshal(itemCommentBytes, &itemComment)
	assert.Nil(suite.T(), err, "could not parse to item comment API response to Item comment")

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	assert.Equal(suite.T(), "this is a test item comment.", itemComment.Description)
	assert.Equal(suite.T(), 1, int(suite.repository.CountCommentsByItem(item)))
}

func (suite *ItemTestSuite) TestUpdateItemComment() {
	suite.DB.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,5,5,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01"),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,5,5,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000', "2099-01-01");
`)

	suite.DB.Exec(`
INSERT INTO item_comments(id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, description, item_id) VALUES
(1, uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, "This is a test item comment.", 1),
(2, uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, "This is another test item comment.", 1)
;
`)

	itemId := uint(1)
	item := suite.repository.FindItemById(itemId)
	assert.NotNil(suite.T(), item)
	item1Comment1 := suite.repository.FindItemCommentById(1)
	item1Comment2 := suite.repository.FindItemCommentById(2)
	assert.NotNil(suite.T(), item1Comment1)
	assert.Equal(suite.T(), "This is a test item comment.", item1Comment1.Description)
	assert.NotNil(suite.T(), item1Comment2)
	assert.Equal(suite.T(), "This is another test item comment.", item1Comment2.Description)

	payload := `
{
	"comment": "this is updated test item comment."
}
`
	// updating comment 1
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId, item1Comment1.ID),
		"PATCH",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		[]byte(payload),
	)
	resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// updating comment 2
	resp, err = client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId, item1Comment2.ID),
		"PATCH",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		[]byte(payload),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	item1Comment1 = suite.repository.FindItemCommentById(1)
	item1Comment2 = suite.repository.FindItemCommentById(2)
	assert.NotNil(suite.T(), item1Comment1)
	assert.Equal(suite.T(), "this is updated test item comment.", item1Comment1.Description)
	assert.NotNil(suite.T(), item1Comment2)
	assert.Equal(suite.T(), "this is updated test item comment.", item1Comment2.Description)
}
