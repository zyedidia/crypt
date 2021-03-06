package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	// subcommands.Register(subcommands.FlagsCommand(), "")
	// subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&lockCmd{}, "")
	subcommands.Register(&unlockCmd{}, "")
	subcommands.Register(&opCmd{}, "")
	subcommands.Register(&versionCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
