package models

type Payment struct {
	OrderId      string `json:"order_id,omitempty"`
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id" validate:"required"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int64  `json:"amount" validate:"gt=0"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int64  `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int64  `json:"goods_total" validate:"gt=0"`
	CustomFee    int64  `json:"custom_fee"`
}
