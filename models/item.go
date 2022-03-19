package models

type Item struct {
	Identity
	CreatedMetaInfo
	UpdatedMetaInfo
	DeletedMetaInfo

	Name		string	`gorm:"type:varchar(512)" json:"name"`
	Description	string	`gorm:"type:varchar(1024)" json:"description"`
}
