package models

import (
	"fmt"
	"github.com/ashishkumar68/auction-api/utils"
	"log"
	"mime/multipart"
)

type ItemImage struct {
	IdentityAuditableModel

	Path             string                `gorm:"column:path;not null" json:"-"`
	ItemId           uint                  `gorm:"column:item_id;index" json:"-"`
	Item             *Item                 `gorm:"foreignKey:ItemId" json:"item"`
	MultiPartImgFile *multipart.FileHeader `gorm:"-" json:"-"`
	Name             string                `gorm:"-" json:"name"`
}

func (ItemImage) TableName() string {
	return "item_images"
}

func NewItemImageFromMultipartFile(item *Item, file *multipart.FileHeader, actionUser *User) (*ItemImage, error) {
	newFileName, err := utils.GetRenamedFileName(file.Filename)
	if err != nil {
		log.Printf(fmt.Sprintf("could not get renamed file name for: %s", file.Filename))
		return nil, err
	}
	baseImgPath := fmt.Sprintf("item/%s/images", item.Uuid)
	itemImage := &ItemImage{
		Path:             fmt.Sprintf("%s/%s", baseImgPath, newFileName),
		ItemId:           item.ID,
		Item:             item,
		MultiPartImgFile: file,
		Name:             newFileName,
	}
	itemImage.SetCreatedBy(actionUser)

	return itemImage, nil
}
