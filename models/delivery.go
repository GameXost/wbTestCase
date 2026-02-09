package models

type Delivery struct {
	Id       int64  `json:"id,omitempty"`
	OrderUId string `json:"order_uid,omitempty"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}
