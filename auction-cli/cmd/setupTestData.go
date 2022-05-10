/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/ashishkumar68/auction-api/client"
	"github.com/ashishkumar68/auction-api/config"
	"github.com/ashishkumar68/auction-api/database"
	"github.com/ashishkumar68/auction-api/migrations"
	"github.com/ashishkumar68/auction-api/models"
	"github.com/ashishkumar68/auction-api/repositories"
	"github.com/ashishkumar68/auction-api/services"
	"github.com/ashishkumar68/auction-api/utils"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	protocol        = "http"
	host            = os.Getenv("HOST")
	port            = os.Getenv("PORT")
	apiBaseRoute    = "/api"
	baseFSItemsPath = fmt.Sprintf("%s/%s", os.Getenv("FILE_UPLOADS_DIR"), models.BaseFSItemsPrefix)
	itemsRoute      = fmt.Sprintf("%s://%s:%s%s/items", protocol, host, port, apiBaseRoute)
)

// setupTestDataCmd represents the setupTestData command
var setupTestDataCmd = &cobra.Command{
	Use:   "setup-test-data",
	Short: "Setup test data",
	Long:  `run auction-cli setup-test-data to setup test records in database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.AppEnvTest != os.Getenv("APP_ENV") {
			fmt.Println("this command is only allowed to run in test environment.")
			return
		}
		addTestData(cmd.Context())
	},
}

func addTestData(c context.Context) {

	db := database.GetDBHandle().WithContext(c)
	repository := repositories.NewRepository(db)
	fmt.Println("Force truncating all tables.")
	migrations.ForceTruncateAllTables(db)
	// clean up items file system.
	err := os.RemoveAll(baseFSItemsPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not clean items upload dir due to err: %s, exiting...", err.Error()))
		return
	}
	fmt.Println("Force truncating complete.")

	fmt.Println("Adding tests data.")
	db.Exec(`
INSERT INTO users(id, uuid, created_at, updated_at, first_name, last_name, email, password, is_active) VALUES 
(1, uuid_v4(), NOW(), NOW(), "John", "Doe 1", "johndoe1@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1),
(2, uuid_v4(), NOW(), NOW(), "John", "Doe 2", "johndoe2@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1),
(3, uuid_v4(), NOW(), NOW(), "John", "Doe 3", "johndoe3@abc.com", "$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUheGU2k.l.5OIf6", 1)
;
`)
	db.Exec(`
INSERT INTO items (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, description, category, brand_name, market_value, last_bid_date) VALUES
(1, "581b7c3c-3fa8-4642-801b-30f63111f621",NOW(),NOW(),NULL,1,1,1,NULL,'ABC Item 1','Item 1 Description','1','ABC','200', "2023-01-01"),
(2, "8c26eb57-40db-4440-8800-f1e854965607",NOW(),NOW(),NULL,1,2,2,NULL,'ABC Item 2','Item 2 Description','0','ABC','200', "2023-02-01"),
(3, "2498c61e-3911-4272-9013-5d21c1e165bc",NOW(),NOW(),NULL,1,3,3,NULL,'ABC Item 3','Item 3 Description','2','ABC','200', "2023-03-01"),
(4, "a06dc759-6a3d-469e-af11-95ded4c3cc90",NOW(),NOW(),NULL,1,1,1,NULL,'ABC Item 4','Item 4 Description','1','ABC','200', "2023-01-01")
;
`)
	db.Exec(`
INSERT INTO bids (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, item_id, value) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL, 1, 210),
(2, uuid_v4(),'2022-04-07 06:46:03.528','2022-04-07 06:46:03.528',NULL,1,2,2,NULL, 1, 211),
(3, uuid_v4(),'2022-04-07 06:46:03.528','2022-04-07 06:46:03.528',NULL,1,3,3,NULL, 1, 212),
(4, uuid_v4(),'2022-04-07 06:46:03.528','2022-04-07 06:46:03.528',NULL,1,1,1,NULL, 2, 211),
(5, uuid_v4(),'2022-04-07 06:46:03.528','2022-04-07 06:46:03.528',NULL,1,2,2,NULL, 2, 212)
;
`)
	db.Exec(`
INSERT INTO reactions (uuid,created_at,updated_at,deleted_at,version,created_by,updated_by,deleted_by,item_id,type) VALUES 
(uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, NULL, 1, 0),
(uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, NULL, 2, 1),
(uuid_v4(), NOW(), NOW(), NULL, 1, 2, 2, NULL, 1, 0),
(uuid_v4(), NOW(), NOW(), NULL, 1, 2, 2, NULL, 2, 1),
(uuid_v4(), NOW(), NOW(), NULL, 1, 3, 3, NULL, 1, 1),
(uuid_v4(), NOW(), NOW(), NULL, 1, 3, 3, NULL, 2, 0),
(uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, NULL, 1, 1),
(uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, NULL, 2, 0)
;
`)
	db.Exec(`
INSERT INTO item_comments(id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, description, item_id) VALUES
(1, uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, "This is test item 1 comment 1.", 1),
(2, uuid_v4(), NOW(), NOW(), NULL, 1, 2, 2, "This is test item 1 comment 2.", 1),
(3, uuid_v4(), NOW(), NOW(), NULL, 1, 3, 3, "This is test item 1 comment 3.", 1),
(4, uuid_v4(), NOW(), NOW(), NULL, 1, 2, 2, "This is test item 1 comment 4.", 1),
(5, uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, "This is test item 2 comment 1.", 2),
(6, uuid_v4(), NOW(), NOW(), NULL, 1, 2, 2, "This is test item 2 comment 2.", 2),
(7, uuid_v4(), NOW(), NOW(), NULL, 1, 3, 3, "This is test item 2 comment 3.", 2),
(8, uuid_v4(), NOW(), NOW(), NULL, 1, 1, 1, "This is test item 2 comment 4.", 2)
;
`)
	db.Exec(`
INSERT INTO item_images (id, uuid, created_at, updated_at, deleted_at, version, created_by, updated_by, deleted_by, name, path, item_id) VALUES
(1, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'guitar_1_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_1_abc.jpg", 1),
(2, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'guitar_2_abc.jpg',"items/581b7c3c-3fa8-4642-801b-30f63111f621/images/guitar_2_abc.jpg", 1),
(3, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'guitar_3_abc.jpg',"items/8c26eb57-40db-4440-8800-f1e854965607/images/guitar_3_abc.jpg", 2),
(4, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'guitar_4_abc.jpg',"items/8c26eb57-40db-4440-8800-f1e854965607/images/guitar_4_abc.jpg", 2),
(5, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'guitar_5_abc.jpg',"items/8c26eb57-40db-4440-8800-f1e854965607/images/guitar_5_abc.jpg", 2),
(6, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'washing_machine_1_abc.png',"items/2498c61e-3911-4272-9013-5d21c1e165bc/images/washing_machine_1_abc.png", 3),
(7, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'washing_machine_2_abc.png',"items/2498c61e-3911-4272-9013-5d21c1e165bc/images/washing_machine_2_abc.png", 3),
(8, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'washing_machine_3_abc.png',"items/2498c61e-3911-4272-9013-5d21c1e165bc/images/washing_machine_3_abc.png", 3),
(9, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'washing_machine_4_abc.png',"items/2498c61e-3911-4272-9013-5d21c1e165bc/images/washing_machine_4_abc.png", 3),
(10, uuid_v4(),'2022-04-06 05:46:03.528','2022-04-06 05:46:03.528',NULL,1,1,1,NULL,'washing_machine_5_abc.png',"items/2498c61e-3911-4272-9013-5d21c1e165bc/images/washing_machine_5_abc.png", 3)
;
`)
	item1Image1, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_1.jpeg", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item1Image2, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_2.jpg", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item2Image1, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_3.jpeg", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item2Image2, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_4.jpeg", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item2Image3, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/guitar_5.jpeg", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item3Image1, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/washing_machine_1.png", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item3Image2, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/washing_machine_2.png", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item3Image3, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/washing_machine_3.png", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item3Image4, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/washing_machine_4.png", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	item3Image5, err := os.Open(fmt.Sprintf("%s/actions/item/fixtures/washing_machine_5.png", os.Getenv("PROJECTDIR")))
	utils.PanicIf(err)
	defer item1Image1.Close()
	defer item1Image2.Close()
	defer item2Image1.Close()
	defer item2Image2.Close()
	defer item2Image3.Close()
	defer item3Image1.Close()
	defer item3Image2.Close()
	defer item3Image3.Close()
	defer item3Image4.Close()
	defer item3Image5.Close()

	fmt.Println("uploading item 1 images...")
	item1UploadFiles := []*os.File{item1Image1, item1Image2}
	item1 := repository.FindItemById(1)
	err = UploadItemImages(item1, item1UploadFiles)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not upload item %d images due to error: %s, exiting...", item1.ID, err.Error()))
		os.Exit(1)
	}
	fmt.Println("uploaded item 1 images...")
	fmt.Println("uploading item 2 images...")
	item2UploadFiles := []*os.File{item2Image1, item2Image2, item2Image3}
	item2 := repository.FindItemById(2)
	err = UploadItemImages(item2, item2UploadFiles)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not upload item %d images due to error: %s, exiting...", item2.ID, err.Error()))
		os.Exit(1)
	}
	fmt.Println("uploaded item 2 images...")
	fmt.Println("uploading item 3 images...")
	item3UploadFiles := []*os.File{item3Image1, item3Image2, item3Image3, item3Image4, item3Image5}
	item3 := repository.FindItemById(3)
	err = UploadItemImages(item3, item3UploadFiles)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not upload item %d images due to error: %s, exiting...", item3.ID, err.Error()))
		os.Exit(1)
	}

	fmt.Println("uploaded item 3 images...")
	fmt.Println("all content generated successfully.")
}

func UploadItemImages(item *models.Item, files []*os.File) error {
	token, err := services.GenerateNewJwtToken(item.UserCreated, services.TokenTypeAccess)
	if err != nil {
		fmt.Println("could not generate new token ", err.Error())
		return err
	}
	payload, contentType, err := client.MakeMultiPartWriterFromFiles("images", files...)
	resp, err := client.MakeRequest(
		fmt.Sprintf("%s/%d/images?removeExisting=true", itemsRoute, item.ID),
		"POST",
		map[string]string{},
		map[string]string{"Authorization": token, "Content-Type": contentType},
		time.Second*10,
		payload,
	)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("could not save item images ", err.Error())
		return err
	} else if resp.StatusCode != http.StatusCreated {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("could not parse response body of save images response, err:", err.Error())
			return err
		}
		fmt.Println("could not save item images err:", string(respBytes))
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(setupTestDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupTestDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupTestDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
