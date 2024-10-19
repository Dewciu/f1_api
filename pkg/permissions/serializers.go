package permissions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionResponse struct {
	ID       uuid.UUID `json:"id"`
	Endpoint string    `json:"endpoint"`
	Method   string    `json:"method"`
} // @name PermissionResponse

type PermissionSerializer struct {
	C *gin.Context
	PermissionModel
}

func (s *PermissionSerializer) Response() PermissionResponse {

	response := PermissionResponse{
		ID:       s.ID,
		Endpoint: s.Endpoint,
		Method:   s.Method,
	}

	return response
}

type PermissionsSerializer struct {
	C           *gin.Context
	Permissions []PermissionModel
}

func (s *PermissionsSerializer) Response() []PermissionResponse {
	var response []PermissionResponse
	for _, permission := range s.Permissions {
		serializer := PermissionSerializer{s.C, permission}
		response = append(response, serializer.Response())
	}

	return response
}

type PermissionGroupResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Permissions []PermissionModel
}

type PermissionGroupSerializer struct {
	C *gin.Context
	PermissionGroupModel
	Permissions []PermissionModel
}

func (s *PermissionGroupSerializer) Response() PermissionGroupResponse {
	response := PermissionGroupResponse{
		ID:          s.ID,
		Name:        s.Name,
		Permissions: s.Permissions,
	}

	return response
}
