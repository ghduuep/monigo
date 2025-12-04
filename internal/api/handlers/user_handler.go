package handlers

import (
	"net/http"
	"strconv"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUsers(c echo.Context) error {

	users, err := database.GetAllUsers(c.Request().Context(), h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var response []dto.UserResponse
	for _, u := range users {
		response = append(response, dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	user, err := database.GetUserByID(c.Request().Context(), h.DB, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found."})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *Handler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	if err = database.DeleteUser(c.Request().Context(), h.DB, id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found."})
	}

	return c.NoContent(http.StatusOK)
}
