package serve

import (
	"context"

	"github.com/labstack/echo/v4"
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
		Name:        "serve",
		Description: "Run server to serve the requests",
		Action: func(_ context.Context, _ *cli.Command) error {
			fx.New(
				fx.Provide(logger.Provide),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
