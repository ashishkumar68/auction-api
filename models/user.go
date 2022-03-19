package models

type User struct {
	Identity
	CreatedMetaInfo
	UpdatedMetaInfo
	DeletedMetaInfo
	VersionMetaInfo

	FirstName		string	`gorm:"type:varchar(512)" json:"firstName"`
	LastName		string	`gorm:"type:varchar(512)" json:"lastName"`
	Email			string	`gorm:"type:varchar(512)" json:"email"`
	plainPassword	string
	Password		string  `gorm:"type:varchar(2048)" json:"-"`
	IsActive		bool	`gorm:"type:tinyint(1)" json:"isActive"`
}

