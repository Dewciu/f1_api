package permissions

import (
	"github.com/dewciu/f1_api/pkg/common"
)

type PermissionModel struct {
	common.BaseModel
	Endpoint string `json:"endpoint" gorm:"not null;index:,unique,composite:idx_endpoint_method"`
	Method   string `json:"method" gorm:"not null;index:,unique,composite:idx_endpoint_method"`
}

func (u *PermissionModel) TableName() string {
	return "permissions"
}

type PermissionGroupModel struct {
	common.BaseModel
	Name        string             `json:"name"`
	Permissions []*PermissionModel `gorm:"many2many:permission_group_permissions;"`
} //@name PermissionGroup

func (u *PermissionGroupModel) TableName() string {
	return "permission_groups"
}
