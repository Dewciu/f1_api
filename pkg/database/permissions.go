package database

import (
	m "github.com/dewciu/f1_api/pkg/models"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{DB: db}
}

func (repo *PermissionRepository) GetPermissionByIDQuery(id string) (m.Permission, error) {
	var permission m.Permission
	err := repo.DB.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return m.Permission{}, err
	}
	return permission, nil
}
