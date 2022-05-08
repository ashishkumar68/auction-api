package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	CategoryElectronicsInt = iota
	CategoryAppliancesInt
	CategoryHomeInt
	CategoryArtInt
)

var (
	EmptyItemBidUserError   = fmt.Errorf("placing item bid requires a user but was found empty")
	DetectedPastLastBidDate = fmt.Errorf("can not set last bid date in the past")
)

type Value float32
type ItemCategory uint8

type Item struct {
	IdentityAuditableModel

	Name        string       `gorm:"type:varchar(512)" json:"name"`
	Description string       `gorm:"type:varchar(1024)" json:"description"`
	Category    ItemCategory `gorm:"type:smallint" json:"category"`
	BrandName   string       `gorm:"type:varchar(1024)" json:"brandName"`
	MarketValue Value        `gorm:"type:float(16,4)" json:"marketValue"`
	LastBidDate time.Time    `gorm:"column:last_bid_date;type:date;not null" json:"lastBidDate"`
	OffBid      bool         `gorm:"column:off_bid;type:tinyint(1);not null;default:0" json:"isOffBid"`

	Bids       []Bid
	ItemImages []*ItemImage `json:"itemImages"`
}

func (Item) TableName() string {
	return "items"
}

func NewItemFromValues(
	name string,
	description string,
	category ItemCategory,
	brandName string,
	value Value,
	lastBidDate time.Time,
) *Item {

	return &Item{
		Name:        name,
		Description: description,
		Category:    category,
		BrandName:   brandName,
		MarketValue: value,
		LastBidDate: lastBidDate,
	}
}

func (item *Item) UpdateFromValues(
	name string,
	description string,
	category ItemCategory,
	brandName string,
	value Value,
	lastBidDate time.Time,
	actionUser *User,
) error {
	if name != "" && name != item.Name {
		item.Name = name
	}
	if description != "" && description != item.Description {
		item.Description = description
	}
	item.Category = category
	if brandName != "" && brandName != item.BrandName {
		item.BrandName = brandName
	}
	if value != 0 && value != item.MarketValue {
		item.MarketValue = value
	}
	if !lastBidDate.IsZero() && !lastBidDate.Equal(item.LastBidDate) {
		if lastBidDate.Before(time.Now()) {
			return DetectedPastLastBidDate
		}
		item.LastBidDate = lastBidDate
	}
	item.UserUpdatedBy = actionUser.ID

	return nil
}

func (item Item) IsOffBid() bool {
	return item.OffBid
}

func (item *Item) MarkOffBid() error {
	if !item.IsOffBid() {
		item.OffBid = true
	}

	return nil
}

func (item Item) IsOwner(user User) bool {
	return item.UserCreated.IsSameAs(user.BaseModel)
}

func (item Item) IsBidEligible() bool {

	if item.IsOffBid() {
		return false
	}
	if time.Now().Before(item.LastBidDate) {
		return true
	}

	return false
}

func GetAvailableItemCategories() []int {
	return []int{CategoryElectronicsInt, CategoryAppliancesInt, CategoryHomeInt, CategoryArtInt}
}

type Bid struct {
	IdentityAuditableModel

	ItemId uint  `gorm:"column:item_id;index" json:"-"`
	Item   *Item `gorm:"foreignKey:ItemId" json:"item"`
	Value  Value `gorm:"type:float(16,4)" json:"bidValue"`
}

func (Bid) TableName() string {
	return "bids"
}

func NewBidFromValues(
	item *Item,
	value Value,
	bidBy *User,
) *Bid {

	newBid := &Bid{
		ItemId: item.ID,
		Item:   nil,
		Value:  value,
	}
	if bidBy != nil {
		newBid.UserCreatedBy = bidBy.ID
		newBid.UserUpdatedBy = bidBy.ID
	}

	return newBid
}

func (bid *Bid) AfterCreate(db *gorm.DB) error {
	if bid.UserCreatedBy == 0 {
		return EmptyItemBidUserError
	}

	return nil
}

func (item *Item) AddImage(img *ItemImage) *Item {
	item.ItemImages = append(item.ItemImages, img)

	return item
}

func (item *Item) RemoveImage(img *ItemImage) *Item {
	deleteImgIndex := -1
	for index, itemImg := range item.ItemImages {
		if itemImg.IsSameAs(img.BaseModel) {
			deleteImgIndex = index
			break
		}
	}
	if deleteImgIndex == -1 {
		return item
	}
	item.ItemImages = append(item.ItemImages[:deleteImgIndex], item.ItemImages[deleteImgIndex+1:]...)

	return item
}
