package models

import (
	"fmt"
	"gorm.io/gorm"
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

	Bids []Bid
}

func NewItemFromValues(
	name string,
	description string,
	category ItemCategory,
	brandName string,
	value Value,
) *Item {

	return &Item{
		Name:        name,
		Description: description,
		Category:    category,
		BrandName:   brandName,
		MarketValue: value,
	}
}

func GetAvailableItemCategories() []int {
	return []int{CategoryElectronicsInt, CategoryAppliancesInt, CategoryHomeInt, CategoryArtInt}
}

type Bid struct {
	IdentityAuditableModel

	ItemId *uint `gorm:"column:item_id;index" json:"-"`
	Item   *Item `gorm:"foreignKey:ItemId" json:"item"`
	Value  Value `gorm:"type:float(16,4)" json:"bidValue"`
}

func NewBidFromValues(
	item *Item,
	value Value,
) *Bid {

	return &Bid{
		ItemId: &item.ID,
		Value:  value,
	}
}

func (bid *Bid) AfterCreate(db *gorm.DB) error {
	if bid.UserCreatedBy == nil {
		return EmptyItemBidUserError
	}

	return nil
}