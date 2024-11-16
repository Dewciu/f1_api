package migrations

import (
	"github.com/dewciu/f1_api/pkg/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) error {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Permission{},
		&models.PermissionGroup{},
	); err != nil {
		return err
	}

	return nil
}
