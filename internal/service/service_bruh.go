package service

import (
	"context"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/models"
	"github.com/go-playground/validator/v10"
	"log"
)

//const MAX_CAPACITY = uint64(10)

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

func (s *Service) LoadCache(ctx context.Context, cacheSize uint64) error {
	ids, err := s.repo.GetRecentIDs(ctx, cacheSize)
	if err != nil {
		return err
	}
	orders := make([]*models.Order, 0)
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

var validate = validator.New()

func ValidateOrder(order *models.Order) error {
	return validate.Struct(order)
}
