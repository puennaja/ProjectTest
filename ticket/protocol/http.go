package protocol

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"ticket/config"
	"ticket/internal/handler/httphdl"
	"ticket/middleware"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func ServeREST() error {
	e := echo.New()
	e.HTTPErrorHandler = middleware.ErrorHandler

	httpHdl := httphdl.NewHTTP(httphdl.Config{
		Validator: app.pkg.validator,
	})
	e.GET("/healthcheck", httpHdl.HealthCheck)
	apiGroup := e.Group("/api/v1")
	{
		apiGroup.Use(
			echomiddleware.CORS(),
			middleware.Logger(app.logger),
		)

		// authGroup := apiGroup.Group("/auth")
		// {
		// 	// Public Route
		// 	authGroup.GET("/login", httpHdl)
		// 	authGroup.POST("/refresh-token", httpHdl.RefreshAccessToken)

		// 	// Private Route
		// 	authGroup.Use(
		// 		middleware.Auth(app.svc.authSvc),
		// 	)
		// 	authGroup.POST("/logout", httpHdl.Logout)
		// }

		userGroup := apiGroup.Group("/user")
		{
			// Private Route
			userGroup.Use(
				middleware.Auth(app.svc.authSvc),
			)
			userGroup.GET("/me", httpHdl.UserMe)
		}
	}

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf("%s:%s", config.GetConfig().Server.Host, config.GetConfig().Server.Port)); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.logger.Info("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
