package routes

import (
	_ "github.com/dewciu/f1_api/docs"
	c "github.com/dewciu/f1_api/pkg/controllers"
	"github.com/dewciu/f1_api/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	UsersEndpoint = "/users"
	AuthEndpoint  = "/auth"
	LoginEndpoint = "/login"
)

var UsersPermissions = []models.Permission{}

func AddUsersRoutes(rg *gin.RouterGroup, db *gorm.DB, middlewareHandlers ...gin.HandlerFunc) {
	users := rg.Group(UsersEndpoint, middlewareHandlers...)
	c := c.NewUserController(db)
	{
		users.GET("/", c.GetAllUsers)
		users.POST("/", c.CreateUser)
		users.GET("/:id", c.GetUserByID)
		users.DELETE("/:id", c.DeleteUserByID)
		users.PUT("/:id", c.UpdateUser)
		users.GET("/:id"+PermissionsEndpoint, c.GetUserWithPermissions)
	}
}

func GetUserPermissions() []models.Permission {
	return []models.Permission{
		{
			Endpoint: UsersEndpoint,
			Method:   "GET",
		},
		{
			Endpoint: UsersEndpoint,
			Method:   "POST",
		},
		{
			Endpoint: UsersEndpoint,
			Method:   "DELETE",
		},
		{
			Endpoint: UsersEndpoint + "/:id",
			Method:   "GET",
		},
		{
			Endpoint: UsersEndpoint + "/:id",
			Method:   "DELETE",
		},
		{
			Endpoint: UsersEndpoint + "/:id",
			Method:   "PUT",
		},
		{
			Endpoint: UsersEndpoint + "/:id" + PermissionsEndpoint,
			Method:   "GET",
		},
	}
}

func AddAuthRoutes(rg *gin.RouterGroup, db *gorm.DB, handlers ...gin.HandlerFunc) {
	auth := rg.Group(AuthEndpoint, handlers...)
	c := c.NewUserController(db)
	{
		auth.POST(LoginEndpoint, c.Login)
	}
}

func GetAuthPermissions() []models.Permission {
	return []models.Permission{
		{
			Endpoint: LoginEndpoint,
			Method:   "POST",
		},
	}
}
