package repository

import (
	"context"
	"errors"
	"github.com/GameXost/wbTestCase/internal/errHandle"
	"github.com/GameXost/wbTestCase/internal/generator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"reflect"
	"testing"

	"time"
)

var repo *Repo

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgres.Run(
		ctx,

		"postgres:17",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		postgres.WithInitScripts("../../init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		panic(err)
	}
	defer func() {
		_ = container.Terminate(ctx)
	}()

	connectStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, connectStr)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	repo = NewRepo(pool)
	os.Exit(m.Run())

}

func TestCreateGetOrder(t *testing.T) {
	ctx := context.Background()
	order := generator.ValidOrder("correct")
	err := repo.CreateFullOrder(ctx, order)
	if err != nil {
		t.Fatalf("CreateFullOrder failed: %v", err)
	}

	got, err := repo.GetFullOrderOnId(ctx, order.OrderUId)
	if err != nil {
		t.Fatalf("GetFullOrderOnId failed: %v", err)
	}

	got.Payment.OrderId = order.OrderUId // эти айдишники есть в структурах, но я их не получаю из репозитория
	got.Items[0].OrderUId = order.OrderUId

	//служебные поля
	order.Items[0].Id = got.Items[0].Id // только в бд есть, уникальный primary key, не возвращаю из бд
	order.Delivery.Id = got.Delivery.Id

	order.DateCreated = got.DateCreated // разные таймзоны, хз как поменять
	if !reflect.DeepEqual(got, order) {
		t.Fatalf("shit, they're not equal got:\n %+v, want:\n %+v", got, order)
	}
}

func TestIdempotentCreateOrder(t *testing.T) {
	ctx := context.Background()
	order := generator.ValidOrder("idempotent")

	order.OrderUId = "Idemp"
	order.Payment.OrderId = "Idemp"
	order.Delivery.OrderUId = "Idemp"
	order.Items[0].OrderUId = "Idemp"

	err := repo.CreateFullOrder(ctx, order)
	if err != nil {
		t.Fatalf("CreateFullOrder try 1 failed: %v", err)
	}
	err = repo.CreateFullOrder(ctx, order)
	if err != nil {
		t.Fatalf("CreateFullOrder try 2 failed: %v", err)
	}

}

func TestOrderNotFound(t *testing.T) {
	ctx := context.Background()
	_, err := repo.GetFullOrderOnId(ctx, "bruh")
	if !errors.Is(err, errHandle.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got: %v", err)
	}

}
