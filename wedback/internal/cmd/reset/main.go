// Package reset clears a guest's RSVP so they can answer again. Plans change
// in the months before a wedding, and the alternative is editing SQLite by
// hand.
package reset

import (
	"context"
	"errors"

	"github.com/parham-alvani/wedding/wedback/internal/cmd/app"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

type params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Service    service.GuestSvc
}

func run(p params, id string) {
	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		go func() {
			defer func() { _ = p.Shutdowner.Shutdown() }()

			guest, err := p.Service.Get(ctx, id)
			if err != nil {
				if errors.Is(err, guestrepo.ErrGuestNotFound) {
					pterm.Error.Printfln("no guest with id %q", id)

					return
				}

				pterm.Error.Printfln("%s", err)

				return
			}

			if guest.Answer == nil {
				pterm.Info.Printfln("%s %s has not answered yet, nothing to reset.", guest.FirstName, guest.LastName)

				return
			}

			if err := p.Service.ResetAnswer(ctx, id); err != nil {
				pterm.Error.Printfln("%s", err)

				return
			}

			pterm.Success.Printfln(
				"Cleared the answer for %s %s; they can reply again.",
				guest.FirstName, guest.LastName,
			)
		}()

		return nil
	}))
}

// Register registers the reset command.
func Register() *cli.Command {
	//nolint: exhaustruct_v5
	return &cli.Command{
		Name:        "reset",
		Usage:       "reset <guest-id>",
		Description: "Clear a guest's RSVP so they can answer again",
		Action: func(_ context.Context, cmd *cli.Command) error {
			id := cmd.Args().First()
			if id == "" {
				return cli.Exit("usage: wedback reset <guest-id>", 1)
			}

			return app.Run(
				app.Providers(),
				fx.Invoke(func(p params) { run(p, id) }),
			)
		},
	}
}
