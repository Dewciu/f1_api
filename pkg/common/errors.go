package common

import (
	"fmt"
)

type ValidationError struct {
	Errors map[string]interface{} `json:"errors"`
}

type AlreadyExistsError struct {
	Column string
}

func (s AlreadyExistsError) Error() string {
	return fmt.Sprintf("record already exist, conflict column: %s", s.Column)
}

func NewValidationError(err error) error {
	res := ValidationError{}
	res.Errors = make(map[string]interface{})
	// errs := err.(validator.ValidationErrors)
	// for _, v := range errs {
	// 	if v.Param() != "" {
	// 		res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
	// 	} else {
	// 		res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
	// 	}
	// }
	return err
}

func NewError(key string, err error) ValidationError {
	res := ValidationError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}
