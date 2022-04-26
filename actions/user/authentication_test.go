package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
)

type RegisterResponseBody struct {
	models.BaseModel

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	IsActive  bool   `json:"isActive"`
}

type LoginResponseBody struct {
	models.LoggedInUser
}

func (suite *UserTestSuite) TestRegisterNewUser() {
	registerRoute := fmt.Sprintf("%s://%s:%s%s/user/register", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)
	email := "johndoe123@abc.com"

	var user models.User
	suite.DB.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	assert.True(suite.T(), user.IsZero())

	payload := `
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "johndoe123@abc.com",
  "password": "secret12"
}
`
	resp, err := http.Post(registerRoute, suite.contentTypeJson, bytes.NewReader([]byte(payload)))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")

	var registerResponse RegisterResponseBody
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(body, &registerResponse)

	assert.Nil(suite.T(), err, "Could not Parse HTTP message response.")
	assert.Equal(suite.T(), "John", registerResponse.FirstName)
	assert.Equal(suite.T(), "Doe", registerResponse.LastName)
	assert.Equal(suite.T(), "johndoe123@abc.com", registerResponse.Email)
	assert.True(suite.T(), registerResponse.IsActive)

	suite.DB.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	assert.False(suite.T(), user.IsZero())
	assert.NotEqual(suite.T(), "secret12", user.Password)
	assert.Equal(suite.T(), "", user.PlainPassword)

	// It should not allow creating user with same email.
	resp, err = http.Post(registerRoute, suite.contentTypeJson, bytes.NewReader([]byte(payload)))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not Parse HTTP message response.")
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *UserTestSuite) TestLogin() {
	suite.DB.Exec(`
INSERT INTO users(uuid, created_at, updated_at, first_name, last_name, email, password) 
VALUES (uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmithtest@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6")
`)
	loginRoute := fmt.Sprintf("%s://%s:%s%s/user/login", suite.protocol, suite.host, suite.port, suite.apiBaseRoute)

	invalidEmail := "blablaemail@abc.bla"
	existingUser := suite.repository.FindUserByEmail(invalidEmail)
	assert.Nil(suite.T(), existingUser)
	payload := `
{
    "email": "blablaemail@abc.bla",
    "password": "blabla"
}
`
	resp, err := http.Post(loginRoute, suite.contentTypeJson, bytes.NewReader([]byte(payload)))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not Parse HTTP message response.")
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	existingUser = suite.repository.FindUserByEmail("johnsmithtest@abc.com")
	assert.False(suite.T(), existingUser.IsZero())
	payload = `
{
    "email": "johnsmithtest@abc.com",
    "password": "blabla"
}
`
	resp, err = http.Post(loginRoute, suite.contentTypeJson, bytes.NewReader([]byte(payload)))
	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	defer resp.Body.Close()
	assert.Nil(suite.T(), err, "Could not Parse HTTP message response.")
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	payload = `
{
    "email": "johnsmithtest@abc.com",
    "password": "secret12"
}
`
	resp, err = http.Post(loginRoute, suite.contentTypeJson, bytes.NewReader([]byte(payload)))

	assert.Nil(suite.T(), err, "Could not detect service available.")
	assert.NotNil(suite.T(), resp, "Could not detect service available.")
	var loginResponse LoginResponseBody
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(body, &loginResponse)
	assert.Nil(suite.T(), err, "Could not Parse HTTP message response.")
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(suite.T(), "John", loginResponse.FirstName)
	assert.Equal(suite.T(), "Smith", loginResponse.LastName)
	assert.Equal(suite.T(), "johnsmithtest@abc.com", loginResponse.Email)
	assert.NotEqual(suite.T(), "", loginResponse.AccessToken)
	assert.NotEqual(suite.T(), "", loginResponse.RefreshToken)
}
