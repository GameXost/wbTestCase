package models

type Payment struct {
	OrderId      string
	Transaction  string
	RequestId    string
	Currency     string
	Provider     string
	Amount       int64
	PaymentDt    int64
	Bank         string
	DeliveryCost int64
	GoodsTotal   int64
	CustomFee    int64
}
