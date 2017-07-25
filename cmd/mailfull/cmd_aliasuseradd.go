package main

import (
	"fmt"
	"strings"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdAliasUserAdd represents a CmdAliasUserAdd.
type CmdAliasUserAdd struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdAliasUserAdd) Synopsis() string {
	return "Create a new aliasuser."
}

// Help returns long-form help text.
func (c *CmdAliasUserAdd) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] address target [target...]

Description:
    %s

Required Args:
    address
        The email address that you want to create.
    target
        Target email addresses.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdAliasUserAdd) Run(args []string) int {
	noCommit, err := noCommitFlag(&args)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	if len(args) < 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	address := args[0]
	targets := args[1:]

	words := strings.Split(address, "@")
	if len(words) != 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}
	aliasUserName := words[0]
	domainName := words[1]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	aliasUser, err := mailfull.NewAliasUser(aliasUserName, targets)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if err := repo.AliasUserCreate(domainName, aliasUser); err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if noCommit {
		return 0
	}

	mailData, err := repo.MailData()
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	err = repo.GenerateDatabases(mailData)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	return 0
}
