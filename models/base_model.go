package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Identity struct {
	ID		uint	`gorm:"name:id;primarykey;AUTO_INCREMENT" json:"id"`
	Uuid	string	`gorm:"name:uuid;type:varchar(36);unique_index" json:"uuid"`
}

func (identity Identity) GetId() uint {
	return identity.ID
}

func (identity Identity) GetUuid() string {
	return identity.Uuid
}

func (identity *Identity) SetId(newId uint) *Identity {
	identity.ID = newId

	return identity
}

func (identity *Identity) SetUuid(newUuid string) *Identity {
	identity.Uuid = newUuid

	return identity
}

func (identity *Identity) AssignUuid() {
	if identity.Uuid == "" {
		identity.Uuid = uuid.NewString()
	}
}

func (identity *Identity) BeforeCreate(tx *gorm.DB) (err error) {
	identity.AssignUuid()

	return nil
}

type CreatedMetaInfo struct {
	CreatedAt		time.Time	`gorm:"name:created_at;type:datetime;autoCreateTime:milli" json:"createdAt"`
	CreatedBy		*uint		`gorm:"name:created_by;index"`
	CreatedByUser	*User		`gorm:"foreignKey:CreatedBy" json:"createdBy,omitempty"`
}

func (info CreatedMetaInfo) GetCreatedAt() time.Time {
	return info.CreatedAt
}

type UpdatedMetaInfo struct {
	UpdatedAt		time.Time	`gorm:"name:updated_at;type:datetime;autoUpdateTime:milli" json:"updatedAt"`
	UpdatedBy		*uint		`gorm:"name:updated_by;index"`
	UpdatedByUser	*User		`gorm:"foreignKey:UpdatedBy" json:"updatedBy,omitempty"`
}

type DeletedMetaInfo struct {
	DeletedAt		*time.Time	`gorm:"name:deleted_at;type:datetime" json:"deletedAt"`
	DeletedBy		*uint		`gorm:"name:deleted_by"`
	DeletedByUser	*User		`gorm:"foreignKey:DeletedBy" json:"deletedBy,omitempty"`
}

func (val DeletedMetaInfo) IsDeleted() bool {
	return val.DeletedAt != nil
}

type ActionedMetaInfo struct {
	ActionedAt	time.Time	`gorm:"name:actioned_at;type:datetime" json:"actionedAt,omitempty"`
	CycledAt	time.Time	`gorm:"name:cycled_at;type:datetime" json:"cycledAt,omitempty"`
}

type RequestedMetaInfo struct {
	RequestedAt time.Time `gorm:"name:requested_at;type:datetime" json:"requestedAt,omitempty"`
}

type VersionMetaInfo struct {
	Version int32 `gorm:"name:version;" json:"version"`
}