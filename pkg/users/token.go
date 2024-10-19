package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/dewciu/f1_api/pkg/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateToken(user_id uuid.UUID) (string, error) {
	config, err := config.GetConfig()
	if err != nil {
		return "", err
	}
	token_hour_lifetime := config.Server.TokenHourLifetime

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_hour_lifetime)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Server.ApiSecret))
}

func ValidateToken(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		config, err := config.GetConfig()
		if err != nil {
			return 0, err
		}
		return []byte(config.Server.ApiSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return tokenString, nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	bearerToken := c.GetHeader("Authorization")
	fmt.Println(bearerToken)
	if bearerToken != "" {
		return bearerToken
	}
	return ""
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			config, err := config.GetConfig()
			if err != nil {
				return 0, err
			}
			return []byte(config.Server.ApiSecret), nil
		})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	user_id := claims["user_id"].(string)

	if user_id == "" {
		return "", errors.New("missing encrypted user ID in token")
	}

	return user_id, nil
}
