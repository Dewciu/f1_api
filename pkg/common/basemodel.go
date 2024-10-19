package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

func (b *BaseModel) BeforeDelete(tx *gorm.DB) error {
	b.DeletedAt = time.Now()
	return nil
}
