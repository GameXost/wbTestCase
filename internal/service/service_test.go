package service

import (
	"context"
	"errors"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/internal/generator"
	"github.com/GameXost/wbTestCase/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidateOrder(t *testing.T) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOrder(&tt.order)
			if !errors.Is(got, tt.want) {
				t.Errorf("feq, we want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestGetOrderCacheHit(t *testing.T) {
	repo := NewMockOrderRepo(t)
	cache := NewMockOrderCache(t)
	serv := NewService(repo, cache)

	ord := &models.Order{OrderUId: "test1"}

	cache.EXPECT().Get("test1").Return(ord, true)
	res, err := serv.GetOrder(context.Background(), "test1")
	assert.NoError(t, err)
	assert.Equal(t, ord, res)
}

func TestGetOrderCacheMiss(t *testing.T) {
	repo := NewMockOrderRepo(t)
	cache := NewMockOrderCache(t)
	serv := NewService(repo, cache)

	ord := &models.Order{OrderUId: "test2"}

	cache.EXPECT().Get("test2").Return(nil, false)
	repo.EXPECT().GetFullOrderOnId(mock.Anything, "test2").Return(ord, nil)
	cache.EXPECT().Set(ord)
	res, err := serv.GetOrder(context.Background(), "test2")
	assert.NoError(t, err)
	assert.Equal(t, ord, res)
}

func TestGetOrderNothingFound(t *testing.T) {
	repo := NewMockOrderRepo(t)
	cache := NewMockOrderCache(t)
	serv := NewService(repo, cache)

	cache.EXPECT().Get("test3").Return(nil, false)
	repo.EXPECT().GetFullOrderOnId(mock.Anything, "test3").Return(nil, errHandle.ErrNotFound)

	res, err := serv.GetOrder(context.Background(), "test3")
	assert.ErrorIs(t, err, errHandle.ErrNotFound)
	assert.Nil(t, res)
}

var cases = []struct {
	name  string
	order models.Order
	want  error
}{
	{
		"correct",
		*generator.ValidOrder("test"),
		nil,
	},
	{
		name: "no Items",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items = []models.Item{}
			return ord
		}(),
		want: errHandle.ErrItemsEmpty,
	},
	{
		name: "no TrackNumber",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.TrackNumber = ""
			return ord
		}(),
		want: errHandle.ErrTrackNumberMissing,
	},
	{
		name: "no OrderUId",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.OrderUId = ""
			return ord
		}(),
		want: errHandle.ErrOrderUIDMissing,
	},
	{
		name: "no Entry",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Entry = ""
			return ord
		}(),
		want: errHandle.ErrEntryMissing,
	},
	{
		name: "no Locale",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Locale = ""
			return ord
		}(),
		want: errHandle.ErrLocaleMissing,
	},
	{
		name: "no CustomerID",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.CustomerId = ""
			return ord
		}(),
		want: errHandle.ErrCustomerIDMissing,
	},
	{
		name: "no DeliveryService",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.DeliveryService = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryServiceMissing,
	},
	{
		name: "no Shardkey",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Shardkey = ""
			return ord
		}(),
		want: errHandle.ErrShardkeyMissing,
	},
	{
		name: "invalid SmId",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.SmId = 0
			return ord
		}(),
		want: errHandle.ErrInvalidSmID,
	},
	{
		name: "no Delivery Name",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Name = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryNameMissing,
	},
	{
		name: "no Delivery Phone",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Phone = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryPhoneMissing,
	},
	{
		name: "no Delivery ZIP",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Zip = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryZIPMissing,
	},
	{
		name: "no Delivery City",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.City = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryCityMissing,
	},
	{
		name: "no Delivery Address",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Address = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryAddressMissing,
	},
	{
		name: "no Delivery Region",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Region = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryRegionMissing,
	},
	{
		name: "no Delivery Email",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Email = ""
			return ord
		}(),
		want: errHandle.ErrDeliveryEmailMissing,
	},
	{
		name: "no Payment Transaction",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.Transaction = ""
			return ord
		}(),
		want: errHandle.ErrPaymentTransactionMissing,
	},
	{
		name: "no Payment Currency",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.Currency = ""
			return ord
		}(),
		want: errHandle.ErrPaymentCurrencyMissing,
	},
	{
		name: "invalid Payment Amount",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.Amount = 0
			return ord
		}(),
		want: errHandle.ErrPaymentAmountInvalid,
	},
	{
		name: "invalid Payment DeliveryCost",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.DeliveryCost = -1
			return ord
		}(),
		want: errHandle.ErrPaymentDeliveryInvalid,
	},
	{
		name: "invalid Payment GoodsTotal",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.GoodsTotal = 0
			return ord
		}(),
		want: errHandle.ErrPaymentGoodsTotalInvalid,
	},
	{
		name: "item name missing",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Name = ""
			return ord
		}(),
		want: errHandle.ErrItemNameMissing,
	},
	{
		name: "item invalid ChrtId",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].ChrtId = 0
			return ord
		}(),
		want: errHandle.ErrItemChrtMissing,
	},
	{
		name: "item invalid Price",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Price = -1
			return ord
		}(),
		want: errHandle.ErrItemPriceInvalid,
	},
	{
		name: "item invalid Sale",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Sale = -1
			return ord
		}(),
		want: errHandle.ErrItemSaleInvalid,
	},
	{
		name: "item invalid TotalPrice",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].TotalPrice = -1
			return ord
		}(),
		want: errHandle.ErrItemTotalPriceInvalid,
	},
	{
		name: "item invalid Status",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Status = -1
			return ord
		}(),
		want: errHandle.ErrStatusCodeInvalid,
	},
}
