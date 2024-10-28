package database

import (
	"fmt"
	"net/url"

	"github.com/dewciu/f1_api/pkg/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(config *config.Config) error {
	dsn := url.URL{
		User:     url.UserPassword(config.Database.User, config.Database.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", config.Database.Host, config.Database.Port),
		Path:     config.Database.Name,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn.String()}), &gorm.Config{})

	if err != nil {
		msg := fmt.Sprintf("Failed to connect to database: %v", err)
		logrus.Errorf(msg)
		return err
	}

	return nil
}

func Disconnect() error {
	db, err := DB.DB()

	if err != nil {
		return err
	}

	return db.Close()
}
