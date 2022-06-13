package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/subcommands"
)

type opCmd struct {
}

func (*opCmd) Name() string     { return "op" }
func (*opCmd) Synopsis() string { return "operate on a crypt file" }
func (*opCmd) Usage() string {
	return `op FILE COMMAND:
  Operate on a crypt file.
`
}

func (p *opCmd) SetFlags(f *flag.FlagSet) {
}

func (p *opCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		return subcommands.ExitFailure
	}

	pw, err := getpasswd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	cryptf := f.Args()[0]
	cmdline := f.Args()[1:]
	if err = unlock(tmp, pw, cryptf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Dir = tmp
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	ents, err := os.ReadDir(tmp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	files := make([]string, 0, len(ents))
	for _, ent := range ents {
		files = append(files, filepath.Join(tmp, ent.Name()))
	}

	if err = lock(cryptf, pw, files...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
