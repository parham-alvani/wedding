package serve

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/http/server"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main(logger *zap.Logger, _ *echo.Echo) {
	logger.Info("welcome to our server")
}

// Register serve command.
func Register() *cli.Command {
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "serve",
		Description: "Run server to serve the requests",
		Action: func(_ context.Context, _ *cli.Command) error {
			fx.New(
				fx.NopLogger,
				fx.Provide(config.Provide),
				fx.Provide(logger.Provide),
				fx.Provide(db.Provide),
				fx.Provide(
					fx.Annotate(repository.ProvideGuestDB, fx.As(new(guestrepo.Repository))),
				),
				fx.Provide(server.Provide),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
