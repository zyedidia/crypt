package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

var kdf *int

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	// subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&lockCmd{}, "")
	subcommands.Register(&unlockCmd{}, "")
	subcommands.Register(&opCmd{}, "")
	subcommands.Register(&versionCmd{}, "")

	kdf = flag.Int("kdf", 32768, "number of iterations to use for key derivation")

	flag.Parse()

	if kdf != nil {
		if *kdf%2 != 0 {
			fmt.Fprintln(os.Stderr, "error: kdf must be a power of 2")
			os.Exit(1)
		}
		if *kdf < 32768 {
			fmt.Fprintln(os.Stderr, "warning: kdf is less than 32768")
		}
	}

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
