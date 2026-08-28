// Package notes reports the free-text parts of the RSVPs: what the kitchen
// needs to know, and what the guests want to hear.
package notes

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/parham-alvani/wedding/wedback/internal/cmd/app"
	"github.com/parham-alvani/wedding/wedback/internal/domain/model"
	"github.com/parham-alvani/wedding/wedback/internal/domain/repository/guestrepo"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
)

type params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Repository guestrepo.Repository
}

// entry pairs a guest's name with one of their notes.
type entry struct {
	who  string
	note string
}

func name(guest model.Guest) string {
	if guest.SpouseFirstName != nil && *guest.SpouseFirstName != "" {
		return guest.FirstName + " & " + *guest.SpouseFirstName + " " + guest.LastName
	}

	return guest.FirstName + " " + guest.LastName
}

// collect gathers one note per guest, skipping the ones who left it blank.
func collect(guests []model.Guest, get func(*model.Answer) string) []entry {
	out := make([]entry, 0, len(guests))

	for _, guest := range guests {
		if guest.Answer == nil {
			continue
		}

		if note := strings.TrimSpace(get(guest.Answer)); note != "" {
			out = append(out, entry{who: name(guest), note: note})
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].who < out[j].who })

	return out
}

func render(title string, entries []entry, empty string) {
	pterm.DefaultSection.Println(title)

	if len(entries) == 0 {
		pterm.Info.Println(empty)

		return
	}

	rows := make([][]string, 0, len(entries)+1)
	rows = append(rows, []string{"Guest", title})

	for _, e := range entries {
		rows = append(rows, []string{e.who, e.note})
	}

	if err := pterm.DefaultTable.WithHasHeader().WithData(rows).WithWriter(os.Stdout).Render(); err != nil {
		pterm.Error.Printfln("%s", err)
	}
}

func run(p params, only string) {
	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		go func() {
			defer func() { _ = p.Shutdowner.Shutdown() }()

			guests, err := p.Repository.List(ctx)
			if err != nil {
				pterm.Error.Printfln("cannot read guests: %s", err)

				return
			}

			if only != "songs" {
				render("Dietary requirements",
					collect(guests, func(a *model.Answer) string { return a.Dietary }),
					"Nobody has told us about a dietary requirement yet.")
			}

			if only != "dietary" {
				render("Song requests",
					collect(guests, func(a *model.Answer) string { return a.Song }),
					"No song requests yet.")
			}
		}()

		return nil
	}))
}

// Register registers the notes command.
func Register() *cli.Command {
	//nolint: exhaustruct_v5
	return &cli.Command{
		Name:        "notes",
		Usage:       "notes [--only dietary|songs]",
		Description: "Show the dietary requirements and song requests guests left with their RSVP",
		Flags: []cli.Flag{
			//nolint: exhaustruct_v5
			&cli.StringFlag{
				Name:  "only",
				Usage: "limit the report to `dietary` or `songs`",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			only := strings.ToLower(strings.TrimSpace(cmd.String("only")))

			switch only {
			case "", "dietary", "songs":
			default:
				return cli.Exit(fmt.Sprintf("--only must be dietary or songs, got %q", only), 1)
			}

			return app.Run(app.Providers(), fx.Invoke(func(p params) { run(p, only) }))
		},
	}
}
