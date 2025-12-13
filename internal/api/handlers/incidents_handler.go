package handlers

import (
	"github.com/ghduuep/pingly/internal/database"
	"github.com/labstack/echo/v4"
	"net/http"
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
