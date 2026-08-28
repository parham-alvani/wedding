// Package app builds and runs the dependency graph shared by every command.
package app

import (
	"fmt"

	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"go.uber.org/fx"
)

// Providers is the object graph every command needs: configuration, logging,
// the database, and the guest repository and service on top of them.
func Providers() fx.Option {
	return fx.Options(
		fx.NopLogger,
		fx.Provide(config.Provide),
		fx.Provide(logger.Provide),
		fx.Provide(db.Provide),
		fx.Provide(
			fx.Annotate(repository.ProvideGuestDB, fx.As(new(guestrepo.Repository))),
		),
		fx.Provide(generator.Provide),
		fx.Provide(service.ProvideGuestSvc),
	)
}

// Run builds the application and runs it. Construction errors are returned
// rather than swallowed: fx.NopLogger keeps the startup noise down, but it
// also hides why a bad configuration value stopped the process.
func Run(opts ...fx.Option) error {
	application := fx.New(opts...)
	if err := application.Err(); err != nil {
		return fmt.Errorf("wedback could not start: %w", err)
	}

	application.Run()

	return nil
}
