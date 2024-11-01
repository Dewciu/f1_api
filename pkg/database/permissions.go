package database

import (
	m "github.com/dewciu/f1_api/pkg/models"
)

func GetPermissionByIDQuery(id string) (m.Permission, error) {
	var permission m.Permission
	err := DB.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return m.Permission{}, err
	}
	return permission, nil
}
