package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/basketball-matchup-analyzer/auth-service/config"
	"github.com/basketball-matchup-analyzer/auth-service/internal/api/handlers"
	"github.com/basketball-matchup-analyzer/auth-service/internal/api/middleware"
	"github.com/basketball-matchup-analyzer/auth-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if c.Request.URL.Path == "/health" {
			return
		}
		slog.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
		)
	}
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

func maxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client) http.Handler {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())
	r.Use(securityHeaders())
	r.Use(maxBodySize(1 << 20)) // 1 MB

	q := repository.New(db)
	h := handlers.NewHandler(cfg, db, rdb, q)

	r.GET("/health", h.Health)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/logout", h.Logout)

		protected := auth.Group("")
		protected.Use(middleware.RequireAuth(cfg.JWTPublicKey, rdb))
		{
			protected.POST("/change-password", h.ChangePassword)
		}
	}

	return r
}
