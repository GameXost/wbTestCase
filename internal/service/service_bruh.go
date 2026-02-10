package service

import (
	"context"
	"fmt"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/models"
	"log"
)

const MAX_CAPACITY = uint64(10)

type OrderRepo interface {
	GetRecentIDs(ctx context.Context, amount uint64) ([]string, error)
	CreateFullOrder(ctx context.Context, order *models.Order) error
	GetFullOrderOnId(ctx context.Context, OrderUId string) (*models.Order, error)
}

type OrderCache interface {
	Get(key string) (*models.Order, bool)
	Set(order *models.Order)
	LoadFull(ids []*models.Order)
}

type Service struct {
	repo  OrderRepo
	cache OrderCache
}

func NewService(repo OrderRepo, cache OrderCache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := ValidateOrder(order); err != nil {
		log.Printf("inalid order data: %v", err)
		return errHandle.ErrValidation
	}
	err := s.repo.CreateFullOrder(ctx, order)
	if err != nil {
		return err
	}
	s.cache.Set(order)
	return nil
}

func (s *Service) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	order, has := s.cache.Get(orderUID)
	if has {
		return order, nil
	}
	order, err := s.repo.GetFullOrderOnId(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	s.cache.Set(order)
	return order, nil
}

func (s *Service) LoadCache(ctx context.Context) error {
	ids, err := s.repo.GetRecentIDs(ctx, MAX_CAPACITY)
	if err != nil {
		return err
	}
	orders := make([]*models.Order, 0, MAX_CAPACITY)
	for _, id := range ids {
		order, err := s.repo.GetFullOrderOnId(ctx, id)
		if err != nil {
			log.Printf("some incorrect data %s with error: %v", id, err)
			continue
		}
		orders = append(orders, order)
	}
	s.cache.LoadFull(orders)
	return nil
}

func ValidateOrder(order *models.Order) error {

	if order.OrderUId == "" {
		return errHandle.ErrOrderUIDMissing
	}
	if order.TrackNumber == "" {
		return errHandle.ErrTrackNumberMissing
	}
	if order.Entry == "" {
		return errHandle.ErrEntryMissing
	}
	if order.Locale == "" {
		return errHandle.ErrLocaleMissing
	}
	if order.CustomerId == "" {
		return errHandle.ErrCustomerIDMissing
	}
	if order.DeliveryService == "" {
		return errHandle.ErrDeliveryServiceMissing
	}
	if order.Shardkey == "" {
		return errHandle.ErrShardkeyMissing
	}
	if order.SmId <= 0 {
		return errHandle.ErrInvalidSmID
	}

	if order.Delivery.Name == "" {
		return errHandle.ErrDeliveryNameMissing
	}
	if order.Delivery.Phone == "" {
		return errHandle.ErrDeliveryPhoneMissing
	}
	if order.Delivery.Zip == "" {
		return errHandle.ErrDeliveryZIPMissing
	}
	if order.Delivery.City == "" {
		return errHandle.ErrDeliveryCityMissing
	}
	if order.Delivery.Address == "" {
		return errHandle.ErrDeliveryAddressMissing
	}
	if order.Delivery.Region == "" {
		return errHandle.ErrDeliveryRegionMissing
	}
	if order.Delivery.Email == "" {
		return errHandle.ErrDeliveryEmailMissing
	}

	if order.Payment.Transaction == "" {
		return errHandle.ErrPaymentTransactionMissing
	}
	if order.Payment.RequestId == "" {
		return errHandle.ErrPaymentRequestMissing
	}
	if order.Payment.Currency == "" {
		return errHandle.ErrPaymentCurrencyMissing
	}
	if order.Payment.Provider == "" {
		return errHandle.ErrPaymentProviderMissing
	}
	if order.Payment.Amount <= 0 {
		return errHandle.ErrPaymentAmountInvalid
	}
	if order.Payment.Bank == "" {
		return errHandle.ErrPaymentBankMissing
	}
	if order.Payment.DeliveryCost < 0 {
		return errHandle.ErrPaymentDeliveryInvalid
	}
	if order.Payment.GoodsTotal <= 0 {
		return errHandle.ErrPaymentGoodsTotalInvalid
	}

	if len(order.Items) == 0 {
		return errHandle.ErrItemsEmpty
	}

	for i := range order.Items {
		if err := validateItem(&order.Items[i]); err != nil {
			return err
		}
	}

	return nil
}

func validateItem(item *models.Item) error {
	if item.Name == "" {
		return fmt.Errorf("id %v %w", item.Id, errHandle.ErrItemNameMissing)
	}
	if item.ChrtId <= 0 {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemChrtMissing)
	}
	if item.TrackNumber == "" {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemTrackNumberMissing)
	}
	if item.RID == "" {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemRIDMissing)
	}
	if item.NmId <= 0 {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemNmIdInvalid)
	}
	if item.Price < 0 {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemPriceInvalid)
	}
	if item.Sale < 0 {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemSaleInvalid)
	}
	if item.TotalPrice < 0 {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrItemTotalPriceInvalid)
	}
	if item.Status < 0 {
		return fmt.Errorf("item: %s %w", item.Name, errHandle.ErrStatusCodeInvalid)
	}
	return nil
}
