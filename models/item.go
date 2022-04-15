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
	EmptyItemBidUserError = fmt.Errorf("placing item bid requires a user but was found empty")
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
	LastBidDate time.Time    `gorm:"name:last_bid_date;type:date;not null" json:"lastBidDate"`

	Bids []Bid
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

func (item Item) IsBidEligible() bool {

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
