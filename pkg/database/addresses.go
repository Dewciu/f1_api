package database

import (
	"github.com/dewciu/f1_api/pkg/common"
	m "github.com/dewciu/f1_api/pkg/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func CreateAddressQuery(address m.Address) error {
	r := DB.Create(&address)
	if r.Error != nil {
		err := r.Error.(*pgconn.PgError)

		if err.Code == "23505" {
			column := common.GetColumnFromUniqueErrorDetails(err.Detail)
			return &common.AlreadyExistsError{Column: column}
		}
	}

	return r.Error
}
