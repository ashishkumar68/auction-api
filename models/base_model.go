package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID		uint	`gorm:"name:id;primarykey;AUTO_INCREMENT" json:"id"`
	Uuid	string	`gorm:"name:uuid;type:varchar(36);unique_index" json:"uuid"`

	CreatedAt		time.Time	`gorm:"name:created_at;type:datetime;autoCreateTime:milli" json:"createdAt"`
	CreatedBy		*uint		`gorm:"name:created_by;index" json:"-"`
	CreatedByUser	*User		`gorm:"foreignKey:CreatedBy" json:"createdBy,omitempty"`

	UpdatedAt		time.Time	`gorm:"name:updated_at;type:datetime;autoUpdateTime:milli" json:"updatedAt"`
	UpdatedBy		*uint		`gorm:"name:updated_by;index" json:"-"`
	UpdatedByUser	*User		`gorm:"foreignKey:UpdatedBy" json:"updatedBy,omitempty"`

	DeletedAt		*time.Time	`gorm:"name:deleted_at;type:datetime" json:"deletedAt"`
	DeletedBy		*uint		`gorm:"name:deleted_by" json:"-"`
	DeletedByUser	*User		`gorm:"foreignKey:DeletedBy" json:"deletedBy,omitempty"`

	Version			int32		`gorm:"name:version;default:1" json:"version"`
}

func (base BaseModel) GetId() uint {
	return base.ID
}

func (base BaseModel) GetUuid() string {
	return base.Uuid
}

func (base BaseModel) IsZero() bool {
	return 0 == base.ID
}

func (base *BaseModel) SetId(newId uint) *BaseModel {
	base.ID = newId

	return base
}

func (base *BaseModel) SetUuid(newUuid string) *BaseModel {
	base.Uuid = newUuid

	return base
}

func (base *BaseModel) AssignUuid() {
	if base.Uuid == "" {
		base.Uuid = uuid.NewString()
	}
}

func (base *BaseModel) BeforeCreate(db *gorm.DB) (err error) {
	base.AssignUuid()
	actionUser := GetActionUser(db)
	if actionUser.IsZero() {
		return nil
	}
	base.CreatedBy = &actionUser.ID
	base.UpdatedBy = &actionUser.ID

	return nil
}

func (base BaseModel) GetCreatedAt() time.Time {
	return base.CreatedAt
}

func (base *BaseModel) BeforeSave(db *gorm.DB) error {
	actionUser := GetActionUser(db)
	if actionUser.IsZero() {
		return nil
	}
	base.UpdatedBy = &actionUser.ID

	return nil
}

func (base BaseModel) IsDeleted() bool {
	return base.DeletedAt != nil
}

type ActionedMetaInfo struct {
	ActionedAt	time.Time	`gorm:"name:actioned_at;type:datetime" json:"actionedAt,omitempty"`
	CycledAt	time.Time	`gorm:"name:cycled_at;type:datetime" json:"cycledAt,omitempty"`
}

type RequestedMetaInfo struct {
	RequestedAt time.Time `gorm:"name:requested_at;type:datetime" json:"requestedAt,omitempty"`
}

func GetActionUser(db *gorm.DB) User {
	var actionUser User
	if user, ok := db.Get("actionUser"); ok {
		actionUser = user.(User)
	}

	return actionUser
}