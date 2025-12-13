package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
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

	from, to, err := parseDataParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date parameters."})
	}

	incidents, err := database.GetIncidentsByUserID(c.Request().Context(), h.DB, userID, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get incidents."})
	}

	return c.JSON(http.StatusOK, incidents)
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
