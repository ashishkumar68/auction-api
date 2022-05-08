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
		bytes.NewReader([]byte(payload)),
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
		bytes.NewReader([]byte(payload)),
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
		bytes.NewReader([]byte(payload)),
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

func (suite *ItemTestSuite) TestItShouldNotAllowUpdateItemCommentWhenEditUserIsDifferent() {
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
INSERT INTO item_comments(id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, description, item_id) VALUES
(1, uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, "This is a test item comment.", 1),
(2, uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, "This is another test item comment.", 1)
;
`)

	itemId := uint(1)
	item := suite.repository.FindItemById(itemId)
	assert.NotNil(suite.T(), item)
	item1Comment1 := suite.repository.FindItemCommentById(1)
	assert.NotNil(suite.T(), item1Comment1)
	assert.Equal(suite.T(), "This is a test item comment.", item1Comment1.Description)

	payload := `
{
	"comment": "this is updated test item comment."
}
`
	actionUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), actionUser)
	token, err := services.GenerateNewJwtToken(actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err, "could not generate new token for action user:", actionUser.GetLoginId())
	// Trying to update item 1 comment 1
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId, item1Comment1.ID),
		"PATCH",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		bytes.NewReader([]byte(payload)),
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)

	// Item comment remains intact.
	item1Comment1 = suite.repository.FindItemCommentById(1)
	assert.NotNil(suite.T(), item1Comment1)
	assert.Equal(suite.T(), "This is a test item comment.", item1Comment1.Description)
}

func (suite *ItemTestSuite) TestDeleteItemCommentByAuthor() {
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

	// updating comment 1
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId, item1Comment1.ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		nil,
	)
	resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	// updating comment 2
	resp, err = client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId, item1Comment2.ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": suite.loggedInToken},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode)

	item1Comment1 = suite.repository.FindItemCommentById(1)
	item1Comment2 = suite.repository.FindItemCommentById(2)
	assert.Nil(suite.T(), item1Comment1)
	assert.Nil(suite.T(), item1Comment2)
}

func (suite *ItemTestSuite) TestDoNotAllowDeleteItemCommentByNonAuthor() {
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
INSERT INTO item_comments(id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, description, item_id) VALUES
(1, uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, "This is a test item comment.", 1),
(2, uuid_v4(), NOW(), NOW(), NULL, 1, 5, 5, "This is another test item comment.", 1)
;
`)

	itemId := uint(1)
	item := suite.repository.FindItemById(itemId)
	assert.NotNil(suite.T(), item)
	item1Comment1 := suite.repository.FindItemCommentById(1)
	assert.NotNil(suite.T(), item1Comment1)
	assert.Equal(suite.T(), "This is a test item comment.", item1Comment1.Description)

	actionUser := suite.repository.FindUserById(6)
	assert.NotNil(suite.T(), actionUser)
	token, err := services.GenerateNewJwtToken(actionUser, services.TokenTypeAccess)
	assert.Nil(suite.T(), err, "could not generate new token for action user:", actionUser.GetLoginId())
	// Trying to update item 1 comment 1
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s://%s:%s%s/items/%d/comment/%d", suite.protocol, suite.host, suite.port, suite.apiBaseRoute, itemId, item1Comment1.ID),
		"DELETE",
		map[string]string{},
		map[string]string{"Authorization": token},
		time.Second*10,
		nil,
	)
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "could not connect to add item comment API")
	assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode)

	// Item comment remains intact.
	item1Comment1 = suite.repository.FindItemCommentById(1)
	assert.NotNil(suite.T(), item1Comment1)
	assert.Equal(suite.T(), "This is a test item comment.", item1Comment1.Description)
}

func (suite *ItemTestSuite) TestFetchItemComments() {
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
INSERT INTO item_comments (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, item_id, description) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,6,6,NULL, 1, "This is item 1 comment 1"),
(2, uuid_v4(),'2022-04-07 06:46:03.528','2022-04-07 06:46:03.528',NULL,1,7,7,NULL, 1, "This is item 1 comment 2")
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
	var comments []*models.ItemComment
	suite.repository.FindCommentsByItem(item).Find(&comments)
	assert.NotNil(suite.T(), comments)
	assert.Len(suite.T(), comments, 2)

	resp, err := client.MakeRequest(
		fmt.Sprintf("%s/%d/comments", suite.itemsRoute, item.ID),
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
	err = json.Unmarshal(pageItems, &comments)
	assert.Nil(suite.T(), err, "Could not deserialize items.")
	assert.Len(suite.T(), comments, 2)
	assert.Equal(suite.T(), comments[0].Description, "This is item 1 comment 2")
	assert.Equal(suite.T(), comments[0].UserCreatedBy, user7.ID)
	assert.Equal(suite.T(), comments[1].Description, "This is item 1 comment 1")
	assert.Equal(suite.T(), comments[1].UserCreatedBy, user6.ID)
}
