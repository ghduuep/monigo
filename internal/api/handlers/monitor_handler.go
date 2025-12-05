package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetMonitors(c echo.Context) error {
	monitors, err := database.GetAllMonitors(c.Request().Context(), h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get monitors."})
	}

	return c.JSON(http.StatusOK, monitors)
}

func (h *Handler) GetMonitorByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	monitor, err := database.GetMonitorByID(c.Request().Context(), h.DB, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	return c.JSON(http.StatusOK, monitor)
}

func (h *Handler) CreateMonitor(c echo.Context) error {
	var req dto.MonitorRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data."})
	}

	monitor := models.Monitor{
		UserID:    req.UserID,
		Target:    req.Target,
		Type:      req.Type,
		Config:    req.Config,
		Interval:  req.Interval,
		CreatedAt: time.Now(),
	}

	if err := database.CreateMonitor(c.Request().Context(), h.DB, &monitor); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Monitor already exists."})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) DeleteMonitor(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	if err := database.DeleteMonitor(c.Request().Context(), h.DB, id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	return c.NoContent(http.StatusOK)
}
