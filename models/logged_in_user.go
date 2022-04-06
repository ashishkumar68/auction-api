package models

type LoggedInUser struct {
	User

	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func CreateLoggedInUserByUser(user User) LoggedInUser {
	return LoggedInUser{
		User: user,
	}
}
