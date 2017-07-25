package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
)

// CmdCatchAll represents a CmdCatchAll.
type CmdCatchAll struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdCatchAll) Synopsis() string {
	return "Show a catchall user."
}

// Help returns long-form help text.
func (c *CmdCatchAll) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s domain

Description:
    %s

Required Args:
    domain
        The domain name.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdCatchAll) Run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	domainName := args[0]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	catchAllUser, err := repo.CatchAllUser(domainName)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if catchAllUser != nil {
		fmt.Fprintf(c.UI.Writer, "%s\n", catchAllUser.Name())
	}

	return 0
}
