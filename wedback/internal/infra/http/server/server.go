package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide(lc fx.Lifecycle, logger *zap.Logger, svc service.GuestSvc) *echo.Echo {
	app := echo.New()

	handler.Healthz{
		Logger: logger.Named("handler").Named("healthz"),
	}.Register(app.Group(""))

	handler.Guest{
		Logger:  logger.Named("handler").Named("guest"),
		Service: svc,
	}.Register(app.Group(""))

	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				sc := echo.StartConfig{Address: ":1378"}
				if err := sc.Start(ctx, app); !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("echo initiation failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()

			return nil
		},
	})

	return app
}
