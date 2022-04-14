package services

import (
	"github.com/ashishkumar68/auction-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *ServiceTestSuite) TestGenerateNewToken() {
	user := models.User{
		BaseModel: models.BaseModel{
			ID:   1,
			Uuid: uuid.NewString(),
		},
		FirstName: "John",
		LastName:  "Smith",
		Email:     "johnsmith@abc.com",
		IsActive:  true,
	}
	tokenString, err := GenerateNewJwtToken(user, TokenTypeAccess)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), tokenString != "")
	tokenVal, err := VerifyJwtToken(tokenString)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), tokenVal.Valid)
}
