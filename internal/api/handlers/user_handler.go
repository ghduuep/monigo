package handlers

import (
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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

func (h *Handler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID must be a number."})
	}

	userID := getUserIdFromToken(c)

	var dto dto.UpdateUserRequest
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data."})
	}

	ctx := c.Request().Context()

	if dto.Email != nil {
		ownerID, err := database.GetUserIDByEmail(ctx, h.DB, *dto.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate e-mail."})
		}

		if ownerID != 0 && ownerID != userID {
			return c.JSON(http.StatusConflict, map[string]string{"error": "E-mail already used."})
		}
	}

	if dto.Password != nil {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate the password hash"})
		}

		*dto.Password = string(hashedBytes)
	}

	err = database.UpdateUser(ctx, h.DB, id, dto)
	if err != nil {
		if err.Error() == "no data" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user."})
	}

	return c.NoContent(http.StatusOK)
}
