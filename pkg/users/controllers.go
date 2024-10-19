package users

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/dewciu/f1_api/docs"
	"github.com/dewciu/f1_api/pkg/common"
	perm "github.com/dewciu/f1_api/pkg/permissions"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

//TODO: Improve error handling -> more descriptive error messages

// @BasePath /api/v1

// LoginController godoc
// @Summary Retrieve JWT API token
// @Description Retrieve JWT API token, when given valid username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param Credentials body LoginValidator true "Login Credentials"
// @Success 200 {object} TokenResponse "Returns JWT token"
// @Router /auth/login [post]
func LoginController(c *gin.Context) {
	var validator LoginValidator

	if err := c.ShouldBindJSON(&validator); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("login", err))
		return
	}

	u := UserModel{Username: validator.Username, Password: validator.Password}

	token, err := u.LoginCheck()

	serializer := TokenSerializer{C: c, Token: token}

	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("login", errors.New("invalid credentials")))
		return
	}

	c.JSON(http.StatusOK, serializer.Response())
}

// @BasePath /api/v1

// GetAllUsersController godoc
// @Summary Get Users
// @Description Retrieves all users from the database, with optional filters
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string false "User's E-mail"
// @Param username query string false "User's username"
// @Param id query string false "User's ID"
// @Success 200 {array} UserResponse "Returns list of users"
// @Router /users [get]
func GetAllUsersController(c *gin.Context) {
	var users []UserModel
	var err error

	if len(c.Request.URL.Query()) == 0 {
		fmt.Println(c.Request.URL.Query())
		users, err = GetAllUsersQuery()
	} else {
		users, err = GetUsersByFilterQuery(c)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("server", err))
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, common.NewError("users", err))
		return
	}

	serializer := UsersSerializer{C: c, Users: users}

	c.JSON(http.StatusOK, serializer.Response())
}

// @BasePath /api/v1

// CreateUserController godoc
// @Summary Create User
// @Description Creates single user in database
// @Tags users
// @Accept json
// @Produce json
// @Param User body UserCreateModelValidator true "User Object"
// @Security ApiKeyAuth
// @Success 200 {object} UserResponse "Returns Created User"
// @Router /users [post]
func CreateUserController(c *gin.Context) {
	validator := UserCreateModelValidator{}
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := CreateUserQuery(validator.UserModel)

	if err != nil {
		fmt.Println(err)
		var er *common.AlreadyExistsError
		if errors.As(err, &er) {
			err := err.(*common.AlreadyExistsError)
			c.JSON(http.StatusConflict, common.NewError("user", err))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewError("database", err))
		return
	}

	serializer := UserSerializer{C: c, UserModel: validator.UserModel}

	c.JSON(http.StatusCreated, serializer.Response())
}

// GetUserByIDController godoc
// @Summary Get User by ID
// @Description Retrieves a user from the database by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "Returns the user"
// @Router /users/{id} [get]
func GetUserByIDController(c *gin.Context) {
	id := c.Param("id")
	user, err := GetUserByIdQuery(id)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, common.NewError("user", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewError("user", err))
		return
	}
	serializer := UserSerializer{C: c, UserModel: user}
	c.JSON(http.StatusOK, serializer.Response())
}

// DeleteUserByIDController godoc
// @Summary Delete User by ID
// @Description Deletes a user from the database by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Router /users/{id} [delete]
func DeleteUserByIDController(c *gin.Context) {
	id := c.Param("id")
	err := DeleteUserByIdQuery(id)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, common.NewError("user", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewError("user", err))
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateUserController godoc
// @Summary Update User by ID
// @Description Updates a user in the database by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param User body UserUpdateModelValidator true "User Object fields to update"
// @Success 200 {object} UserResponse "Returns the updated user"
// @Router /users/{id} [put]
func UpdateUserController(c *gin.Context) {
	//TODO Finish update user controller
	id := c.Param("id")
	validator := UserUpdateModelValidator{}
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	user, err := UpdateUserByIdQuery(id, validator)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, common.NewError("user", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewError("user", err))
		return
	}
	serializer := UserSerializer{C: c, UserModel: user}
	c.JSON(http.StatusOK, serializer.Response())
}

// TODO Make this controlleer better with serializer and figure out what to do with permissions directory to avoid import cycle

// GetUserWithPermissionsController godoc
// @Summary Retrieve Permissions for the user by ID
// @Description Retrieves permission list for specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} []PermissionResponse "Returns user's permissions"
// @Router /users/{id}/permissions [get]
func GetUserWithPermissionsController(c *gin.Context) {
	id := c.Param("id")

	permissions, err := GetPermissionsForUserIDQuery(id)
	if err != nil || len(permissions) == 0 {
		c.JSON(http.StatusNotFound, common.NewError("permissions", errors.New("permissions not found")))
		return
	}

	serializer := perm.PermissionsSerializer{C: c, Permissions: permissions}
	c.JSON(http.StatusOK, serializer.Response())
}

// TODO: Fix errors! Post password: "testtest" email: "test" username: "test"
