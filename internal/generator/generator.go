package generator

import (
	"github.com/GameXost/wbTestCase/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"time"
)

func ValidOrder(orderUID string) *models.Order {
	return &models.Order{
		OrderUId:          orderUID,
		TrackNumber:       gofakeit.Noun(),
		Entry:             gofakeit.Noun(),
		Locale:            gofakeit.Language(),
		InternalSignature: gofakeit.Noun(),
		CustomerId:        gofakeit.ID(),
		DeliveryService:   gofakeit.Dessert(),
		Shardkey:          gofakeit.HackerNoun(),
		SmId:              int64(gofakeit.Number(1, 100_000_000)),
		DateCreated:       time.Time{}.UTC(),
		OofShard:          gofakeit.Noun(),
		Payment: models.Payment{
			OrderId:      orderUID,
			Transaction:  gofakeit.UUID(),
			RequestId:    gofakeit.Noun(),
			Currency:     gofakeit.Currency().Short,
			Provider:     gofakeit.BeerName(),
			Amount:       int64(gofakeit.Number(1, 100_000_000)),
			PaymentDt:    int64(gofakeit.Number(1, 100_000_000)),
			Bank:         gofakeit.BankName(),
			DeliveryCost: int64(gofakeit.Number(1, 100_000_000)),
			GoodsTotal:   int64(gofakeit.Number(1, 100_000_000)),
			CustomFee:    int64(gofakeit.Number(1, 100_000_000)),
		},
		Items: []models.Item{
			{
				Id:          0,
				OrderUId:    orderUID,
				ChrtId:      int64(gofakeit.Number(1, 100_000_000)),
				TrackNumber: gofakeit.Noun(),
				Price:       int64(gofakeit.Number(1, 100_000_000)),
				RID:         gofakeit.Noun(),
				Name:        gofakeit.Name(),
				Sale:        int64(gofakeit.Number(0, 100_000_000)),
				Size:        gofakeit.RandomString([]string{"m", "l", "s"}),
				TotalPrice:  int64(gofakeit.Number(1, 100_000_000)),
				NmId:        int64(gofakeit.Number(1, 100_000_000)),
				Brand:       gofakeit.BeerName(),
				Status:      int64(gofakeit.Number(100, 600)),
			},
		},
		Delivery: models.Delivery{
			Id:       0,
			OrderUId: orderUID,
			Name:     gofakeit.Name(),
			Phone:    gofakeit.Phone(),
			Zip:      gofakeit.Zip(),
			City:     gofakeit.City(),
			Address:  gofakeit.Address().Address,
			Region:   gofakeit.Country(),
			Email:    gofakeit.Email(),
		},
	}
}

func InvalidOrder(orderUID string) *models.Order {
	order := ValidOrder(orderUID)
	// слайс функций, из которого рандомно выбирается элемент, который портит модель заказа
	mutations := []func(*models.Order){
		func(o *models.Order) { o.OrderUId = "" },
		func(o *models.Order) { o.TrackNumber = "" },
		func(o *models.Order) { o.Entry = "" },
		func(o *models.Order) { o.Locale = "" },
		func(o *models.Order) { o.CustomerId = "" },
		func(o *models.Order) { o.DeliveryService = "" },
		func(o *models.Order) { o.Shardkey = "" },
		func(o *models.Order) { o.SmId = 0 },
		func(o *models.Order) { o.Items = []models.Item{} },
		func(o *models.Order) { o.Payment.Amount = 0 },
		func(o *models.Order) { o.Payment.Transaction = "" },
		func(o *models.Order) { o.Payment.Currency = "" },
		func(o *models.Order) { o.Payment.GoodsTotal = 0 },
		func(o *models.Order) { o.Payment.DeliveryCost = -1 },
		func(o *models.Order) { o.Items[0].Price = -1 },
		func(o *models.Order) { o.Items[0].TotalPrice = -1 },
		func(o *models.Order) { o.Items[0].Sale = -1 },
		func(o *models.Order) { o.Items[0].Status = -1 },
		func(o *models.Order) { o.Items[0].Name = "" },
		func(o *models.Order) { o.Items[0].ChrtId = 0 },
		func(o *models.Order) { o.Delivery.Name = "" },
		func(o *models.Order) { o.Delivery.Phone = "" },
		func(o *models.Order) { o.Delivery.Email = "" },
		func(o *models.Order) { o.Delivery.City = "" },
		func(o *models.Order) { o.Delivery.Address = "" },
		func(o *models.Order) { o.Delivery.Zip = "" },
		func(o *models.Order) { o.Delivery.Region = "" },
	}
	ind := gofakeit.Number(0, len(mutations)-1)
	if len(order.Items) == 0 && ind >= 14 && ind <= 19 {
		ind = 0
	}
	mutations[ind](order)

	return order
}
