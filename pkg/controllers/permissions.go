package controllers

import (
	"errors"
	"net/http"

	"github.com/dewciu/f1_api/pkg/common"
	d "github.com/dewciu/f1_api/pkg/database"
	s "github.com/dewciu/f1_api/pkg/serializers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionController struct {
	permRepo *d.PermissionRepository
}

func NewPermissionController(db *gorm.DB) *PermissionController {
	permRepo := d.NewPermissionRepository(db)
	return &PermissionController{permRepo: permRepo}
}

func (pc *PermissionController) GetPermissionByIDController(c *gin.Context) {
	id := c.Param("id")

	permission, err := pc.permRepo.GetPermissionByIDQuery(id)
	// TODO: Make better error handling, ex. return 500 if something else goes wrong
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("permissions", errors.New("permission not found")))
		return
	}

	serializer := s.PermissionSerializer{C: c, Permission: permission}
	c.JSON(http.StatusOK, serializer.Response())
}
