package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
)

type RegisterResponseBody struct {
	models.BaseModel

	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	Email		string	`json:"email"`
	IsActive	bool	`json:"isActive"`
}

type LoginResponseBody struct {
	models.LoggedInUser
}

var _ = Describe("Auth Tests", func() {
	protocol := "http"
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	prefix := "/api"
	registerRoute := fmt.Sprintf("%s://%s:%s%s/register", protocol, host, port, prefix)
	loginRoute := fmt.Sprintf("%s://%s:%s%s/login", protocol, host, port, prefix)

	contentTypeJson := "application/json"
	var dbConnection *gorm.DB
	var userRepository *repositories.UserRepository
	cleanUpTables := func() {
		dbConnection = database.NewConnectionWithContext(context.TODO())
		userRepository = repositories.NewUserRepository(dbConnection)
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
				Expect(registerResponse.IsActive).To(Equal(true))

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

	Context("I should be able to login with correct credentials", func() {
		BeforeEach(func() {
			dbConnection.Exec(`
INSERT INTO users(uuid, created_at, updated_at, first_name, last_name, email, password) 
VALUES ("40449c42-1a4d-4dad-b942-48ded845329e", NOW(), NOW(), "John", "Smith", "johnsmith24@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6")`)
		})

		It("should not allow login with non-existing account", func() {
			invalidEmail := "blablaemail@abc.bla"
			existingUser := userRepository.FindByEmail(invalidEmail)
			Expect(existingUser.IsZero()).To(BeTrue())
			payload := `
{
    "email": "blablaemail@abc.bla",
    "password": "blabla"
}
`
			resp, err := http.Post(loginRoute, contentTypeJson, bytes.NewReader([]byte(payload)))
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			defer resp.Body.Close()
			Expect(err).To(BeNil(), "Could not Parse HTTP message response.")
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should not allow login with incorrect credentials", func() {
			existingUser := userRepository.FindByEmail("johnsmith24@abc.com")
			Expect(existingUser.IsZero()).To(BeFalse())
			payload := `
{
    "email": "johnsmith24@abc.com",
    "password": "blabla"
}
`
			resp, err := http.Post(loginRoute, contentTypeJson, bytes.NewReader([]byte(payload)))
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			defer resp.Body.Close()
			Expect(err).To(BeNil(), "Could not Parse HTTP message response.")
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should allow login with correct credentials", func() {
			existingUser := userRepository.FindByEmail("johnsmith24@abc.com")
			Expect(existingUser.IsZero()).To(BeFalse())
			payload := `
{
    "email": "johnsmith24@abc.com",
    "password": "secret12"
}
`
			resp, err := http.Post(loginRoute, contentTypeJson, bytes.NewReader([]byte(payload)))
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			var loginResponse LoginResponseBody
			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			err = json.Unmarshal(body, &loginResponse)
			Expect(err).To(BeNil(), "Could not Parse HTTP message response.")
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(loginResponse.FirstName).To(Equal("John"))
			Expect(loginResponse.LastName).To(Equal("Smith"))
			Expect(loginResponse.Email).To(Equal("johnsmith24@abc.com"))
			Expect(loginResponse.AccessToken).ToNot(Equal(""))
			Expect(loginResponse.RefreshToken).ToNot(Equal(""))
		})
	})
})