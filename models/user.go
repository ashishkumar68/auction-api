package models

import (
	"gorm.io/gorm"
)

type User struct {
	BaseModel

	FirstName     string `gorm:"type:varchar(512)" json:"firstName"`
	LastName      string `gorm:"type:varchar(512)" json:"lastName"`
	Email         string `gorm:"type:varchar(512)" json:"email"`
	PlainPassword string `gorm:"-" json:"-"`
	Password      string `gorm:"type:varchar(2048)" json:"-"`
	IsActive      bool   `gorm:"type:tinyint(1)" json:"isActive"`

	CreatedBy     *uint `gorm:"index" json:"-"`
	CreatedByUser *User `gorm:"foreignKey:CreatedBy" json:"createdBy"`

	UpdatedBy     *uint `gorm:"index" json:"-"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy" json:"updatedBy"`

	DeletedBy     *uint `gorm:"index" json:"-"`
	DeletedByUser *User `gorm:"foreignKey:DeletedBy" json:"deletedBy"`
}

func NewUserFromValues(
	firstName string,
	lastName string,
	email string,
	plainPass string) User {

	return User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		PlainPassword: plainPass,
		// TODO: Currently setting user active by default but it needs to be email verified to get activated.
		IsActive: true,
	}
}

func (user User) GetLoginId() string {
	return user.Email
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	user.AssignUuid()
	actionUser := GetActionUser(db)
	if actionUser.IsZero() {
		return nil
	}
	user.CreatedBy = &actionUser.ID
	user.UpdatedBy = &actionUser.ID

	return nil
}
