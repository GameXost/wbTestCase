package models

type Item struct {
	Id          int64  `json:"id,omitempty"`
	OrderUId    string `json:"order_uid,omitempty"`
	ChrtId      int64  `json:"chrt_id" validate:"gt=0"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int64  `json:"price" validate:"gte=0"`
	RID         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int64  `json:"sale" validate:"gte=0"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price" validate:"gte=0"`
	NmId        int64  `json:"nm_id" validate:"gt=0"`
	Brand       string `json:"brand"`
	Status      int64  `json:"status" validate:"gte=0"`
}
