package main

import (
	"context"
	"errors"
	"github.com/GameXost/wbTestCase/internal/config"
	"github.com/GameXost/wbTestCase/internal/kafka"
	repository "github.com/GameXost/wbTestCase/internal/repository"
	"github.com/GameXost/wbTestCase/internal/repository/cache"
	"github.com/GameXost/wbTestCase/internal/server"
	service "github.com/GameXost/wbTestCase/internal/service"
	"github.com/GameXost/wbTestCase/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// повыносить из main всякую хуйню
func main() {
	//config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	//context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//bd pool
	pool, err := initDB(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer pool.Close()

	//services
	orderService, orderHandler := initLayers(pool, cfg)

	//cache preload
	if err := orderService.LoadCache(ctx, cfg.Cache.Size); err != nil {
		log.Printf("failed to restore cache from db: %v", err)
	} else {
		log.Println("Cache loaded")
	}

	//kafka
	consumerErrs, err := startKafka(ctx, cfg, orderService)
	if err != nil {
		log.Fatalf("failed to init kafka: %v", err)
	}
	//prometheus
	registerMetrics()

	//router handlers
	router := SetupRouter(orderHandler)
	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("HTTP port: %s", cfg.Server.Port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	//shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-consumerErrs:
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("kafka died - critical error: %v", err)
			return
		}
	case err := <-serverErrors:
		log.Printf("server error: %v", err)
		return
	case sig := <-shutdown:
		log.Printf("shitdown signal: %v", sig)
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		if errHttp := httpServer.Close(); errHttp != nil {
			log.Printf("closing httpServer failed: %v", errHttp)
		}
	}
	log.Println("server stopped gracefully")
}

func SetupRouter(handler *server.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/order/{order_uid}", handler.GetOrder)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	filesDir := http.Dir("web")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	r.Handle("/*", http.FileServer(filesDir))
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func initDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConf, err := pgxpool.ParseConfig(cfg.DB.DSN())
	if err != nil {
		return nil, err
	}
	poolConf.MaxConns = int32(cfg.DB.PoolMaxConns)
	poolConf.MinConns = int32(cfg.DB.PoolMinConns)
	poolConf.MaxConnIdleTime = cfg.DB.PoolMaxIdleTime
	poolConf.MaxConnLifetime = cfg.DB.PoolMaxLifeTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	log.Println("postgress good")
	return pool, nil
}

func initLayers(pool *pgxpool.Pool, cfg *config.Config) (*service.Service, *server.Handler) {
	orderRepo := repository.NewRepo(pool)
	orderCache := cache.NewCache(cfg.Cache.Size)
	orderService := service.NewService(orderRepo, orderCache)
	orderHandler := server.NewHandler(orderService)
	log.Println("initialized all layers")
	return orderService, orderHandler
}

func startKafka(ctx context.Context, cfg *config.Config, srvs *service.Service) (<-chan error, error) {

	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.Group,
		srvs,
		cfg.Kafka.DLQTopic,
	)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 1)

	go func() {
		log.Println("Kafka consumer started")
		if err := consumer.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	return errChan, nil
}

func registerMetrics() {
	prometheus.MustRegister(
		metrics.RequestsTotal,
		metrics.RequestsServerError,
		metrics.RequestsBadRequest,
		metrics.RequestsNotFound,
		metrics.CacheHits,
		metrics.CacheMisses,
		metrics.RequestsSuccess,
	)
}
