package middleware

import (
	"crypto/ed25519"
	"log/slog"
	"net/http"
	"strings"

	"github.com/basketball-matchup-analyzer/auth-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RequireAuth(publicKey ed25519.PublicKey, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Warn("missing authorization header", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			slog.Warn("invalid authorization header format", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString, publicKey)
		if err != nil {
			slog.Warn("invalid token", "error", err, "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if claims.TokenType != "access" {
			slog.Warn("wrong token type", "type", claims.TokenType, "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		_, err = rdb.Get(c.Request.Context(), "blocklist:"+tokenString).Result()
		if err == nil {
			slog.Warn("blocklisted token used", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if err != redis.Nil {
			slog.Error("redis blocklist check failed", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
