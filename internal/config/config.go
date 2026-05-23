package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddress       string
	DatabaseDSN       string
	SchedulerInterval time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		HTTPAddress:       getEnv("HTTP_ADDRESS", ":8080"),
		DatabaseDSN:       getEnv("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/go_reminder?charset=utf8mb4&parseTime=True&loc=Local"),
		SchedulerInterval: getDurationFromSeconds("SCHEDULER_INTERVAL_SECONDS", 60*time.Second),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getDurationFromSeconds(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return fallback
	}

	return time.Duration(seconds) * time.Second
}
