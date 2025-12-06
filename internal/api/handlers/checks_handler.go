package handlers

import (
	"net/http"
	"strconv"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetChecksByMonitorID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	checks, err := database.GetChecksByMonitor(c.Request().Context(), h.DB, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid data."})
	}

	return c.JSON(http.StatusOK, checks)
}
