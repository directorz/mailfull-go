package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
)

// CommitCommand represents a CommitCommand.
type CommitCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *CommitCommand) Synopsis() string {
	return "Create databases from the structure of the MailData directory."
}

// Help returns long-form help text.
func (c *CommitCommand) Help() string {
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
func (c *CommitCommand) Run(args []string) int {
	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
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
