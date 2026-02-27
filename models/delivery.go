package models

type Delivery struct {
	Id       int64  `json:"id,omitempty"`
	OrderUId string `json:"order_uid,omitempty"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Zip      string `json:"zip" validate:"required"`
	City     string `json:"city" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Region   string `json:"region" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
