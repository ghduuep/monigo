package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/labstack/echo/v4"
)

// @Summary Get user channels
// @Description List all notification channels configured by the user
// @Tags channels
// @Security BearerAuth
// @Success 200 {array} models.NotificationChannel
// @Router /channels [get]
func (h *Handler) GetChannels(c echo.Context) error {
	userID := getUserIdFromToken(c)

	cacheKey := fmt.Sprintf("user:%d:channels", userID)

	var channels []models.NotificationChannel

	if h.getCache(c.Request().Context(), cacheKey, &channels) {
		return c.JSON(http.StatusOK, channels)
	}

	channels, err := database.GetUserChannels(c.Request().Context(), h.DB, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch channels."})
	}

	go h.setCache(context.Background(), cacheKey, channels, 1*time.Hour)

	return c.JSON(http.StatusOK, channels)
}

// @Summary Create channel
// @Description Add a new notification channel (email, telegram)
// @Tags channels
// @Security BearerAuth
// @Param request body dto.CreateChannelRequest true "Channel Info"
// @Success 201 {object} models.NotificationChannel
// @Router /channels [post]
func (h *Handler) CreateChannel(c echo.Context) error {
	var req dto.CreateChannelRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data."})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	userID := getUserIdFromToken(c)

	channel := models.NotificationChannel{
		UserID: userID,
		Type:   models.NotificationType(req.Type),
		Target: req.Target,
	}

	if err := database.CreateChannel(c.Request().Context(), h.DB, &channel); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create channel."})
	}

	cacheKey := fmt.Sprintf("user:%d:channels", userID)
	h.invalidateCache(c.Request().Context(), cacheKey)

	return c.JSON(http.StatusCreated, channel)
}

// @Summary Delete channel
// @Description Remove a notification channel
// @Tags channels
// @Security BearerAuth
// @Param id path int true "Channel ID"
// @Success 204
// @Router /channels/{id} [delete]
func (h *Handler) DeleteChannel(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID."})
	}

	userID := getUserIdFromToken(c)

	if err := database.DeleteChannel(c.Request().Context(), h.DB, id, userID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Channel not found."})
	}

	cacheKey := fmt.Sprintf("user:%d:channels", userID)
	h.invalidateCache(c.Request().Context(), cacheKey)

	return c.NoContent(http.StatusNoContent)
}
