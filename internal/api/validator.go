package api

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, formatValidationError(err))
	}
	return nil
}

func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			if e.Param() != "" {
				errors[e.Field()] = fmt.Sprintf("%s=%s", e.Tag(), e.Param())
			} else {
				errors[e.Field()] = e.Tag()
			}
		}
	}
	return errors
}
