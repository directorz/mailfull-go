package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdInit represents a CmdInit.
type CmdInit struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdInit) Synopsis() string {
	return "Initializes current directory as a Mailfull repository."
}

// Help returns long-form help text.
func (c *CmdInit) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s

Description:
    %s
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdInit) Run(args []string) int {
	if err := mailfull.InitRepository("."); err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	return 0
}
