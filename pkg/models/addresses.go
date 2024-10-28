package models

// TODO: Relate Address to UserModel etc...
type Address struct {
	Model
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	UserID      string `json:"user_id"`
	User        User
} //@name Address
