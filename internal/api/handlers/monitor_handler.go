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

// @Summary Get all monitors
// @Description Retrieve all monitors created by the authenticated user.
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Monitor
// @Failure 500 {object} map[string]string
// @Router /monitors [get]
func (h *Handler) GetMonitors(c echo.Context) error {
	userID := getUserIdFromToken(c)
	monitors, err := database.GetMonitorsByUserID(c.Request().Context(), h.DB, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get monitors."})
	}

	return c.JSON(http.StatusOK, monitors)
}

// @Summary Get monitor by ID
// @Description Retrieve a specific monitor details.
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} models.Monitor
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /monitors/{id} [get]
func (h *Handler) GetMonitorByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	userID := getUserIdFromToken(c)

	monitor, err := database.GetMonitorByIDAndUser(c.Request().Context(), h.DB, id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	return c.JSON(http.StatusOK, monitor)
}

// @Summary Create a new monitor
// @Description Create a new monitor for a target (HTTP, DNS, etc).
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MonitorRequest true "Monitor Configuration"
// @Success 201 {object} nil
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /monitors [post]
func (h *Handler) CreateMonitor(c echo.Context) error {
	var req dto.MonitorRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data."})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	userID := getUserIdFromToken(c)

	intervalDuration, _ := time.ParseDuration(req.Interval)
	timeoutDuration, _ := time.ParseDuration(req.Timeout)

	monitor := models.Monitor{
		UserID:    userID,
		Target:    req.Target,
		Type:      req.Type,
		Config:    req.Config,
		Interval:  intervalDuration,
		Timeout:   timeoutDuration,
		CreatedAt: time.Now(),
	}

	if err := database.CreateMonitor(c.Request().Context(), h.DB, &monitor); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Monitor already exists."})
	}

	return c.NoContent(http.StatusCreated)
}

// @Summary Delete a monitor
// @Description Remove a monitor by its ID.
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /monitors/{id} [delete]
func (h *Handler) DeleteMonitor(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	userID := getUserIdFromToken(c)

	if err := database.DeleteMonitor(c.Request().Context(), h.DB, id, userID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Get monitor stats
// @Description Get uptime and latency stats for a monitor
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Monitor ID"
// @Success 200 {object} dto.MonitorStatsResponse
// @Router /monitors/{id}/stats [get]
func (h *Handler) GetMonitorStats(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	userID := getUserIdFromToken(c)

	_, err = database.GetMonitorByIDAndUser(c.Request().Context(), h.DB, id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	var from, to time.Time
	to = time.Now()

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	period := c.QueryParam("period")

	if startDateStr != "" && endDateStr != "" {
		layout := "2006-01-02"

		parsedStart, err := time.Parse(layout, startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date. Use YYYY-MM-DD format"})
		}

		from = parsedStart

		parsedEnd, err := time.Parse(layout, endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date. Use YYYY-MM-DD format"})
		}

		to = parsedEnd.Add(24 * time.Hour).Add(-1 * time.Second)

	} else {
		switch period {
		case "7d":
			from = time.Now().AddDate(0, 0, -7)
		case "30d":
			from = time.Now().AddDate(0, 0, -30)
		case "24h":
			from = time.Now().Add(-24 * time.Hour)
		default:
			from = time.Now().Add(-24 * time.Hour)
		}
	}

	stats, err := database.GetMonitorStats(c.Request().Context(), h.DB, id, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get stats."})
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetMonitorLastChecks(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	userID := getUserIdFromToken(c)

	_, err = database.GetMonitorByIDAndUser(c.Request().Context(), h.DB, id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	checks, err := database.GetLastChecks(c.Request().Context(), h.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get checks."})
	}

	return c.JSON(http.StatusOK, checks)
}

func (h *Handler) GetMonitorLastIncidents(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	userID := getUserIdFromToken(c)

	_, err = database.GetMonitorByIDAndUser(c.Request().Context(), h.DB, id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Monitor not found."})
	}

	incidents, err := database.GetIncidentsByID(c.Request().Context(), h.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get incidents."})
	}

	return c.JSON(http.StatusOK, incidents)
}
