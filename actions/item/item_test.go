package item

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/morkid/paginate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"time"
)

var _ = Describe("Item Tests", func() {
	protocol := "http"
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	prefix := "/api"
	itemsRoute := fmt.Sprintf("%s://%s:%s%s/items", protocol, host, port, prefix)

	contentTypeJson := "application/json"
	var dbConnection *gorm.DB
	var itemRepository *repositories.ItemRepository
	var userRepository *repositories.UserRepository
	var user *models.User
	cleanUpTables := func() {
		dbConnection = database.NewConnectionWithContext(context.TODO())
		itemRepository = repositories.NewItemRepository(dbConnection)
		userRepository = repositories.NewUserRepository(dbConnection)
		dbConnection.Exec(`SET foreign_key_checks = 0;`)
		dbConnection.Exec(`TRUNCATE TABLE users;`)
		dbConnection.Exec(`TRUNCATE TABLE items;`)
		dbConnection.Exec(`TRUNCATE TABLE bids;`)
		dbConnection.Exec(`SET foreign_key_checks = 1;`)

		dbConnection.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (1, uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmith24@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
		user = userRepository.FindByEmail("johnsmith24@abc.com")

	}
	BeforeEach(func() {
		cleanUpTables()
	})

	Context("Don't allow create item as an Anonymous user.", func() {
		BeforeEach(func() {
			items := itemRepository.FindByName("ABC Washing Machine")
			Expect(items).To(BeEmpty())
			Expect(user.IsZero()).To(BeFalse())
		})

		payload := `
{
    "name": "ABC Washing Machine",
    "description": "A washing machine",
    "category": 1,
    "brandName": "ABC",
    "marketValue": 20000
}
`

		It("should not allow new item as anonymous user in system.", func() {
			resp, err := http.Post(itemsRoute, contentTypeJson, bytes.NewReader([]byte(payload)))
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})

		It("should allow creating item as a logged in user", func() {
			token, err := services.GenerateNewJwtToken(user, services.TokenTypeAccess)
			Expect(err).To(BeNil(), "Could not generate new token for create item test.")
			resp, err := client.MakeRequest(
				itemsRoute,
				"POST",
				map[string]string{},
				map[string]string{"Authorization": token},
				time.Second*10,
				[]byte(payload),
			)
			defer resp.Body.Close()

			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			var item models.Item
			body, err := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &item)

			Expect(err).To(BeNil(), "Could not Parse HTTP message response.")
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			items := itemRepository.FindByName("ABC Washing Machine")
			Expect(len(items)).To(Equal(1))

			actualItem := itemRepository.FindByUuid(item.Uuid)
			Expect(actualItem).ToNot(BeNil())
			Expect(actualItem.Name).To(ContainSubstring("ABC Washing Machine"))
			//createdBy := userRepository.Find(*actualItem.User)
			Expect(actualItem.UserCreated).ToNot(BeNil())
			Expect(actualItem.UserCreated.Email).To(Equal("johnsmith24@abc.com"))
		})
	})

	Context("I should be able to list the items anonymously as well as logged in user.", func() {
		BeforeEach(func() {
			dbConnection.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value) VALUES 
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000'),
(2, uuid_v4(),'2022-04-06 06:46:03.528','2022-04-06 06:46:03.528',NULL,1,1,1,NULL,'ABC Item 2','Item 2 Description','1','ABC','22000');
`)
			items := itemRepository.FindByName("ABC Item")
			Expect(len(items)).To(Equal(2))
		})

		It("should allow listing items anonymously.", func() {
			resp, err := client.MakeRequest(
				itemsRoute,
				"GET",
				map[string]string{},
				map[string]string{},
				time.Second*10,
				nil,
			)
			defer resp.Body.Close()
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			var page paginate.Page
			responseBytes, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil(), "Could not read from response message body.")
			err = json.Unmarshal(responseBytes, &page)
			Expect(err).To(BeNil(), "Could not parse items list from response message.")
			pageItems, err := json.Marshal(page.Items)
			Expect(err).To(BeNil(), "Could not serialize payload to bytes.")
			var items []models.Item
			err = json.Unmarshal(pageItems, &items)
			Expect(err).To(BeNil(), "Could not deserialize items.")

			Expect(len(items)).To(Equal(2))
			Expect(items[0].Name).To(ContainSubstring("ABC Item"))
			Expect(items[1].Name).To(ContainSubstring("ABC Item"))
		})

		It("should allow listing items with authenticated user.", func() {
			token, err := services.GenerateNewJwtToken(user, services.TokenTypeAccess)
			Expect(err).To(BeNil(), "Could not generate new token for create item test.")
			resp, err := client.MakeRequest(
				itemsRoute,
				"GET",
				map[string]string{},
				map[string]string{"Authorization": token},
				time.Second*10,
				nil,
			)
			defer resp.Body.Close()
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			var page paginate.Page
			responseBytes, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil(), "Could not read from response message body.")
			err = json.Unmarshal(responseBytes, &page)
			Expect(err).To(BeNil(), "Could not parse items list from response message.")
			pageItems, err := json.Marshal(page.Items)
			Expect(err).To(BeNil(), "Could not serialize payload to bytes.")
			var items []models.Item
			err = json.Unmarshal(pageItems, &items)
			Expect(err).To(BeNil(), "Could not deserialize items.")

			Expect(len(items)).To(Equal(2))
			Expect(items[0].Name).To(ContainSubstring("ABC Item"))
			Expect(items[1].Name).To(ContainSubstring("ABC Item"))
		})
	})
})
