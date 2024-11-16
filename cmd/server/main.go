package main

import (
	"fmt"
	"os"

	"github.com/dewciu/f1_api/pkg/config"
	"github.com/dewciu/f1_api/pkg/database"
	"github.com/dewciu/f1_api/pkg/migrations"
	"github.com/dewciu/f1_api/pkg/routes"
	"github.com/dewciu/f1_api/pkg/seeding"
	"github.com/sirupsen/logrus"

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
	if len(os.Args) < 2 {
		config.CONFIG_PATH = "app-config.yaml"
	} else {
		config.CONFIG_PATH = os.Args[1]
	}

	conf, err := config.GetConfig()

	if err != nil {
		logrus.Panicf("Failed to get configuration: %v", err)
	}

	DB, err := database.Connect(conf)

	if err != nil {
		msg := fmt.Sprintf("Failed to connect to DB: %v", err)
		panic(msg)
	}
	router := routes.SetupRouter(DB)

	defer database.Disconnect(DB)

	if err = migrations.Migrate(DB); err != nil {
		msg := fmt.Sprintf("Failed to migrate DB: %v", err)
		panic(msg)
	}

	if err = seeding.Seed(DB); err != nil {
		msg := fmt.Sprintf("Failed to seed DB: %v", err)
		panic(msg)
	}

	hostname := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	router.Run(hostname)
}
