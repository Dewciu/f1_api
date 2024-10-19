package addresses

import (
	"github.com/dewciu/f1_api/pkg/common"
	"github.com/jackc/pgx/v5/pgconn"
)

func CreateAddressQuery(address AddressModel) error {
	r := common.DB.Create(&address)
	if r.Error != nil {
		err := r.Error.(*pgconn.PgError)

		if err.Code == "23505" {
			column := common.GetColumnFromUniqueErrorDetails(err.Detail)
			return &common.AlreadyExistsError{Column: column}
		}
	}

	return r.Error
}
