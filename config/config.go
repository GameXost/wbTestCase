package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB     DBConfig
	Kafka  KafkaConfig
	Server ServerConfig
	Cache  CacheConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string

	PoolMaxConns    int
	PoolMinConns    int
	PoolMaxLifeTime time.Duration
	PoolMaxIdleTime time.Duration
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	Group    string
	DLQTopic string
}

type ServerConfig struct {
	Port string
}

type CacheConfig struct {
	Size uint64
}

// посморел, что хорошая практика писать отдельный пакет для загрузки конфига :3
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("BD_PORT", "5432"),
			Name:            getEnv("DB_NAME", "wbcase"),
			User:            getEnv("DB_USER", "wbuser"),
			Password:        getEnv("DB_PASSWORD", "wbpass"),
			PoolMaxConns:    getIntEnv("DB_POOL_MAX_CONNS", 10),
			PoolMinConns:    getIntEnv("DB_POOL_MIN_CONNS", 2),
			PoolMaxLifeTime: getDurationEnv("DB_POOL_MAX_LIFE_TIME", time.Hour),
			PoolMaxIdleTime: getDurationEnv("DB_POOL_MAX_IDLE_TIME", 30*time.Minute),
		},
		Kafka: KafkaConfig{
			Brokers:  []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topic:    getEnv("KAFKA_TOPIC", "orders"),
			Group:    getEnv("KAFKA_GROUP", "order_consumers"),
			DLQTopic: getEnv("KAFKA_TOPIC_DLQ", "orders.dlq"),
		},
		Server: ServerConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Cache: CacheConfig{
			Size: uint64(getIntEnv("CACHE_SIZE", 10)),
		},
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil

}
func (c *Config) Validate() error {
	if c.DB.Host == "" {
		return fmt.Errorf("DB_HOST is empty")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("DB_PASSOWRD is empty")
	}
	if c.Cache.Size <= 0 {
		return fmt.Errorf("CACHE_SIZE is lower or is 0")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.Name)
}
