package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdCommit represents a CmdCommit.
type CmdCommit struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdCommit) Synopsis() string {
	return "Create databases from the structure of the MailData directory."
}

// Help returns long-form help text.
func (c *CmdCommit) Help() string {
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
func (c *CmdCommit) Run(args []string) int {
	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
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
