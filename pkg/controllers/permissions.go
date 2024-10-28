package controllers

import (
	"errors"
	"net/http"

	"github.com/dewciu/f1_api/pkg/common"
	d "github.com/dewciu/f1_api/pkg/database"
	s "github.com/dewciu/f1_api/pkg/serializers"
	"github.com/gin-gonic/gin"
)

func GetPermissionByIDController(c *gin.Context) {
	id := c.Param("id")

	permission, err := d.GetPermissionByIDQuery(id)
	// TODO: Make better error handling, ex. return 500 if something else goes wrong
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("permissions", errors.New("permission not found")))
		return
	}

	serializer := s.PermissionSerializer{C: c, Permission: permission}
	c.JSON(http.StatusOK, serializer.Response())
}
