package handlers

import (
	"errors"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/dto"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
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

func (h *Handler) GetUser(c echo.Context) error {
	userID := getUserIdFromToken(c)

	user, err := database.GetUserByID(c.Request().Context(), h.DB, userID)
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
	userID := getUserIdFromToken(c)

	if err := database.DeleteUser(c.Request().Context(), h.DB, userID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found."})
	}

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

	return c.NoContent(http.StatusOK)
}
