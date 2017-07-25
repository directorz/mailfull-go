package main

import (
	"fmt"
	"strings"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdUserDel represents a CmdUserDel.
type CmdUserDel struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdUserDel) Synopsis() string {
	return "Delete and backup a user."
}

// Help returns long-form help text.
func (c *CmdUserDel) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] address

Description:
    %s

Required Args:
    address
        The email address that you want to delete.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdUserDel) Run(args []string) int {
	noCommit, err := noCommitFlag(&args)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	if len(args) != 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	address := args[0]
	words := strings.Split(address, "@")
	if len(words) != 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	userName := words[0]
	domainName := words[1]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if userName == "postmaster" {
		c.Meta.Errorf("Cannot delete postmaster.\n")
		return 1
	}

	if err := repo.UserRemove(domainName, userName); err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if noCommit {
		return 0
	}
	if err = repo.GenerateDatabases(); err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	return 0
}
