package tests

import (
	"context"
	"fmt"
	"time"

	m "github.com/dewciu/f1_api/pkg/models"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

const POSTGRES_PORT = "5432"

var models = []interface{}{
	m.User{},
	m.Address{},
	m.Permission{},
	m.PermissionGroup{},
}

func setupDB(tablesAffected []string) {
	ctx := context.Background()
	postgres_user := "testuser"
	postgres_password := "testpassword"
	postgres_db := "testdb"

	req := tc.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{POSTGRES_PORT + "/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     postgres_user,
			"POSTGRES_PASSWORD": postgres_password,
			"POSTGRES_DB":       postgres_db,
		},
		WaitingFor: wait.ForListeningPort(POSTGRES_PORT + "/tcp").WithStartupTimeout(60 * time.Second),
	}

	pg, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	host, _ := pg.Host(ctx)
	mappedPort, _ := pg.MappedPort(ctx, POSTGRES_PORT)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, mappedPort.Port(), postgres_user, postgres_password, postgres_db)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(models...)

	for _, table := range tablesAffected {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table))
	}
}
