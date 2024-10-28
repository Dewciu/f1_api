package validators

import (
	"github.com/dewciu/f1_api/pkg/common"
	"github.com/dewciu/f1_api/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressModelValidator struct {
	Address struct {
		Street      string `form:"street" json:"street" binding:"required,max=255"`
		City        string `form:"city" json:"city" binding:"required,max=255"`
		HouseNumber string `form:"house_number" json:"house_number" binding:"required,max=5"`
		State       string `form:"state" json:"state" binding:"required,max=255"`
	} `json:"address"`
	addressModel models.Address `json:"-"`
}

func (s *AddressModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}

	s.addressModel.ID = uuid.New()
	s.addressModel.Street = s.Address.Street
	s.addressModel.City = s.Address.City
	s.addressModel.HouseNumber = s.Address.HouseNumber
	s.addressModel.State = s.Address.State

	return nil
}
