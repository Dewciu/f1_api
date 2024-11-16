package routes

import (
	"github.com/dewciu/f1_api/pkg/middleware"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(DB *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Handler()
	v1 := r.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(DB)
	addSwaggerRoutes(v1)
	AddAuthRoutes(
		v1,
		DB,
	)
	AddUsersRoutes(
		v1,
		DB,
		authMiddleware.CheckJWT(),
		authMiddleware.CheckPermissions(v1.BasePath()),
	)
	return r
}

func addSwaggerRoutes(rg *gin.RouterGroup) {
	swag := rg.Group("/swagger")
	swag.GET("/*any", swagger.WrapHandler(files.Handler))
}
