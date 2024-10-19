package permissions

import (
	"errors"
	"net/http"

	"github.com/dewciu/f1_api/pkg/common"
	"github.com/gin-gonic/gin"
)

func GetPermissionByIDController(c *gin.Context) {
	id := c.Param("id")

	permission, err := GetPermissionByIDQuery(id)
	// TODO: Make better error handling, ex. return 500 if something else goes wrong
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("permissions", errors.New("permission not found")))
		return
	}

	serializer := PermissionSerializer{C: c, PermissionModel: permission}
	c.JSON(http.StatusOK, serializer.Response())
}
