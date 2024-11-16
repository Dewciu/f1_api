package controllers

import (
	"errors"
	"fmt"
	"net/http"

	_ "github.com/dewciu/f1_api/docs"
	"github.com/dewciu/f1_api/pkg/common"
	d "github.com/dewciu/f1_api/pkg/database"
	m "github.com/dewciu/f1_api/pkg/models"
	s "github.com/dewciu/f1_api/pkg/serializers"
	v "github.com/dewciu/f1_api/pkg/validators"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	DB       *gorm.DB
	userRepo *d.UserRepository
}

func NewUserController(db *gorm.DB) *UserController {
	userRepo := d.NewUserRepository(db)
	return &UserController{DB: db, userRepo: userRepo}
}

//TODO: Improve error handling -> more descriptive error messages

// @BasePath /api/v1

// Login godoc
// @Summary Retrieve JWT API token
// @Description Retrieve JWT API token, when given valid username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param Credentials body LoginValidator true "Login Credentials"
// @Success 200 {object} TokenResponse "Returns JWT token"
// @Router /auth/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var validator v.LoginValidator

	if err := c.ShouldBindJSON(&validator); err != nil {
		c.JSON(http.StatusBadRequest, common.NewError("login", err))
		return
	}
	u := m.User{Username: validator.Username, Password: validator.Password}

	token, err := uc.userRepo.LoginCheck(u)

	serializer := s.TokenSerializer{C: c, Token: token}

	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewError("login", errors.New("invalid credentials")))
		return
	}

	c.JSON(http.StatusOK, serializer.Response())
}

// @BasePath /api/v1

// GetAllUsers godoc
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
func (uc *UserController) GetAllUsers(c *gin.Context) {
	var users []m.User
	var err error

	if len(c.Request.URL.Query()) == 0 {
		users, err = uc.userRepo.GetAllUsersQuery()
	} else {
		users, err = uc.userRepo.GetUsersByFilterQuery(c)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewError("server", err))
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, common.NewError("users", err))
		return
	}

	serializer := s.UsersSerializer{C: c, Users: users}

	c.JSON(http.StatusOK, serializer.Response())
}

// @BasePath /api/v1

// CreateUser godoc
// @Summary Create User
// @Description Creates single user in database
// @Tags users
// @Accept json
// @Produce json
// @Param User body UserCreateModelValidator true "User Object"
// @Security ApiKeyAuth
// @Success 200 {object} UserResponse "Returns Created User"
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	validator := v.UserCreateModelValidator{}
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user := validator.User

	err := uc.userRepo.CreateUserQuery(user)

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

	serializer := s.UserSerializer{C: c, User: user}

	c.JSON(http.StatusCreated, serializer.Response())
}

// GetUserByID godoc
// @Summary Get User by ID
// @Description Retrieves a user from the database by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "Returns the user"
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := uc.userRepo.GetUserByIdQuery(id)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, common.NewError("user", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewError("user", err))
		return
	}
	serializer := s.UserSerializer{C: c, User: user}
	c.JSON(http.StatusOK, serializer.Response())
}

// DeleteUserByID godoc
// @Summary Delete User by ID
// @Description Deletes a user from the database by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")

	err := uc.userRepo.DeleteUserByIdQuery(id)
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

// UpdateUser godoc
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
func (uc *UserController) UpdateUser(c *gin.Context) {
	//TODO Finish update user controller
	id := c.Param("id")
	validator := v.UserUpdateModelValidator{}
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := uc.userRepo.UpdateUserByIdQuery(id, validator)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, common.NewError("user", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewError("user", err))
		return
	}
	serializer := s.UserSerializer{C: c, User: user}
	c.JSON(http.StatusOK, serializer.Response())
}

// GetUserWithPermissions godoc
// @Summary Retrieve Permissions for the user by ID
// @Description Retrieves permission list for specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} []PermissionResponse "Returns user's permissions"
// @Router /users/{id}/permissions [get]
func (uc *UserController) GetUserWithPermissions(c *gin.Context) {
	id := c.Param("id")

	permissions, err := uc.userRepo.GetPermissionsForUserIDQuery(id)
	if err != nil || len(permissions) == 0 {
		c.JSON(http.StatusNotFound, common.NewError("permissions", errors.New("permissions not found")))
		return
	}

	serializer := s.PermissionsSerializer{C: c, Permissions: permissions}
	c.JSON(http.StatusOK, serializer.Response())
}
