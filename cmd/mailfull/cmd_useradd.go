package main

import (
	"fmt"
	"strings"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdUserAdd represents a CmdUserAdd.
type CmdUserAdd struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdUserAdd) Synopsis() string {
	return "Create a new user."
}

// Help returns long-form help text.
func (c *CmdUserAdd) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] address

Description:
    %s

Required Args:
    address
        The email address that you want to create.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdUserAdd) Run(args []string) int {
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

	user, err := mailfull.NewUser(userName, mailfull.NeverMatchHashedPassword, nil)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if err := repo.UserCreate(domainName, user); err != nil {
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
