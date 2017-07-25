package main

import (
	"fmt"

	mailfull "github.com/directorz/mailfull-go"
)

// CmdAliasDomainAdd represents a CmdAliasDomainAdd.
type CmdAliasDomainAdd struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdAliasDomainAdd) Synopsis() string {
	return "Create a new aliasdomain."
}

// Help returns long-form help text.
func (c *CmdAliasDomainAdd) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] domain target

Description:
    %s

Required Args:
    domain
        The domain name that you want to create.
    target
        The target domain name.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdAliasDomainAdd) Run(args []string) int {
	noCommit, err := noCommitFlag(&args)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	if len(args) != 2 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	aliasDomainName := args[0]
	targetDomainName := args[1]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	aliasDomain, err := mailfull.NewAliasDomain(aliasDomainName, targetDomainName)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if err := repo.AliasDomainCreate(aliasDomain); err != nil {
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
