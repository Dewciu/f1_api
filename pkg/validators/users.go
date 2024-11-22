package validators

import (
	"github.com/dewciu/f1_api/pkg/common"
	m "github.com/dewciu/f1_api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/google/uuid"
)

type UserCreateModelValidator struct {
	Username string `form:"username" json:"username" binding:"required,alphanum,min=4,max=255"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=255"`
	User     m.User `json:"-"`
} // @name UserCreateModelValidator

var g = galidator.New()

// TODO: Add custom error messages
func (s *UserCreateModelValidator) Bind(c *gin.Context) interface{} {
	err := common.Bind(c, s)
	customizer := g.Validator(UserCreateModelValidator{})

	if err != nil {
		return customizer.DecryptErrors(err)
	}

	s.User.ID = uuid.New()
	s.User.Username = s.Username
	s.User.Email = s.Email
	s.User.Password = s.Password

	return nil
}

type UserUpdateModelValidator struct {
	Username string `json:"username" binding:"omitempty,alphanum,min=4,max=255" `
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8,max=255"`
} // @name UserUpdateModelValidator

func (s *UserUpdateModelValidator) Bind(c *gin.Context) interface{} {

	customizer := g.Validator(UserUpdateModelValidator{})

	err := c.ShouldBindJSON(s)
	if err != nil {
		return customizer.DecryptErrors(err)
	}

	return nil
}

type LoginValidator struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name LoginValidator

func (s *LoginValidator) Bind(c *gin.Context) interface{} {
	customizer := g.Validator(LoginValidator{})
	err := common.Bind(c, s)
	if err != nil {
		return customizer.DecryptErrors(err)
	}

	return nil
}
