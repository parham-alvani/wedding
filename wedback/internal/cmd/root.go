package cmd

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func Execute() {
	cmd := &cli.Command{
		Name:  "wedback",
		Usage: "Parham and Elahe's Wedding Backend",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
