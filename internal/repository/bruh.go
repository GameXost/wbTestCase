package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}
func GetOrder(ctx context.Context, OrderUId string) models.Order {
}

func CreateOrder() {

}
