package seeding

import (
	"github.com/dewciu/f1_api/pkg/database"
	"github.com/dewciu/f1_api/pkg/models"
	"github.com/dewciu/f1_api/pkg/routes"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// TODO: Improve seeding
func Seed(DB *gorm.DB) error {

	adminName := "admin"
	userRepo := database.NewUserRepository(DB)

	if DB.First(&models.User{}, "username = ?", adminName).RowsAffected <= 0 {
		var permissions [][]models.Permission = [][]models.Permission{
			routes.GetUserPermissions(),
			routes.GetAuthPermissions(),
		}

		var batchPermissions []models.Permission

		for _, permission := range permissions {
			batchPermissions = append(batchPermissions, permission...)
		}

		if err := DB.Create(&batchPermissions).Error; err != nil {
			err := err.(*pgconn.PgError)

			if err.Code != "23505" {
				return err
			}
		}
		err := userRepo.CreateUserQuery(models.User{
			Username:    adminName,
			Password:    "admin",
			Permissions: batchPermissions,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
