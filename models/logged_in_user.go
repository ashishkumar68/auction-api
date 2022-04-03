package models

type LoggedInUser struct {
	BaseModel

	FirstName		string	`json:"firstName"`
	LastName		string	`json:"lastName"`
	Email			string	`json:"email"`
	AccessToken		string	`json:"accessToken"`
	RefreshToken	string	`json:"refreshToken"`
}

func CreateLoggedInUserByUser(user User) LoggedInUser {
	return LoggedInUser{
		BaseModel: BaseModel{
			ID:   user.ID,
			Uuid: user.Uuid,
			CreatedAt: user.CreatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.UpdatedBy,
		},
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
	}
}