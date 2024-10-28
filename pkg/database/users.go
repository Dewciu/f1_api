package database

import (
	"errors"

	a "github.com/dewciu/f1_api/pkg/auth"
	"github.com/dewciu/f1_api/pkg/common"
	"github.com/dewciu/f1_api/pkg/models"
	m "github.com/dewciu/f1_api/pkg/models"
	v "github.com/dewciu/f1_api/pkg/validators"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetAllUsersQuery() ([]m.User, error) {
	var users []m.User
	err := common.DB.Find(&users).Error
	return users, err
}
func CreateUserQuery(user m.User) error {
	r := common.DB.Create(&user)

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
	query := common.DB

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
	err := common.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func DeleteUserByIdQuery(id string) error {
	err := common.DB.Where("id = ?", id).Delete(&m.User{}).Error
	return err
}

func UpdateUserByIdQuery(id string, userToUpdate v.UserUpdateModelValidator) (m.User, error) {
	var user m.User

	if userToUpdate.Password != "" {
		hash, err := a.GeneratePassword(userToUpdate.Password)
		if err != nil {
			return m.User{}, err
		}
		userToUpdate.Password = hash
	}

	if err := common.DB.Model(&user).Where("id = ?", id).Updates(userToUpdate).First(&user).Error; err != nil {
		return m.User{}, err
	}

	return user, nil
}

func GetPermissionsForUserIDQuery(id string) ([]models.Permission, error) {

	var user m.User

	err := common.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		return []models.Permission{}, err
	}

	var permissions []models.Permission

	err = common.DB.Model(&user).Association("Permissions").Find(&permissions)

	if err != nil {
		return []models.Permission{}, err
	}

	if len(permissions) == 0 {
		// TODO: Better errors
		return permissions, errors.New("no permissions found for user")
	}

	return permissions, nil
}
