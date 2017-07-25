package main

import (
	"fmt"
	"strings"

	"github.com/directorz/mailfull-go"
)

// CmdUserDel represents a CmdUserDel.
type CmdUserDel struct {
	Meta
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
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if userName == "postmaster" {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] Cannot delete postmaster.\n")
		return 1
	}

	if err := repo.UserRemove(domainName, userName); err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if noCommit {
		return 0
	}

	mailData, err := repo.MailData()
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	err = repo.GenerateDatabases(mailData)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	return 0
}
