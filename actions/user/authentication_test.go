package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"net/http"
	"os"
)

type RegisterResponseBody struct {
	models.Identity

	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	Email		string	`json:"email"`
	IsActive	bool	`json:"isActive"`
}

var _ = Describe("Auth Tests", func() {
	protocol := "http"
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	prefix := "/api"
	registerRoute := fmt.Sprintf("%s://%s:%s%s/register", protocol, host, port, prefix)
	contentTypeJson := "application/json"
	var dbConnection *gorm.DB
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	}
	cleanUpTables := func() {
		dbConnection = database.NewConnection(config)
		dbConnection.Exec(`SET foreign_key_checks = 0;`)
		dbConnection.Exec(`TRUNCATE TABLE users;`)
		dbConnection.Exec(`TRUNCATE TABLE items;`)
		dbConnection.Exec(`SET foreign_key_checks = 1;`)
	}
	BeforeEach(func() {
		cleanUpTables()
	})

	Context("I should be able to register as a new user.", func() {
		email := "johndoe123@abc.com"
		It("should not have the user in database initially", func() {
			var user models.User
			dbConnection.Where("email = ? AND deleted_at IS NULL", email).First(&user)
			Expect(user.IsZero()).To(BeTrue())
		})
		When("I request Submit a new register request", func() {
			payload := `
					{
					  "firstName": "John",
					  "lastName": "Doe",
					  "email": "johndoe123@abc.com",
					  "password": "secret12"
					}
`
			It("should create a new user in system.", func() {
				resp, err := http.Post(registerRoute, contentTypeJson, bytes.NewReader([]byte(payload)))
				Expect(err).To(BeNil(), "Could not detect service available.")
				Expect(resp).To(Not(BeNil()), "Could not detect service available.")
				var registerResponse RegisterResponseBody
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				err = json.Unmarshal(body, &registerResponse)
				Expect(err).To(BeNil(), "Could not Parse HTTP message response.")
				Expect(registerResponse.FirstName).To(Equal("John"))
				Expect(registerResponse.LastName).To(Equal("Doe"))
				Expect(registerResponse.Email).To(Equal("johndoe123@abc.com"))
				Expect(registerResponse.IsActive).To(Equal(false))

				var user models.User
				dbConnection.Where("email = ? AND deleted_at IS NULL", email).First(&user)
				Expect(user.IsZero()).To(BeFalse())
				Expect(user.Password).ToNot(Equal("secret12"))
				Expect(user.PlainPassword).To(Equal(""))

				// It should not allow creating user with same email.
				resp, err = http.Post(registerRoute, contentTypeJson, bytes.NewReader([]byte(payload)))
				Expect(err).To(BeNil(), "Could not detect service available.")
				Expect(resp).To(Not(BeNil()), "Could not detect service available.")
				defer resp.Body.Close()
				Expect(err).To(BeNil(), "Could not Parse HTTP message response.")
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})