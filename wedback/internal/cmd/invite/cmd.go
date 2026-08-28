package invite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/parham-alvani/wedding/wedback/internal/cmd/app"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/parham-alvani/wedding/wedback/internal/infra/wedding"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

// qrFileMode is the permission for generated QR files and their directory.
const (
	qrFileMode = 0o600
	qrDirMode  = 0o750
)

type params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Shutdowner  fx.Shutdowner
	Repository  guestrepo.Repository
	WeddingConf wedding.Config
}

type options struct {
	onlyWaiting bool
	guestID     string
	qrDir       string
	separator   string
}

// selectGuests narrows the list down to what the flags asked for.
func selectGuests(guests []model.Guest, opts options) []model.Guest {
	if opts.guestID != "" {
		for _, guest := range guests {
			if guest.ID == opts.guestID {
				return []model.Guest{guest}
			}
		}

		return nil
	}

	if opts.onlyWaiting {
		return Waiting(guests)
	}

	return guests
}

func writeQRCodes(rendered []Rendered, dir string) error {
	if err := os.MkdirAll(dir, qrDirMode); err != nil {
		return fmt.Errorf("cannot create %s: %w", dir, err)
	}

	for _, item := range rendered {
		png, err := QRCode(item.Invitation.Link)
		if err != nil {
			return err
		}

		// The id is generated from an alphanumeric alphabet, so it is always a
		// safe file name.
		path := filepath.Join(dir, item.Invitation.ID+".png")
		if err := os.WriteFile(path, png, qrFileMode); err != nil {
			return fmt.Errorf("cannot write %s: %w", path, err)
		}
	}

	pterm.Success.Printfln("Wrote %d QR code(s) to %s", len(rendered), dir)

	return nil
}

func run(p params, opts options) {
	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		go func() {
			defer func() { _ = p.Shutdowner.Shutdown() }()

			guests, err := p.Repository.List(ctx)
			if err != nil {
				pterm.Error.Printfln("cannot read guests: %s", err)

				return
			}

			selected := selectGuests(guests, opts)

			rendered, err := Render(
				p.WeddingConf.InviteTemplate,
				selected,
				p.WeddingConf.BaseURL,
				p.WeddingConf.HusbandName,
				p.WeddingConf.WifeName,
			)
			if err != nil {
				pterm.Error.Printfln("%s", err)

				return
			}

			// Messages go to stdout so they can be piped or redirected; the
			// separator keeps them machine-splittable.
			for i, item := range rendered {
				if i > 0 {
					_, _ = fmt.Fprintln(os.Stdout, opts.separator)
				}

				_, _ = fmt.Fprintln(os.Stdout, item.Message)
			}

			if opts.qrDir != "" {
				if err := writeQRCodes(rendered, opts.qrDir); err != nil {
					pterm.Error.Printfln("%s", err)

					return
				}
			}

			pterm.Info.Printfln("Rendered %d invitation(s).", len(rendered))
		}()

		return nil
	}))
}

// Register registers the invite command.
func Register() *cli.Command {
	//nolint: exhaustruct_v5
	return &cli.Command{
		Name:  "invite",
		Usage: "invite [guest-id]",
		Description: "Render each guest's invitation message and personal RSVP link. " +
			"Set wedding.invite_template to change the wording; the template is a Go " +
			"text/template over .Names, .FirstName, .LastName, .SpouseFirstName, " +
			".SpouseLastName, .Link, .ID, .IsFamily, .Children, .Husband, .Wife, " +
			".Answered and .Coming.",
		Flags: []cli.Flag{
			//nolint: exhaustruct_v5
			&cli.BoolFlag{
				Name:  "waiting",
				Usage: "only guests who have not replied yet",
			},
			//nolint: exhaustruct_v5
			&cli.StringFlag{
				Name:  "qr-dir",
				Usage: "also write a QR code PNG per guest into this directory",
			},
			//nolint: exhaustruct_v5
			&cli.StringFlag{
				Name:  "separator",
				Usage: "printed between messages",
				Value: "---",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			opts := options{
				onlyWaiting: cmd.Bool("waiting"),
				guestID:     strings.TrimSpace(cmd.Args().First()),
				qrDir:       cmd.String("qr-dir"),
				separator:   cmd.String("separator"),
			}

			return app.Run(app.Providers(), fx.Invoke(func(p params) { run(p, opts) }))
		},
	}
}
