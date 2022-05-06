package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IdentityAuditableModel struct {
	BaseModel
	AuditModel
}

type BaseModel struct {
	ID   uint   `gorm:"name:id;primaryKey;AUTO_INCREMENT" json:"id"`
	Uuid string `gorm:"name:uuid;type:varchar(36);unique_index" json:"uuid"`

	CreatedAt time.Time  `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoCreateTime:milli" json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`

	Version int32 `gorm:"name:version;default:1" json:"version"`
}

type AuditModel struct {
	UserCreatedBy uint  `gorm:"column:created_by;index;not null" json:"-"`
	UserCreated   *User `gorm:"foreignKey:UserCreatedBy" json:"createdBy,omitempty"`

	UserUpdatedBy uint  `gorm:"column:updated_by;index;not null" json:"-"`
	UserUpdated   *User `gorm:"foreignKey:UserUpdatedBy" json:"updatedBy,omitempty"`

	UserDeletedBy *uint `gorm:"column:deleted_by;index" json:"-"`
	UserDeleted   *User `gorm:"foreignKey:UserDeletedBy" json:"deletedBy,omitempty"`
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

func (auditModel *AuditModel) SetCreatedBy(actionUser *User) *AuditModel {
	auditModel.UserCreatedBy = actionUser.ID
	auditModel.UserUpdatedBy = actionUser.ID

	return auditModel
}

func (auditModel *AuditModel) SetUpdatedBy(actionUser *User) *AuditModel {
	auditModel.UserUpdatedBy = actionUser.ID

	return auditModel
}

func (auditModel *AuditModel) SetDeletedBy(actionUser *User) *AuditModel {
	auditModel.UserDeletedBy = &actionUser.ID

	return auditModel
}

func (base *BaseModel) AssignUuid() {
	if base.Uuid == "" {
		base.Uuid = uuid.NewString()
	}
}

func (base *IdentityAuditableModel) BeforeCreate(db *gorm.DB) (err error) {
	base.AssignUuid()

	return nil
}

func (base BaseModel) GetCreatedAt() time.Time {
	return base.CreatedAt
}

func (base BaseModel) IsDeleted() bool {
	return base.DeletedAt != nil
}

func (base BaseModel) IsSameAs(model BaseModel) bool {
	return base.ID == model.ID && base.Uuid == model.Uuid
}

type ActionedMetaInfo struct {
	ActionedAt time.Time `gorm:"name:actioned_at;type:datetime" json:"actionedAt,omitempty"`
	CycledAt   time.Time `gorm:"name:cycled_at;type:datetime" json:"cycledAt,omitempty"`
}

type RequestedMetaInfo struct {
	RequestedAt time.Time `gorm:"name:requested_at;type:datetime" json:"requestedAt,omitempty"`
}
