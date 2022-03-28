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
	PlainPassword	string	`gorm:"-" json:"-"`
	Password		string  `gorm:"type:varchar(2048)" json:"-"`
	IsActive		bool	`gorm:"type:tinyint(1)" json:"isActive"`
}

func NewUserFromValues(
	firstName string,
	lastName string,
	email string,
	plainPass string) *User {

	return &User{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		PlainPassword: plainPass,
		// TODO: Currently setting user active by default but it needs to be email verified to get activated.
		IsActive: true,
	}
}

func (user User) IsZero() bool {
	return user == User{}
}

func (user User) GetLoginId() string {
	return user.Email
}