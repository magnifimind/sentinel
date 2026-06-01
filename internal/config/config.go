package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ListenAddr string
	DB         DBConfig
	Kafka      KafkaConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

type KafkaConfig struct {
	Brokers        []string
	ConsumerGroup  string
	Topics         []string
	SessionTimeout time.Duration
}

func Load() Config {
	return Config{
		ListenAddr: envOr("LISTEN_ADDR", ":8080"),
		DB: DBConfig{
			Host:     envOr("DB_HOST", "localhost"),
			Port:     envIntOr("DB_PORT", 5432),
			User:     envOr("DB_USER", "sentinel"),
			Password: envOr("DB_PASSWORD", "sentinel"),
			Name:     envOr("DB_NAME", "sentinel"),
			SSLMode:  envOr("DB_SSLMODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers:        []string{envOr("KAFKA_BROKERS", "localhost:9092")},
			ConsumerGroup:  envOr("KAFKA_CONSUMER_GROUP", "sentinel-api"),
			Topics:         []string{"agent.decisions", "agent.actions", "agent.outcomes"},
			SessionTimeout: 10 * time.Second,
		},
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
