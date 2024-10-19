package permissions

import (
	"github.com/dewciu/f1_api/pkg/common"
)

func GetPermissionByIDQuery(id string) (PermissionModel, error) {
	var permission PermissionModel
	err := common.DB.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return PermissionModel{}, err
	}
	return permission, nil
}
