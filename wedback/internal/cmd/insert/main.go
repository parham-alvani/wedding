package serve

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main(logger *zap.Logger, _ *echo.Echo) {
	logger.Info("welcome to our server")
}

// Register server command.
func Register() *cli.Command {
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "insert",
		Description: "Insert a new guest",
		Action: func(_ context.Context, _ *cli.Command) error {
			fx.New(
				fx.Provide(config.Provide),
				fx.Provide(logger.Provide),
				fx.Provide(db.Provide),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
