package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/basketball-matchup-analyzer/auth-service/config"
	"github.com/basketball-matchup-analyzer/auth-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	cfg *config.Config
	db  *pgxpool.Pool
	rdb *redis.Client
	q   *repository.Queries
}

func NewHandler(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, q *repository.Queries) *Handler {
	return &Handler{
		cfg: cfg,
		db:  db,
		rdb: rdb,
		q:   q,
	}
}

func (h *Handler) Health(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "ok"
	err := h.db.Ping(ctx)
	if err != nil {
		dbStatus = "down"
	}

	redisStatus := "ok"
	err = h.rdb.Ping(ctx).Err()
	if err != nil {
		redisStatus = "down"
	}

	httpStatus := http.StatusOK
	overallStatus := "ok"
	if dbStatus == "down" || redisStatus == "down" {
		httpStatus = http.StatusServiceUnavailable
		overallStatus = "degraded"
	}

	c.JSON(httpStatus, gin.H{
		"status":   overallStatus,
		"postgres": dbStatus,
		"redis":    redisStatus,
	})

}
