package httphdl

import (
	"github.com/labstack/echo/v4"
)

// @Tag HealthCheck
// @Accept json
// @Produce json
// @Success 200 {object} domain.Response{}
// @Router /healthcheck [GET]
func (hdl *HTTPHandler) HealthCheck(ctx echo.Context) error {
	resp := map[string]string{"message": "success"}
	return ctx.JSON(200, resp)
}
