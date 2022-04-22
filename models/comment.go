package models

import "fmt"

var (
	ItemCommentNotFound = fmt.Errorf("item comment was not found")
)

type ItemComment struct {
	IdentityAuditableModel

	Description string `gorm:"description;not null" json:"description"`
	ItemId      uint   `gorm:"column:item_id;index" json:"-"`
	Item        *Item  `gorm:"foreignKey:ItemId" json:"item"`
}

func (ItemComment) TableName() string {
	return "item_comments"
}

func NewItemComment(description string, item *Item, actionUser *User) *ItemComment {
	newItemComment := &ItemComment{
		Description: description,
		ItemId:      item.ID,
	}

	newItemComment.UserCreatedBy = actionUser.ID
	newItemComment.UserUpdatedBy = actionUser.ID

	return newItemComment
}

func (comment ItemComment) IsAuthor(user User) bool {

	if comment.UserCreated.IsSameAs(user.BaseModel) {
		return true
	}
	return false
}
