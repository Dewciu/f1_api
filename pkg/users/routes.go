package users

import (
	_ "github.com/dewciu/f1_api/docs"
	perm "github.com/dewciu/f1_api/pkg/permissions"
	"github.com/gin-gonic/gin"
)

const (
	UsersEndpoint       = "/users"
	AuthEndpoint        = "/auth"
	LoginEndpoint       = "/login"
	PermissionsEndpoint = "/permissions"
)

var UsersPermissions = []perm.PermissionModel{}

func AddUsersRoutes(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	users := rg.Group(UsersEndpoint, handlers...)
	{
		users.GET("/", GetAllUsersController)
		users.POST("/", CreateUserController)
		users.GET("/:id", GetUserByIDController)
		users.DELETE("/:id", DeleteUserByIDController)
		users.PUT("/:id", UpdateUserController)
		users.GET("/:id"+PermissionsEndpoint, GetUserWithPermissionsController)
	}
}

func GetUserPermissions() []perm.PermissionModel {
	return []perm.PermissionModel{
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

func AddAuthRoutes(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	auth := rg.Group(AuthEndpoint, handlers...)
	{
		auth.POST(LoginEndpoint, LoginController)
	}
}

func GetAuthPermissions() []perm.PermissionModel {
	return []perm.PermissionModel{
		{
			Endpoint: LoginEndpoint,
			Method:   "POST",
		},
	}
}
