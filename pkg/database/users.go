package database

import (
	"errors"

	"github.com/dewciu/f1_api/pkg/auth"
	"github.com/dewciu/f1_api/pkg/common"
	m "github.com/dewciu/f1_api/pkg/models"
	v "github.com/dewciu/f1_api/pkg/validators"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsersQuery() ([]m.User, error) {
	var users []m.User
	err := DB.Find(&users).Error
	return users, err
}
func CreateUserQuery(user m.User) error {
	r := DB.Create(&user)

	if r.Error != nil {
		err := r.Error.(*pgconn.PgError)

		if err.Code == "23505" {
			column := common.GetColumnFromUniqueErrorDetails(err.Detail)
			return &common.AlreadyExistsError{Column: column}
		}

		return err
	}

	return nil
}

func GetUsersByFilterQuery(c *gin.Context) ([]m.User, error) {
	var users []m.User
	query := DB

	if username := c.Query("username"); username != "" {
		query = query.Where("username = ?", username)
	}
	if email := c.Query("email"); email != "" {
		query = query.Where("email = ?", email)
	}
	if id := c.Query("id"); id != "" {
		query = query.Where("id = ?", id)
	}

	if err := query.Find(&users).Error; err != nil {
		return []m.User{}, err
	}

	return users, nil
}

func GetUserByIdQuery(id string) (m.User, error) {
	var user m.User
	err := DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func DeleteUserByIdQuery(id string) error {
	err := DB.Where("id = ?", id).Delete(&m.User{}).Error
	return err
}

func UpdateUserByIdQuery(id string, userToUpdate v.UserUpdateModelValidator) (m.User, error) {
	var user m.User

	if userToUpdate.Password != "" {
		hash, err := auth.GeneratePassword(userToUpdate.Password)
		if err != nil {
			return m.User{}, err
		}
		userToUpdate.Password = hash
	}

	if err := DB.Model(&user).Where("id = ?", id).Updates(userToUpdate).First(&user).Error; err != nil {
		return m.User{}, err
	}

	return user, nil
}

func GetPermissionsForUserIDQuery(id string) ([]m.Permission, error) {

	var user m.User

	err := DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		return []m.Permission{}, err
	}

	var permissions []m.Permission

	err = DB.Model(&user).Association("Permissions").Find(&permissions)

	if err != nil {
		return []m.Permission{}, err
	}

	if len(permissions) == 0 {
		// TODO: Better errors
		return permissions, errors.New("no permissions found for user")
	}

	return permissions, nil
}

func LoginCheck(u m.User) (string, error) {

	var user m.User

	result := DB.Model(&user).Where("username = ?", u.Username).First(&user)
	err := result.Error

	if err != nil {
		return "", err
	}

	err = auth.VerifyPassword(u.Password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := auth.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	err = DB.Model(&user).Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return "", err
	}

	return token, nil
}
