package forms

import "github.com/ashishkumar68/auction-api/models"

type AuditableForm struct {
	ActionUser *models.User `binding:"required"`
}
