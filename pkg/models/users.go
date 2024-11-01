package models

import (
	a "github.com/dewciu/f1_api/pkg/auth"
	"gorm.io/gorm"
)

// TODO Add permissions to endpoints for the user
type User struct {
	Model
	Username    string        `gorm:"unique;not null; type:varchar(255)" json:"username"`
	Email       string        `gorm:"unique;not null; type:varchar(255)" json:"email"`
	Password    string        `gorm:"not null" json:"password"`
	Permissions []*Permission `gorm:"many2many:user_permissions;"`
} //@name User

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password == "" {
		return nil
	}

	hash, err := a.GeneratePassword(u.Password)

	if err != nil {
		return err
	}

	u.Password = hash

	return u.Model.BeforeCreate(tx)
}
