package main

import (
	"fmt"

	"github.com/dewciu/f1_api/pkg/config"
	"github.com/dewciu/f1_api/pkg/database"
	"github.com/dewciu/f1_api/pkg/middleware"
	"github.com/dewciu/f1_api/pkg/models"
	"github.com/dewciu/f1_api/pkg/routes"
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

	if err = database.Connect(conf); err != nil {
		msg := fmt.Sprintf("Failed to connect to DB: %v", err)
		panic(msg)
	}

	defer database.Disconnect()

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
	if err := database.DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Permission{},
		&models.PermissionGroup{},
	); err != nil {
		return err
	}

	return nil
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Handler()
	v1 := r.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(database.DB)
	addSwaggerRoutes(v1)
	routes.AddAuthRoutes(
		v1,
		database.DB,
	)
	routes.AddUsersRoutes(
		v1,
		database.DB,
		authMiddleware.CheckJWT(),
		authMiddleware.CheckPermissions(v1.BasePath()),
	)
	return r
}

// TODO: Improve seeding
func Seed() error {

	adminName := "admin"
	repo := database.NewUserRepository(database.DB)

	if database.DB.First(&models.User{}, "username = ?", adminName).RowsAffected <= 0 {
		err := repo.CreateUserQuery(models.User{
			Username: adminName,
			Password: "admin",
		})
		if err != nil {
			return err
		}
	}

	var permissions [][]models.Permission = [][]models.Permission{
		routes.GetUserPermissions(),
		routes.GetAuthPermissions(),
	}

	var batchPermissions []models.Permission

	for _, permission := range permissions {
		batchPermissions = append(batchPermissions, permission...)
	}

	if err := database.DB.Create(&batchPermissions).Error; err != nil {
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

//TODO: Consider changing the folder structure (ex. controllers together)
