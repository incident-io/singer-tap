package main

import (
	"context"

	"github.com/alecthomas/kingpin/v2"
	"github.com/incident-io/singer-tap/cmd/incident-tap/cmd"
)

func main() {
	if err := cmd.Run(context.Background()); err != nil {
		kingpin.Fatalf(err.Error())
	}
}
