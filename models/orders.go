package models

import "time"

type Order struct {
	OrderUId          string    `json:"order_uid,omitempty"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int64     `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Delivery          Delivery  `json:"delivery"`
}
