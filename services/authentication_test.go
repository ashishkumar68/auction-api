package services

import (
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Authentication service tests", func() {
	var dbConnection *gorm.DB
	cleanUpTables := func() {
		dbConnection = database.NewConnection()
		dbConnection.Exec(`SET foreign_key_checks = 0;`)
		dbConnection.Exec(`TRUNCATE TABLE users;`)
		dbConnection.Exec(`TRUNCATE TABLE items;`)
		dbConnection.Exec(`SET foreign_key_checks = 1;`)
	}
	Context("It should generate JWT token for an identity", func() {
		BeforeEach(func() {
			cleanUpTables()
		})
		It("Should generate new JWT token for user which can be verified successfully.", func() {
			user := models.User{
				Identity: models.Identity{
					ID: 1,
					Uuid: uuid.NewString(),
				},
				FirstName: "John",
				LastName: "Smith",
				Email: "johnsmith@abc.com",
				IsActive: true,
			}
			tokenString, err := GenerateNewJwtToken(user, TokenTypeAccess)
			Expect(err).To(BeNil())
			Expect(tokenString).ToNot(Equal(""))
			tokenVal, err := VerifyJwtToken(tokenString)
			Expect(err).To(BeNil())
			Expect(tokenVal.Valid).To(BeTrue())
		})
	})
})