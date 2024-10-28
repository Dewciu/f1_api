package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dewciu/f1_api/pkg/auth"
	"github.com/dewciu/f1_api/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.ValidateToken(c)
		if err != nil {
			// TODO: Make logging everywhere and more informative
			logrus.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		req_user_id, err := auth.ExtractUserIDFromToken(token)

		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("req_user_id", req_user_id)
	}
}

func PermAuthMiddleware(basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.ReplaceAll(c.FullPath(), basePath, "")
		fmt.Println(path)
		req_user_id, ok := c.Get("req_user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		permissions, err := database.GetPermissionsForUserIDQuery(req_user_id.(string))

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		fmt.Println(permissions)

		for _, perm := range permissions {
			fmt.Println(perm.Endpoint)
			if perm.Endpoint == path {
				return
			}

			if perm.Endpoint+"/" == path {
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
