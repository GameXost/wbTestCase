package service

import (
	"context"
	"github.com/GameXost/wbTestCase/internal/apperror"
	"github.com/GameXost/wbTestCase/internal/generator"
	"github.com/GameXost/wbTestCase/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidateOrder(t *testing.T) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateOrder(&tt.order)
			if got == nil {
				if tt.want != nil {
					t.Errorf("feq, we want %v, got %v", tt.want, got)
				}
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
	repo.EXPECT().GetFullOrderOnId(mock.Anything, "test3").Return(nil, apperror.ErrNotFound)

	res, err := serv.GetOrder(context.Background(), "test3")
	assert.ErrorIs(t, err, apperror.ErrNotFound)
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
		want: apperror.ErrItemsEmpty,
	},
	{
		name: "no TrackNumber",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.TrackNumber = ""
			return ord
		}(),
		want: apperror.ErrTrackNumberMissing,
	},
	{
		name: "no OrderUId",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.OrderUId = ""
			return ord
		}(),
		want: apperror.ErrOrderUIDMissing,
	},
	{
		name: "no Entry",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Entry = ""
			return ord
		}(),
		want: apperror.ErrEntryMissing,
	},
	{
		name: "no Locale",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Locale = ""
			return ord
		}(),
		want: apperror.ErrLocaleMissing,
	},
	{
		name: "no CustomerID",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.CustomerId = ""
			return ord
		}(),
		want: apperror.ErrCustomerIDMissing,
	},
	{
		name: "no DeliveryService",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.DeliveryService = ""
			return ord
		}(),
		want: apperror.ErrDeliveryServiceMissing,
	},
	{
		name: "no Shardkey",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Shardkey = ""
			return ord
		}(),
		want: apperror.ErrShardkeyMissing,
	},
	{
		name: "invalid SmId",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.SmId = 0
			return ord
		}(),
		want: apperror.ErrInvalidSmID,
	},
	{
		name: "no Delivery Name",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Name = ""
			return ord
		}(),
		want: apperror.ErrDeliveryNameMissing,
	},
	{
		name: "no Delivery Phone",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Phone = ""
			return ord
		}(),
		want: apperror.ErrDeliveryPhoneMissing,
	},
	{
		name: "no Delivery ZIP",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Zip = ""
			return ord
		}(),
		want: apperror.ErrDeliveryZIPMissing,
	},
	{
		name: "no Delivery City",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.City = ""
			return ord
		}(),
		want: apperror.ErrDeliveryCityMissing,
	},
	{
		name: "no Delivery Address",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Address = ""
			return ord
		}(),
		want: apperror.ErrDeliveryAddressMissing,
	},
	{
		name: "no Delivery Region",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Region = ""
			return ord
		}(),
		want: apperror.ErrDeliveryRegionMissing,
	},
	{
		name: "no Delivery Email",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Delivery.Email = ""
			return ord
		}(),
		want: apperror.ErrDeliveryEmailMissing,
	},
	{
		name: "no Payment Transaction",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.Transaction = ""
			return ord
		}(),
		want: apperror.ErrPaymentTransactionMissing,
	},
	{
		name: "no Payment Currency",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.Currency = ""
			return ord
		}(),
		want: apperror.ErrPaymentCurrencyMissing,
	},
	{
		name: "invalid Payment Amount",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.Amount = 0
			return ord
		}(),
		want: apperror.ErrPaymentAmountInvalid,
	},
	{
		name: "invalid Payment DeliveryCost",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.DeliveryCost = -1
			return ord
		}(),
		want: apperror.ErrPaymentDeliveryInvalid,
	},
	{
		name: "invalid Payment GoodsTotal",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Payment.GoodsTotal = 0
			return ord
		}(),
		want: apperror.ErrPaymentGoodsTotalInvalid,
	},
	{
		name: "item name missing",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Name = ""
			return ord
		}(),
		want: apperror.ErrItemNameMissing,
	},
	{
		name: "item invalid ChrtId",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].ChrtId = 0
			return ord
		}(),
		want: apperror.ErrItemChrtMissing,
	},
	{
		name: "item invalid Price",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Price = -1
			return ord
		}(),
		want: apperror.ErrItemPriceInvalid,
	},
	{
		name: "item invalid Sale",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Sale = -1
			return ord
		}(),
		want: apperror.ErrItemSaleInvalid,
	},
	{
		name: "item invalid TotalPrice",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].TotalPrice = -1
			return ord
		}(),
		want: apperror.ErrItemTotalPriceInvalid,
	},
	{
		name: "item invalid Status",
		order: func() models.Order {
			ord := *generator.ValidOrder("test")
			ord.Items[0].Status = -1
			return ord
		}(),
		want: apperror.ErrStatusCodeInvalid,
	},
}
