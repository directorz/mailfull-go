package command

import (
	"bytes"
	"flag"

	"github.com/mitchellh/cli"
)

// Meta is for `*Command` struct.
type Meta struct {
	UI         *cli.BasicUi
	CmdName    string
	SubCmdName string
	Version    string
}

// noCommitFlag returns true if `pargs` has "-n" flag.
// `pargs` is overwrites with non-flag arguments.
func noCommitFlag(pargs *[]string) (bool, error) {
	nFlag := false

	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.SetOutput(&bytes.Buffer{})
	flagSet.BoolVar(&nFlag, "n", nFlag, "")
	err := flagSet.Parse(*pargs)
	*pargs = flagSet.Args()

	return nFlag, err
}
