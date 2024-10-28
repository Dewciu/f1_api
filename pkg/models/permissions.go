package models

type Permission struct {
	Model
	Endpoint string `json:"endpoint" gorm:"not null;index:,unique,composite:idx_endpoint_method"`
	Method   string `json:"method" gorm:"not null;index:,unique,composite:idx_endpoint_method"`
} //@name Permission

type PermissionGroup struct {
	Model
	Name        string        `json:"name"`
	Permissions []*Permission `gorm:"many2many:permission_group_permissions;"`
} //@name PermissionGroup
