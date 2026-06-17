package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/basketball-matchup-analyzer/auth-service/internal/model"
	"github.com/basketball-matchup-analyzer/auth-service/internal/repository"
	"github.com/basketball-matchup-analyzer/auth-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Register(c *gin.Context) {

	var input model.RegisterInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	user, err := h.q.CreateUser(c.Request.Context(), repository.CreateUserParams{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if utils.IsDuplicateEmail(err) {
			c.JSON(http.StatusConflict, model.ErrorResponse{Error: "email already in use"})
			return
		}
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Message: "account created",
		Data: gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {

	var input model.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}
	user, err := h.q.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
			return
		}
		slog.Error("failed to get user", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	ok, err := utils.VerifyPassword(user.PasswordHash, input.Password)
	if err != nil {
		slog.Error("failed to verify password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.cfg.JWTSecret, h.cfg.AccessTokenTTLMinutes)
	if err != nil {
		slog.Error("failed to generate access token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, h.cfg.JWTSecret, h.cfg.RefreshTokenTTLDays)
	if err != nil {
		slog.Error("failed to generate refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refreshToken, h.cfg.RefreshTokenTTLDays*24*60*60, "/", "", h.cfg.Environment == "production", true)

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Message: "logged in",
		Data: gin.H{
			"access_token": accessToken,
		},
	})
}

func (h *Handler) Logout(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 {
		accessToken := parts[1]
		claims, err := utils.ValidateToken(accessToken, h.cfg.JWTSecret)
		if err == nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				if err := h.rdb.Set(c.Request.Context(), "blocklist:"+accessToken, "1", ttl).Err(); err != nil {
					slog.Error("failed to blocklist access token", "error", err)
				}
			}
		}
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		claims, err := utils.ValidateToken(refreshToken, h.cfg.JWTSecret)
		if err == nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				if err := h.rdb.Set(c.Request.Context(), "blocklist:"+refreshToken, "1", ttl).Err(); err != nil {
					slog.Error("failed to blocklist refresh token", "error", err)
				}
			}
		}
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cfg.Environment == "production", true)

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Message: "logged out",
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "missing refresh token"})
		return
	}

	claims, err := utils.ValidateToken(refreshToken, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid or expired refresh token"})
		return
	}

	if claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid token type"})
		return
	}

	blocked, err := h.rdb.Exists(c.Request.Context(), "blocklist:"+refreshToken).Result()
	if err != nil {
		slog.Error("failed to check refresh token block", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	if blocked > 0 {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "refresh token has been revoked"})
		return
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	err = h.rdb.Set(c.Request.Context(), "blocklist:"+refreshToken, "1", ttl).Err()
	if err != nil {
		slog.Error("failed to blocklist old refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	user, err := h.q.GetUserByID(c.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "user no longer exists"})
			return
		}
		slog.Error("failed to fetch user on refresh", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.cfg.JWTSecret, h.cfg.AccessTokenTTLMinutes)

	if err != nil {
		slog.Error("failed to generate access token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, h.cfg.JWTSecret, h.cfg.RefreshTokenTTLDays)
	if err != nil {
		slog.Error("failed to generate refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", newRefreshToken, h.cfg.RefreshTokenTTLDays*24*60*60, "/", "", h.cfg.Environment == "production", true)

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Message: "token refreshed",
		Data: gin.H{
			"access_token": accessToken,
		},
	})
}

func (h *Handler) ForgotPassword(c *gin.Context) {

	var input model.ForgotPasswordInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	user, err := h.q.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "if that email exists, a reset link has been sent"})
			return
		}
		slog.Error("failed to look up email for password reset", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	rawToken, err := utils.GenerateSecureToken()
	if err != nil {
		slog.Error("failed to generate reset token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	hashedToken := utils.HashToken(rawToken)

	expiresAt := time.Now().Add(time.Hour)
	_, err = h.q.SetPasswordResetToken(c.Request.Context(), repository.SetPasswordResetTokenParams{
		Email:                  user.Email,
		PasswordResetToken:     pgtype.Text{String: hashedToken, Valid: true},
		PasswordResetExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		slog.Error("failed to save password reset token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	// TODO: send rawToken to user.Email via SendGrig

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Message: "if that email exists, a reset link has been sent",
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {

	var input model.ResetPasswordInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	hashedToken := utils.HashToken(input.Token)
	user, err := h.q.GetUserByPasswordResetToken(c.Request.Context(), pgtype.Text{String: hashedToken, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid or expired reset token"})
			return
		}
		slog.Error("failed to look up password reset token", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	newHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		slog.Error("failed to hash new password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	_, err = h.q.ResetPassword(c.Request.Context(), repository.ResetPasswordParams{
		PasswordHash: newHash,
		ID:           user.ID,
	})
	if err != nil {
		slog.Error("failed to reset password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Message: "password reset",
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {

	var input model.ChangePasswordInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	user, err := h.q.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		slog.Error("failed to fetch user for password change", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}

	ok, err := utils.VerifyPassword(user.PasswordHash, input.CurrentPassword)
	if err != nil {
		slog.Error("failed to verify password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
		return
	}

	newHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	_, err = h.q.UpdatePassword(c.Request.Context(), repository.UpdatePasswordParams{
		PasswordHash: newHash,
		ID:           userID,
	})
	if err != nil {
		slog.Error("failed to update password", "error", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Message: "pasword change",
	})

}
