package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	vv9 "gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *vv9.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
