package csvio

import (
	"context"
	"os"
	"strings"

	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/domain/service"
	"github.com/parham-alvani/wedding/wedback/internal/infra/config"
	"github.com/parham-alvani/wedding/wedback/internal/infra/db"
	"github.com/parham-alvani/wedding/wedback/internal/infra/generator"
	"github.com/parham-alvani/wedding/wedback/internal/infra/logger"
	"github.com/parham-alvani/wedding/wedback/internal/infra/repository"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

type importParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Shutdowner  fx.Shutdowner
	Service     service.GuestSvc
	WeddingConf wedding.Config
}

func readRows(path string) ([]Row, error) {
	// The path comes from the operator's own command line, not from a request.
	file, err := os.Open(path) // nolint: gosec
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	defer func() { _ = file.Close() }()

	return Parse(file)
}

func reportDryRun(rows []Row) {
	pterm.Info.Printfln("Dry run: %d guest(s) parsed, nothing written.", len(rows))

	for _, row := range rows {
		pterm.Info.Printfln("line %d: %s %s", row.Line, row.FirstName, row.LastName)
	}
}

func createRows(ctx context.Context, p importParams, rows []Row) {
	created, failed := 0, 0

	for _, row := range rows {
		guest, err := p.Service.New(
			ctx,
			row.FirstName,
			row.LastName,
			row.SpouseFirstName,
			row.SpouseLastName,
			row.IsFamily,
			row.Children,
		)
		if err != nil {
			failed++

			pterm.Warning.Printfln("line %d (%s %s): %s", row.Line, row.FirstName, row.LastName, err)

			continue
		}

		created++

		pterm.Success.Printfln(
			"%s %s -> %s/guests/%s",
			guest.FirstName, guest.LastName, strings.TrimSuffix(p.WeddingConf.BaseURL, "/"), guest.ID,
		)
	}

	pterm.Info.Printfln("Imported %d guest(s), %d skipped.", created, failed)
}

func runImport(p importParams, path string, dryRun bool) {
	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		go func() {
			defer func() { _ = p.Shutdowner.Shutdown() }()

			rows, err := readRows(path)
			if err != nil {
				pterm.Error.Printfln("%s", err)

				return
			}

			if dryRun {
				reportDryRun(rows)

				return
			}

			createRows(ctx, p, rows)
		}()

		return nil
	}))
}

type exportParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Shutdowner  fx.Shutdowner
	Repository  guestrepo.Repository
	WeddingConf wedding.Config
}

func writeExport(path string, guests []model.Guest, baseURL string) error {
	if path == "" {
		return Export(os.Stdout, guests, baseURL)
	}

	// The path comes from the operator's own command line, not from a request.
	file, err := os.Create(path) // nolint: gosec
	if err != nil {
		return err // nolint: wrapcheck
	}

	defer func() { _ = file.Close() }()

	if err := Export(file, guests, baseURL); err != nil {
		return err
	}

	pterm.Success.Printfln("Wrote %d guest(s) to %s", len(guests), path)

	return nil
}

func runExport(p exportParams, path string, withLinks bool) {
	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		go func() {
			defer func() { _ = p.Shutdowner.Shutdown() }()

			guests, err := p.Repository.List(ctx)
			if err != nil {
				pterm.Error.Printfln("cannot read guests: %s", err)

				return
			}

			baseURL := ""
			if withLinks {
				baseURL = p.WeddingConf.BaseURL
			}

			if err := writeExport(path, guests, baseURL); err != nil {
				pterm.Error.Printfln("%s", err)
			}
		}()

		return nil
	}))
}

func providers() fx.Option {
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

// RegisterImport registers the bulk import command.
func RegisterImport() *cli.Command {
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "import",
		Usage:       "import <file.csv>",
		Description: "Bulk import guests from a CSV file (columns: " + strings.Join(Columns(), ", ") + ")",
		Flags: []cli.Flag{
			//nolint: exhaustruct
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "parse and report without writing anything",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			path := cmd.Args().First()
			if path == "" {
				return cli.Exit("usage: wedback import <file.csv>", 1)
			}

			dryRun := cmd.Bool("dry-run")

			fx.New(providers(), fx.Invoke(func(p importParams) {
				runImport(p, path, dryRun)
			})).Run()

			return nil
		},
	}
}

// RegisterExport registers the CSV export command.
func RegisterExport() *cli.Command {
	//nolint: exhaustruct
	return &cli.Command{
		Name:        "export",
		Usage:       "export [file.csv]",
		Description: "Export every guest and their RSVP as CSV (stdout when no file is given)",
		Flags: []cli.Flag{
			//nolint: exhaustruct
			&cli.BoolFlag{
				Name:  "links",
				Usage: "add a column with each guest's invitation link",
				Value: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			fx.New(providers(), fx.Invoke(func(p exportParams) {
				runExport(p, cmd.Args().First(), cmd.Bool("links"))
			})).Run()

			return nil
		},
	}
}
