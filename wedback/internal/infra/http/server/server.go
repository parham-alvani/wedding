package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide(lc fx.Lifecycle, logger *zap.Logger, svc service.GuestSvc) *echo.Echo {
	app := echo.New()

	app.Use(middleware.CORS())

	handler.Healthz{
		Logger: logger.Named("handler").Named("healthz"),
	}.Register(app.Group(""))

	handler.Guest{
		Logger:  logger.Named("handler").Named("guest"),
		Service: svc,
	}.Register(app.Group(""))

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := app.Start(":1378"); !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("echo initiation failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: app.Shutdown,
	})

	return app
}
