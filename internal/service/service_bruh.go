package service

import (
	"context"
	"errors"
	"github.com/GameXost/wbTestCase/cache"
	"github.com/GameXost/wbTestCase/internal/repository"
	"github.com/GameXost/wbTestCase/models"
	"log"
)

const MAX_CAPACITY = uint64(10)

type Service struct {
	repo  *repository.Repo
	cache *cache.Cache
}

func NewService(repo *repository.Repo, cache *cache.Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *models.Order) error {
	if order.OrderUId == "" {
		return errors.New("need order_uid")
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
