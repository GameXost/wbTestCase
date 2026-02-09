package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/models"
	"log"
)

var ()

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
	if order.OrderUId == "" {
		return errors.New("need order_uid")
	}
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
	ids, err1 := s.repo.GetRecentIDs(ctx, MAX_CAPACITY)
	if err1 != nil {
		return err1
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
		return errors.New("order_uid is missing")
	}
	if order.TrackNumber == "" {
		return errors.New("track number is missing")
	}
	if order.Entry == "" {
		return errors.New("entry is missing")
	}
	if order.Locale == "" {
		return errors.New("locale is missing")
	}
	if order.CustomerId == "" {
		return errors.New("customer id is missing")
	}
	if order.DeliveryService == "" {
		return errors.New("delivery service is missing")
	}
	if order.Shardkey == "" {
		return errors.New("shardkey is missing")
	}
	if order.SmId <= 0 {
		return errors.New("invalid sm id")
	}

	if order.Delivery.Name == "" {
		return errors.New("delivery name is missing")
	}
	if order.Delivery.Phone == "" {
		return errors.New("delivery phone is missing")
	}
	if order.Delivery.Zip == "" {
		return errors.New("delivery ZIP is missing")
	}
	if order.Delivery.City == "" {
		return errors.New("delivery city is missing")
	}
	if order.Delivery.Address == "" {
		return errors.New("delivery address is missing")
	}
	if order.Delivery.Region == "" {
		return errors.New("delivery region is missing")
	}
	if order.Delivery.Email == "" {
		return errors.New("delivery email is missing")
	}

	if order.Payment.Transaction == "" {
		return errors.New("payment transaction is missing")
	}
	if order.Payment.RequestId == "" {
		return errors.New("payment request id is missing")
	}
	if order.Payment.Currency == "" {
		return errors.New("payment currency is missing")
	}
	if order.Payment.Provider == "" {
		return errors.New("payment provider is missing")
	}
	if order.Payment.Amount <= 0 {
		return errors.New("payment amount is invalid")
	}
	if order.Payment.Bank == "" {
		return errors.New("payment bank is missing")
	}
	if order.Payment.DeliveryCost < 0 {
		return errors.New("payment delivery cost is invalid, lower zero")
	}
	if order.Payment.GoodsTotal <= 0 {
		return errors.New("payment goods total is invalid")
	}

	if len(order.Items) == 0 {
		return errors.New("items empty")
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
		return fmt.Errorf("item name is missing, id %v", item.Id)
	}
	if item.ChrtId <= 0 {
		return fmt.Errorf("item %s chrt id is invalid", item.Name)
	}
	if item.TrackNumber == "" {
		return fmt.Errorf("item %s track number is missing", item.Name)
	}
	if item.RID == "" {
		return fmt.Errorf("item %s RID is missing", item.Name)
	}
	if item.NmId <= 0 {
		return fmt.Errorf("item %s nm id is invalid", item.Name)
	}
	if item.Price < 0 {
		return fmt.Errorf("item %s price is invalid", item.Name)
	}
	if item.Sale < 0 {
		return fmt.Errorf("item %s sale is invalid", item.Name)
	}
	if item.TotalPrice < 0 {
		return fmt.Errorf("item %s total price is invalid", item.Name)
	}
	if item.Status < 0 {
		return fmt.Errorf("item %s status code is invalid", item.Name)
	}
	return nil
}
