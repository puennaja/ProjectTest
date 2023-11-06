package middleware

import (
	"context"
	"strings"

	"ticket/internal/core/domain"
	"ticket/internal/core/port"
	"ticket/pkg/errors"

	"github.com/labstack/echo/v4"
)

func Auth(authSvc port.AuthService) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			token := c.Request().Header.Get("Authorization")
			token = strings.Replace(token, "Bearer ", "", 1)
			path := strings.Replace(c.Request().URL.Path, "/api/v1", "", 1)
			method := c.Request().Method

			// authentication
			authentication, err := authSvc.Authentication(ctx, token)
			if err != nil {
				return errors.ErrUnauthorized.SetError(err)
			}

			// authorization
			authorization, err := authSvc.Authorization(ctx, authentication.User.Role, method, path)
			if err != nil {
				return errors.ErrUnauthorized.SetError(err)
			}

			if !authorization.GrantAccess {
				return errors.ErrPermissionDenied
			}

			ctxToken := context.WithValue(c.Request().Context(), domain.ContextKeyToken, token)
			ctxAuth := context.WithValue(ctxToken, domain.ContextKeyAuth, authentication)
			c.SetRequest(c.Request().WithContext(ctxAuth))
			return next(c)
		}
	}
}
