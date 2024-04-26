package serve

import (
	"context"

	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

func main() {
}

// Register insert command.
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
				fx.Provide(
					fx.Annotate(repository.ProvideGuestDB, fx.As(new(guestrepo.Repository))),
				),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
