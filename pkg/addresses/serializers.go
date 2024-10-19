package addresses

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressResponse struct {
	ID     uuid.UUID `json:"-"`
	Street string    `json:"street"`
	City   string    `json:"city"`
	State  string    `json:"state"`
} //@name Address

type AddressSerializer struct {
	C *gin.Context
	AddressModel
}

func (s *AddressSerializer) Response() AddressResponse {
	response := AddressResponse{
		ID:     s.ID,
		Street: s.Street,
		City:   s.City,
		State:  s.State,
	}

	return response
}
