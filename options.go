package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pborman/getopt/v2"
)

type Options struct {
	Help      bool
	Version   bool
	Servers   []string
	Split     bool
	OutputDir string
}

func Help() {
	getopt.PrintUsage(os.Stdout)
	os.Exit(0)
}

func Version() {
    fmt.Printf("masto-emoji-pack v%s\n", VERSION)
	os.Exit(0)
}

func Usage(err error) {
	fmt.Fprintln(os.Stderr, err)
	getopt.Usage()
	os.Exit(2)
}

func Parse() {
	if err := getopt.Getopt(nil); err != nil {
		Usage(err)
	}
}

func parseOptions() (options Options) {
	getopt.SetParameters("DOMAIN...")
	help := getopt.BoolLong("help", 'h', "show help message")
	version := getopt.BoolLong("version", 'v', "show version info")
	split := getopt.BoolLong("split", 's', "split emoji pack via category")
	dir := getopt.StringLong("path", 'p', "/tmp", "generate emoji pack directory", "PATH")

	Parse()

	options = Options{
		Help:      *help,
		Version:   *version,
		Servers:   getopt.Args(),
		Split:     *split,
		OutputDir: filepath.Clean(*dir),
	}

	if options.Help {
		Help()
	}
	if options.Version {
		Version()
	}
	if len(options.Servers) == 0 {
		Usage(errors.New("must be specified: DOMAIN..."))
	}

	return
}
