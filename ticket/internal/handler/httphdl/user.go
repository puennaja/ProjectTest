package httphdl

import (
	"net/http"
	"ticket/internal/core/domain"
	"ticket/pkg/errors"

	"github.com/labstack/echo/v4"
)

func (hdl *HTTPHandler) UserMe(c echo.Context) error {
	ctx := c.Request().Context()
	auth, ok := ctx.Value(domain.ContextKeyAuth).(*domain.AuthenticationResponse)
	if !ok {
		return errors.ErrUnauthorizedContext
	}

	resp, err := hdl.userService.GetUserByID(ctx, auth.User.ID)
	if err != nil {
		return errors.NewResponseErr(err)
	}
	return c.JSON(http.StatusOK, resp)
}
