package addresses

import (
	"github.com/dewciu/f1_api/pkg/common"
	"github.com/dewciu/f1_api/pkg/users"
)

// TODO: Relate AddressModel to UserModel etc...
type AddressModel struct {
	common.BaseModel
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	UserID      string `json:"user_id"`
	User        users.UserModel
}

func (s *AddressModel) TableName() string {
	return "addresses"
}
