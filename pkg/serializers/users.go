package serializers

import (
	m "github.com/dewciu/f1_api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
} //@name UserResponse

type UserSerializer struct {
	C *gin.Context
	m.User
}

func (s *UserSerializer) Response() UserResponse {

	response := UserResponse{
		ID:       s.ID,
		Username: s.Username,
		Email:    s.Email,
	}

	return response
}

type UsersSerializer struct {
	C     *gin.Context
	Users []m.User
}

func (s *UsersSerializer) Response() []UserResponse {
	var response []UserResponse
	for _, user := range s.Users {
		serializer := UserSerializer{s.C, user}
		response = append(response, serializer.Response())
	}

	return response
}

type TokenResponse struct {
	Token string `json:"token"`
} //@name TokenResponse

type TokenSerializer struct {
	C     *gin.Context
	Token string
}

func (s *TokenSerializer) Response() TokenResponse {
	response := TokenResponse{
		Token: s.Token,
	}

	return response
}
