package httphdl

import (
	"net/http"
	"ticket/internal/core/domain"
	"ticket/pkg/errors"

	"github.com/labstack/echo/v4"
)

func (hdl *HTTPHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var request domain.LoginRequest
	if err := c.Bind(&request); err != nil {
		return errors.NewResponseErr(err)
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.authService.Login(ctx, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) RefreshAccessToken(c echo.Context) error {
	ctx := c.Request().Context()

	var request domain.RefreshTokenRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	resp, err := hdl.authService.RefreshAccessToken(ctx, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (hdl *HTTPHandler) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	token, ok := ctx.Value(domain.ContextKeyToken).(string)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	request := domain.LogoutRequest{
		AccessToken: token,
	}

	if err := hdl.validator.StrcutWithTranslateError(request); err != nil {
		return errors.NewResponseErr(errors.ErrValidation, errors.OptionErr(err))
	}

	err := hdl.authService.Logout(ctx, &request)
	if err != nil {
		return errors.NewResponseErr(err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}
