package main

import (
	"bytes"
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
	pipe  bool
	write bool
}

func (*opCmd) Name() string     { return "op" }
func (*opCmd) Synopsis() string { return "operate on a crypt file" }
func (*opCmd) Usage() string {
	return `op FILE COMMAND:
  Operate on a crypt file.
`
}

func (p *opCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.pipe, "pipe", false, "pipe encrypted file instead of using temp files")
	f.BoolVar(&p.write, "write", false, "allow writing the encrypted file in pipe mode")
}

func (p *opCmd) ExecutePipe(f *flag.FlagSet) subcommands.ExitStatus {
	if len(f.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		return subcommands.ExitFailure
	}

	pw, err := getpasswd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	cryptf := f.Args()[0]
	cmdline := f.Args()[1:]
	buf := &bytes.Buffer{}
	hdr, err := unlockTo(buf, pw, cryptf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "op pipe: %v", err)
		return subcommands.ExitFailure
	}
	output := &bytes.Buffer{}
	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Stdin = buf
	if !p.write {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = output
	}
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	if !p.write {
		return subcommands.ExitSuccess
	}
	if err = lockFrom(cryptf, pw, hdr, output.Bytes()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (p *opCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(f.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		return subcommands.ExitFailure
	}

	if p.pipe {
		return p.ExecutePipe(f)
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
