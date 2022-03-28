package models

type LoggedInUser struct {
	Identity
	CreatedMetaInfo
	UpdatedMetaInfo
	VersionMetaInfo

	FirstName		string	`json:"firstName"`
	LastName		string	`json:"lastName"`
	Email			string	`json:"email"`
	AccessToken		string	`json:"accessToken"`
	RefreshToken	string	`json:"refreshToken"`
}

func CreateLoggedInUserByUser(user User) LoggedInUser {
	return LoggedInUser{
		Identity: Identity{
			ID:   user.ID,
			Uuid: user.Uuid,
		},
		CreatedMetaInfo: CreatedMetaInfo{
			CreatedAt: user.CreatedAt,
			CreatedBy: user.CreatedBy,
		},
		UpdatedMetaInfo: UpdatedMetaInfo{
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.UpdatedBy,
		},
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
	}
}