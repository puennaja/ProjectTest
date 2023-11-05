package middleware

import (
	"ticket/pkg/errors"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if errResponse := errors.IsResponseErr(err); errResponse != nil {
		_ = c.JSON(errResponse.HTTPStatus(), errResponse)
		return
	}

	if errEcho, ok := err.(*echo.HTTPError); ok {
		err := errors.NewResponseErr(errEcho.Unwrap(), errors.OptionHttpStatus(errEcho.Code))
		_ = c.JSON(err.HTTPStatus(), err)
		return
	}

	defaultErr := errors.NewResponseErr(err)
	_ = c.JSON(defaultErr.HTTPStatus(), defaultErr)
}
