package forms

type RegisterNewUserForm struct {

	FirstName string `json:"firstName" binding:"required,min=3,max=200"`
	LastName  string `json:"lastName" binding:"required,min=3,max=200"`
	Email     string `json:"email" binding:"required,min=3,max=200,email"`
	Password  string `json:"password" binding:"required,min=8,max=80"`
}