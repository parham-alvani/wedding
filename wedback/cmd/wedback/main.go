package main

import (
	"github.com/parham-alvani/wedding/wedback/internal/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
