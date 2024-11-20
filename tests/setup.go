package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/dewciu/f1_api/pkg/config"
	"github.com/dewciu/f1_api/pkg/migrations"
	"github.com/dewciu/f1_api/pkg/seeding"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

const POSTGRES_PORT = "5432"

func SetupDB(tablesAffected []string) (*gorm.DB, tc.Container, context.Context) {
	config.CONFIG_PATH = "../app-config.yaml"
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

	migrations.Migrate(db)
	for _, table := range tablesAffected {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table))
	}
	seeding.Seed(db)

	return db, pg, ctx
}
