package handlers

import (
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/labstack/echo/v4"
)

// @Summary Get all incidents
// @Description Retrieve all incidents for the authenticated user across all monitors.
// @Tags monitors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Incident
// @Router /incidents [get]
func (h *Handler) GetIncidents(c echo.Context) error {
	userID := getUserIdFromToken(c)
	page, limit, offset := getPaginationParams(c)

	qStart := c.QueryParam("start_date")
	qEnd := c.QueryParam("end_date")
	qPeriod := c.QueryParam("period")

	qTarget := c.QueryParam("target")
	qError := c.QueryParam("error_cause")

	cacheKey := fmt.Sprintf("user:%d:incidents:%s:%s:%s:%s:%s%d:%d", userID, qStart, qEnd, qPeriod, qTarget, qError, page, limit)

	var incidents []*models.Incident
	if h.getCache(c.Request().Context(), cacheKey, &incidents) {
		return c.JSON(http.StatusOK, incidents)
	}

	from, to, err := parseDataParams(c)
	if err != nil {
		if err.Error() == "data requested exceeds the 1 year retention policy" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot query data older than 1 year."})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date parameters."})
	}

	incidents, total, err := database.GetIncidentsByUserID(c.Request().Context(), h.DB, userID, limit, offset, from, to, qTarget, qError)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get incidents."})
	}

	if incidents == nil {
		incidents = []*models.Incident{}
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	response := dto.PaginatedResponse{
		Data: incidents,
		Meta: dto.Meta{
			CurrentPage: page,
			Perpage:     limit,
			Total:       total,
			LastPage:    lastPage,
		},
	}

	go h.setCache(c.Request().Context(), cacheKey, incidents, 30*time.Second)

	return c.JSON(http.StatusOK, response)
}

// @Summary Export incidents to CSV
// @Description Download a CSV file with incident history for all monitors.
// @Tags incidents
// @Security BearerAuth
// @Success 200 {file} string "CSV content"
// @Router /incidents/export [get]
func (h *Handler) ExportIncidentsCSV(c echo.Context) error {
	userID := getUserIdFromToken(c)

	from, to, err := parseDataParams(c)
	if err != nil {
		if err.Error() == "data requested exceeds the 1 year retention policy" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot query data older than 1 year."})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date parameters."})
	}

	filename := fmt.Sprintf("incidents_report_%s.csv", time.Now().Format("20060102"))
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+filename)
	c.Response().WriteHeader(http.StatusOK)

	rows, err := database.ExportIncidents(c.Request().Context(), h.DB, userID, from, to)
	if err != nil {
		return nil
	}
	defer rows.Close()

	writer := csv.NewWriter(c.Response().Writer)
	defer writer.Flush()

	writer.Write([]string{"ID", "Target", "Type", "Started At", "Resolved At", "Duration", "Error Cause"})

	for rows.Next() {
		var id int
		var target, mType string
		var startedAt time.Time
		var resolvedAt *time.Time
		var duration *time.Duration
		var errorCause string

		if err := rows.Scan(&id, &target, &mType, &startedAt, &resolvedAt, &duration, &errorCause); err != nil {
			continue
		}

		resolvedStr := "On going"
		if resolvedAt != nil {
			resolvedStr = resolvedAt.Format(time.RFC3339)
		}

		durationStr := "-"
		if duration != nil {
			durationStr = duration.Round(time.Second).String()
		}

		writer.Write([]string{
			fmt.Sprintf("%d", id),
			target,
			mType,
			startedAt.Format(time.RFC3339),
			resolvedStr,
			durationStr,
			errorCause,
		})
	}

	return nil
}

func (h *Handler) GetIncidentsSummary(c echo.Context) error {
	userID := getUserIdFromToken(c)

	summary, err := database.GetIncidentSummary(c.Request().Context(), h.DB, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch incidents summary."})
	}

	return c.JSON(http.StatusOK, summary)
}
