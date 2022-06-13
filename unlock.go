package main

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
)

type unlockCmd struct {
	out string
}

func (*unlockCmd) Name() string     { return "unlock" }
func (*unlockCmd) Synopsis() string { return "unlock crypt files" }
func (*unlockCmd) Usage() string {
	return `unlock [OPTS] FILE:
  Unlock crypt files.
`
}

func (p *unlockCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.out, "o", "", "output directory or file name")
}

func (p *unlockCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	pw, err := getpasswd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	if p.out == "" {
		p.out = "."
	} else {
		os.MkdirAll(p.out, os.ModePerm)
	}
	err = unlock(p.out, pw, f.Args()...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func unlock(out string, pw []byte, files ...string) error {
	for _, arg := range files {
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			return err
		}
		udata, err := Decrypt(pw, data)
		if err != nil {
			return err
		}
		tr := tar.NewReader(bytes.NewReader(udata))
		if err = extract(out, tr, pw); err != nil {
			return err
		}
	}
	return nil
}

func extract(base string, tr *tar.Reader, pw []byte) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeReg:
			data, err := io.ReadAll(tr)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(filepath.Join(base, hdr.Name), data, fs.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
		case tar.TypeDir:
			err = os.MkdirAll(filepath.Join(base, hdr.Name), fs.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
		default:
			// TODO
			return errors.New("unsupported tar filetype")
		}
	}

	return nil
}
