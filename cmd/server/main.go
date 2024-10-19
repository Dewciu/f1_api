package main

import (
	"fmt"

	"github.com/dewciu/f1_api/pkg/addresses"
	"github.com/dewciu/f1_api/pkg/common"
	"github.com/dewciu/f1_api/pkg/config"
	"github.com/dewciu/f1_api/pkg/middleware"
	perm "github.com/dewciu/f1_api/pkg/permissions"
	"github.com/dewciu/f1_api/pkg/users"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"

	files "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	_ "github.com/dewciu/f1_api/docs"
)

// @title F1 API
// @version 1.0
// @description This is an API for F1 application
// @termsOfService http://swagger.io/terms/

// @contact.name Kacper Król
// @contact.email kacperkrol99@icloud.com

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @Schemes http https
func main() {
	conf, err := config.GetConfig()

	if err != nil {
		logrus.Panicf("Failed to get configuration: %v", err)
	}

	router := SetupRouter()

	if err = common.Connect(conf); err != nil {
		msg := fmt.Sprintf("Failed to connect to DB: %v", err)
		panic(msg)
	}

	defer common.Disconnect()

	if err = Migrate(); err != nil {
		msg := fmt.Sprintf("Failed to migrate DB: %v", err)
		panic(msg)
	}

	if err = Seed(); err != nil {
		msg := fmt.Sprintf("Failed to seed DB: %v", err)
		panic(msg)
	}

	hostname := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	router.Run(hostname)
}

func Migrate() error {
	if err := common.DB.AutoMigrate(
		&users.UserModel{},
		&addresses.AddressModel{},
		&perm.PermissionModel{},
		&perm.PermissionGroupModel{},
	); err != nil {
		return err
	}

	return nil
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	addSwaggerRoutes(v1)
	users.AddAuthRoutes(v1)
	users.AddUsersRoutes(v1, middleware.JwtAuthMiddleware(), middleware.PermAuthMiddleware(v1.BasePath()))
	return r
}

// TODO: Improve seeding
func Seed() error {

	adminName := "admin"

	if common.DB.First(&users.UserModel{}, "username = ?", adminName).RowsAffected <= 0 {
		err := users.CreateUserQuery(users.UserModel{
			Username: adminName,
			Password: "admin",
		})
		if err != nil {
			return err
		}
	}

	var permissions [][]perm.PermissionModel = [][]perm.PermissionModel{
		users.GetUserPermissions(),
		users.GetAuthPermissions(),
	}

	var batchPermissions []perm.PermissionModel

	for _, permission := range permissions {
		batchPermissions = append(batchPermissions, permission...)
	}

	if err := common.DB.Create(&batchPermissions).Error; err != nil {
		err := err.(*pgconn.PgError)

		if err.Code != "23505" {
			return err
		}
	}

	return nil
}

func addSwaggerRoutes(rg *gin.RouterGroup) {
	swag := rg.Group("/swagger")
	swag.GET("/*any", swagger.WrapHandler(files.Handler))
}
