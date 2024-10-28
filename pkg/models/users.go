package models

import (
	a "github.com/dewciu/f1_api/pkg/auth"
	"github.com/dewciu/f1_api/pkg/common"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) LoginCheck() (string, error) {

	var user User

	result := common.DB.Model(&user).Where("username = ?", u.Username).First(&user)
	err := result.Error

	if err != nil {
		return "", err
	}

	err = a.VerifyPassword(u.Password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := a.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	err = common.DB.Model(&user).Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return "", err
	}

	return token, nil
}
