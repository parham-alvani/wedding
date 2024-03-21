package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"
)

func Execute() {
	cmd := &cli.Command{
		Name:        "wedback",
		Description: "Parham and Elahe's Wedding Backend",
		Authors: []any{
			"Parham Alvani <parham.alvani@gmail.com>",
			"Elahe Dastan <elahe.dstn@gmail.com>",
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
