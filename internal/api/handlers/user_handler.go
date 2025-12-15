package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) GetUser(c echo.Context) error {
	userID := getUserIdFromToken(c)

	cacheKey := fmt.Sprintf("user:%d:profile", userID)

	var response dto.UserResponse
	if h.getCache(c.Request().Context(), cacheKey, &response) {
		return c.JSON(http.StatusOK, response)
	}

	user, err := database.GetUserByID(c.Request().Context(), h.DB, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found."})
	}

	response = dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	go h.setCache(context.Background(), cacheKey, response, 1*time.Hour)

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	userID := getUserIdFromToken(c)

	if err := database.DeleteUser(c.Request().Context(), h.DB, userID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found."})
	}

	cacheKey := fmt.Sprintf("user:%d:profile", userID)
	cacheChannels := fmt.Sprintf("user:%d:channels", userID)
	h.invalidateCache(c.Request().Context(), cacheKey, cacheChannels)

	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	userID := getUserIdFromToken(c)

	var dto dto.UpdateUserRequest
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data."})
	}

	if err := c.Validate(&dto); err != nil {
		return err
	}

	ctx := c.Request().Context()

	if dto.Password != nil {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate the password hash"})
		}

		tokenString := extractTokenString(c)
		*dto.Password = string(hashedBytes)
		_ = database.AddToBlackList(ctx, h.RDB, tokenString, 72*time.Hour)
	}

	err := database.UpdateUser(ctx, h.DB, userID, dto)
	if err != nil {
		if err.Error() == "no data" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "E-mail already used."})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user."})
	}

	cacheKey := fmt.Sprintf("user:%d:profile", userID)
	h.invalidateCache(ctx, cacheKey)

	return c.NoContent(http.StatusOK)
}
