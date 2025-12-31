package models

import "time"

type Order struct {
	OrderUId          string
	TrackNumber       string
	Entry             string
	Locale            string
	InternalSignature string
	CustomerId        string
	DeliveryService   string
	Shardkey          string
	SmId              int64
	DateCreated       time.Time
	OofShard          string
	Payment           Payment
	Items             []Item
	Delivery          Delivery
}
