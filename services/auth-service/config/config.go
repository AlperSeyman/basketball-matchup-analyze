package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DBUrl                 string
	RedisUrl              string
	JWTSecret             string
	AccessTokenTTLMinutes int
	RefreshTokenTTLDays   int
	Port                  string
	Environment           string
}

func Load() (*Config, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		slog.Warn("no .env file found, reading from environment")
	}

	dbUrl, err := buildDBUrl()
	if err != nil {
		return nil, err
	}

	jwtSecret, err := requireEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters, got %d", len(jwtSecret))
	}

	accessTTL, err := getIntEnv("JWT_ACCESS_TOKEN_TTL_MINUTES", 15)
	if err != nil {
		return nil, err
	}

	refreshTTL, err := getIntEnv("JWT_REFRESH_TOKEN_TTL_DAYS", 7)
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUrl:                 dbUrl,
		RedisUrl:              getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:             jwtSecret,
		AccessTokenTTLMinutes: accessTTL,
		RefreshTokenTTLDays:   refreshTTL,
		Port:                  getEnv("AUTH_SERVICE_PORT", "8001"),
		Environment:           getEnv("ENVIRONMENT", "development"),
	}, nil
}

func ConnectDb(cfg *Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	slog.Info("connected to PostgreSQL")
	return pool, nil
}

func ConnectRedis(cfg *Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	rdb := redis.NewClient(opts)
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	slog.Info("connected to Redis")
	return rdb, nil
}

func buildDBUrl() (string, error) {
	user, err := requireEnv("POSTGRES_USER")
	if err != nil {
		return "", err
	}
	password, err := requireEnv("POSTGRES_PASSWORD")
	if err != nil {
		return "", err
	}
	db, err := requireEnv("POSTGRES_DB")
	if err != nil {
		return "", err
	}

	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	sslmode := getEnv("POSTGRES_SSLMODE", "disable")

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   "/" + db,
	}
	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func requireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return value, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for environment variable %s: %w", key, err)
	}
	return parsed, nil
}
