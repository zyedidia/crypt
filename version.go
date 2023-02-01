package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

var Version = "1.2.0"

type versionCmd struct {
	outname string
}

func (*versionCmd) Name() string     { return "version" }
func (*versionCmd) Synopsis() string { return "show version" }
func (*versionCmd) Usage() string {
	return `version:
  Show version
`
}

func (p *versionCmd) SetFlags(f *flag.FlagSet) {
}

func (p *versionCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Println(Version)
	return subcommands.ExitSuccess
}
