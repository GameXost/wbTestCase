package models

type Delivery struct {
	Id       int64
	OrderUId string
	Name     string
	Phone    string
	Zip      string
	City     string
	Address  string
	Region   string
	Email    string
}
