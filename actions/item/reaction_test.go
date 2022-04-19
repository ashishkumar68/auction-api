package item

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"time"
)

var _ = Describe("Reaction Tests", func() {
	protocol := "http"
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	prefix := "/api"
	baseItemsRoute := fmt.Sprintf("%s://%s:%s%s/items", protocol, host, port, prefix)

	contentType := "application/json"
	var dbConnection *gorm.DB
	var repository *repositories.Repository

	var user *models.User
	cleanUpTables := func() {
		dbConnection = database.GetDBHandle().WithContext(context.TODO())
		repository = repositories.NewRepository(dbConnection)
		migrations.ForceTruncateAllTables(dbConnection)

		dbConnection.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) 
VALUES (1, uuid_v4(), NOW(), NOW(), "John", "Smith", "johnsmith24@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
`)
		dbConnection.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'ABC Item 1','Item 1 Description','1','ABC','20000', "2099-01-01");
`)
		user = repository.FindUserByEmail("johnsmith24@abc.com")
	}
	BeforeEach(func() {
		cleanUpTables()

		Expect(user).To(Not(BeNil()))
	})
	AfterEach(func() {
		migrations.ForceTruncateAllTables(dbConnection)
	})

	Context("Add reaction to items", func() {
		var item *models.Item
		var addNewItemReactionRoute string
		var newToken string
		BeforeEach(func() {
			item = repository.FindItemById(1)
			Expect(item).To(Not(BeNil()))

			reaction := repository.FindReactionByItemAndUser(item, user)
			Expect(reaction).To(BeNil())
			addNewItemReactionRoute = fmt.Sprintf("%s/%d/reaction", baseItemsRoute, item.ID)
			token, err := services.GenerateNewJwtToken(user, services.TokenTypeAccess)
			Expect(err).To(BeNil(), "Could not generate new token for create item test.")
			newToken = token
		})

		It("should allow adding a new reaction to an item by a user", func() {
			payload := `{"reactionType": 0}`
			resp, err := client.MakeRequest(
				addNewItemReactionRoute,
				"POST",
				map[string]string{},
				map[string]string{"Authorization": newToken, "Content-Type": contentType},
				time.Second*10,
				[]byte(payload),
			)
			defer resp.Body.Close()
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			responseBytes, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil(), "Could not read from API response.")
			var reaction models.Reaction
			err = json.Unmarshal(responseBytes, &reaction)
			Expect(err).To(BeNil(), "Could not parse API response.")
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			newR := repository.FindReactionByItemAndUser(item, user)
			Expect(newR).To(Not(BeNil()))
			Expect(newR.Type).To(Equal(uint8(models.ReactionTypeLike)))
		})

		It("should allow updating an existing reaction on an item by a user", func() {
			dbConnection.Exec(`
INSERT INTO reactions (uuid,created_at,updated_at,deleted_at,version,created_by,updated_by,deleted_by,item_id,type) VALUES 
(uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, NULL, 1, 0)
;
`)
			existingReaction := repository.FindReactionByItemAndUser(item, user)
			Expect(existingReaction).To(Not(BeNil()))
			Expect(existingReaction.Type).To(Equal(uint8(models.ReactionTypeLike)))

			payload := `{"reactionType": 1}`
			resp, err := client.MakeRequest(
				addNewItemReactionRoute,
				"POST",
				map[string]string{},
				map[string]string{"Authorization": newToken, "Content-Type": contentType},
				time.Second*10,
				[]byte(payload),
			)
			defer resp.Body.Close()
			Expect(err).To(BeNil(), "Could not detect service available.")
			Expect(resp).To(Not(BeNil()), "Could not detect service available.")
			responseBytes, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil(), "Could not read from API response.")
			var reaction models.Reaction
			err = json.Unmarshal(responseBytes, &reaction)
			Expect(err).To(BeNil(), "Could not parse API response.")
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			newR := repository.FindReactionByItemAndUser(item, user)
			Expect(newR).To(Not(BeNil()))
			Expect(newR.Type).To(Equal(uint8(models.ReactionTypeDislike)))
		})
	})
})
