package main

import (
	"context"
	"errors"
	cache "github.com/GameXost/wbTestCase/cache"
	"github.com/GameXost/wbTestCase/config"
	repository "github.com/GameXost/wbTestCase/internal/repository"
	"github.com/GameXost/wbTestCase/internal/server"
	service "github.com/GameXost/wbTestCase/internal/service"
	"github.com/GameXost/wbTestCase/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolConf, err := pgxpool.ParseConfig(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("failed to parse config, %v", err)
	}
	poolConf.MaxConns = int32(cfg.DB.PoolMaxConns)
	poolConf.MinConns = int32(cfg.DB.PoolMinConns)
	poolConf.MaxConnIdleTime = cfg.DB.PoolMaxIdleTime
	poolConf.MaxConnLifetime = cfg.DB.PoolMaxLifeTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	log.Println("postgress good")

	orderRepo := repository.NewRepo(pool)
	orderCache := cache.NewCache(cfg.Cache.Size)
	orderService := service.NewService(orderRepo, orderCache)
	orderServer := server.NewHandler(orderService)
	log.Println("initialized all layers")

	if err = orderService.LoadCache(ctx); err != nil {
		log.Printf("failed to restore cache from db: %v", err)
	} else {
		log.Println("Cache loaded")
	}
	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Group, orderService)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	consumerErrs := make(chan error, 1)
	go func() {
		log.Println("kafka consumer started")
		consumerErrs <- consumer.Start(ctx)
	}()

	router := SetupRouter(orderServer)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("HTTP port: %s", cfg.Server.Port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {

	case err = <-consumerErrs:
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Fatalf("kafka died - critical error: %v", err)
		}
	case err = <-serverErrors:
		log.Fatalf("server error: %v", err)
	case sig := <-shutdown:
		log.Printf("shitdown signal: %v", sig)
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		if err = httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
			httpServer.Close()
		}
		log.Println("server stopped gracefully")
	}
}

func SetupRouter(handler *server.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/order/{order_uid}", handler.GetOrder)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	filesDir := http.Dir("web")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	r.Handle("/*", http.FileServer(filesDir))

	return r
}
