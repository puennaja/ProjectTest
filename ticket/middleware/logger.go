package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Logger(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		var body, resp interface{}

		_ = json.Unmarshal(reqBody, &body)
		_ = json.Unmarshal(resBody, &resp)

		if c.Response().Status >= 400 {
			logger.Warnw(fmt.Sprintf("[*] request method %s path %s", c.Request().Method, c.Request().URL),
				"method", c.Request().Method,
				"path", c.Request().URL,
				"body", body,
				"resp", resp,
				"ip", c.RealIP(),
			)
		} else if c.Response().Status > 499 {
			logger.Errorw(fmt.Sprintf("[!] request method %s path %s", c.Request().Method, c.Request().URL),
				"method", c.Request().Method,
				"path", c.Request().URL,
				"body", body,
				"resp", resp,
				"ip", c.RealIP(),
			)
		} else {
			logger.Errorw(fmt.Sprintf("[~] request method %s path %s", c.Request().Method, c.Request().URL),
				"method", c.Request().Method,
				"path", c.Request().URL,
				"body", body,
				"resp", resp,
				"ip", c.RealIP(),
			)
		}
	})
}
