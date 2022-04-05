package models

const (
	CategoryElectronicsInt	=	iota
	CategoryAppliancesInt
	CategoryHomeInt
	CategoryArtInt
)

type Value float32
type ItemCategory uint8

type Item struct {
	IdentityAuditableModel

	Name		string			`gorm:"type:varchar(512)" json:"name"`
	Description	string			`gorm:"type:varchar(1024)" json:"description"`
	Category	ItemCategory	`gorm:"type:smallint" json:"category"`
	BrandName	string			`gorm:"type:varchar(1024)" json:"brandName"`
	MarketValue	Value			`gorm:"type:float(16,4)" json:"marketValue"`
}

func NewItemFromValues(
	name string,
	description string,
	category ItemCategory,
	brandName string,
	value Value) *Item {

	return &Item{
		Name: name,
		Description: description,
		Category: category,
		BrandName: brandName,
		MarketValue: value,
	}
}

func GetAvailableItemCategories() []int {
	return []int{CategoryElectronicsInt, CategoryAppliancesInt, CategoryHomeInt, CategoryArtInt}
}