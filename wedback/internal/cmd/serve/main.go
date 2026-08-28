package serve

import (
	"context"

	"github.com/labstack/echo/v5"
	"github.com/parham-alvani/wedding/wedback/internal/cmd/app"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/server"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main(logger *zap.Logger, _ *echo.Echo) {
	logger.Info("welcome to our server")
}

// Register serve command.
func Register() *cli.Command {
	//nolint: exhaustruct_v5
	return &cli.Command{
		Name:        "serve",
		Description: "Run server to serve the requests",
		Action: func(_ context.Context, _ *cli.Command) error {
			return app.Run(
				app.Providers(),
				fx.Provide(server.Provide),
				fx.Invoke(main),
			)
		},
	}
}
