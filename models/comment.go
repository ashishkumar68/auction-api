package models

type ItemComment struct {
	IdentityAuditableModel

	Description	string	`gorm:"description;not null" json:"description"`
	ItemId		uint	`gorm:"column:item_id;index" json:"-"`
	Item		*Item	`gorm:"foreignKey:ItemId" json:"item"`
}

func (ItemComment) TableName() string {
	return "item_comments"
}