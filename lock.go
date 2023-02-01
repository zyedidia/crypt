package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"unicode"

	"github.com/google/subcommands"
	"golang.org/x/term"
)

type lockCmd struct {
	outname string
}

func (*lockCmd) Name() string     { return "lock" }
func (*lockCmd) Synopsis() string { return "lock files or directories" }
func (*lockCmd) Usage() string {
	return `lock [OPTS] FILES:
  Archive and encrypt files or directories.
`
}

func (p *lockCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.outname, "o", "", "output crypt file name")
}

func (p *lockCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	pw, err := getpasswd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	if p.outname == "" {
		if len(f.Args()) == 1 && f.Args()[0] != "" && isword(f.Args()[0][0]) {
			p.outname = f.Args()[0] + ".crypt"
		} else {
			p.outname = "archive.crypt"
		}
	}

	err = lock(p.outname, pw, f.Args()...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func lockFrom(outname string, pw []byte, hdr *tar.Header, data []byte) error {
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	tw := tar.NewWriter(zw)
	err := archiveBytes(data, hdr, tw)
	if err != nil {
		return err
	}
	tw.Close()
	zw.Close()

	ldata, err := Encrypt(pw, buf.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outname, ldata, 0666)
}

func lock(outname string, pw []byte, files ...string) error {
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	tw := tar.NewWriter(zw)
	for _, arg := range files {
		err := archive("", arg, tw)
		if err != nil {
			return err
		}
	}
	tw.Close()
	zw.Close()

	ldata, err := Encrypt(pw, buf.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outname, ldata, 0666)
}

func archiveBytes(data []byte, hdr *tar.Header, tw *tar.Writer) error {
	// TODO: link target
	hdr.Size = int64(len(data))
	err := tw.WriteHeader(hdr)
	if err != nil {
		return err
	}
	_, err = tw.Write(data)
	return err
}

func archive(base, path string, tw *tar.Writer) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	// TODO: link target
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	hdr.Name = base + "/" + fi.Name()
	err = tw.WriteHeader(hdr)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		ents, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, ent := range ents {
			e := archive(base+"/"+fi.Name(), path+"/"+ent.Name(), tw)
			if e != nil {
				err = e
			}
		}
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = tw.Write(data)
	return err
}

func getpasswd() ([]byte, error) {
	fmt.Print("password:")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return pw, err
}

func isword(b byte) bool {
	r := rune(b)
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}
