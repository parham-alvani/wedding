package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/parham-alvani/wedding/wedback/internal/cmd/csvio"
	"github.com/parham-alvani/wedding/wedback/internal/cmd/insert"
	"github.com/parham-alvani/wedding/wedback/internal/cmd/list"
	"github.com/parham-alvani/wedding/wedback/internal/cmd/serve"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/urfave/cli/v3"
)

func Execute() {
	// nolint: exhaustruct_v5
	cmd := &cli.Command{
		Name:        "wedback",
		Description: "Parham and Elaheh's Wedding Backend (fork and customize for your own wedding!)",
		Authors: []any{
			"Parham Alvani <parham.alvani@gmail.com>",
			"Elaheh Dastan <elahe.dstn@gmail.com>",
		},
		Before: func(ctx context.Context, _ *cli.Command) (context.Context, error) {
			// The banner is decoration, not data: keeping it on stderr means
			// `wedback export > guests.csv` produces a clean CSV.
			center := pterm.DefaultCenter.WithWriter(os.Stderr)

			center.Println("Elaheh and Parham's Wedding")

			s, _ := pterm.DefaultBigText.WithLetters(putils.LettersFromString("Wedding")).Srender()
			center.Println(s)

			center.WithCenterEachLineSeparately().Println("Parham Alvani\nApril 2024")

			return ctx, nil
		},
		Commands: []*cli.Command{
			serve.Register(),
			list.Register(),
			insert.Register(),
			csvio.RegisterImport(),
			csvio.RegisterExport(),
		},
		Version: func() string {
			revision := ""
			timestamp := ""
			modified := ""

			if info, ok := debug.ReadBuildInfo(); ok {
				for _, setting := range info.Settings {
					switch setting.Key {
					case "vcs.revision":
						revision = setting.Value
					case "vcs.time":
						timestamp = setting.Value
					case "vcs.modified":
						modified = setting.Value
					}
				}
			}

			if revision == "" {
				return ""
			}

			if modified == "true" {
				return fmt.Sprintf("%s (%s) [dirty]", revision, timestamp)
			}

			return fmt.Sprintf("%s (%s)", revision, timestamp)
		}(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
