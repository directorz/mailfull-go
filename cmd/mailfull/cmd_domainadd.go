package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdDomainAdd represents a CmdDomainAdd.
type CmdDomainAdd struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdDomainAdd) Synopsis() string {
	return "Create a new domain and postmaster."
}

// Help returns long-form help text.
func (c *CmdDomainAdd) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] domain

Description:
    %s

Required Args:
    domain
        The domain name that you want to create.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdDomainAdd) Run(args []string) int {
	noCommit, err := noCommitFlag(&args)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	if len(args) != 1 {
		fmt.Fprintf(c.UI.ErrorWriter, "%v\n", c.Help())
		return 1
	}

	domainName := args[0]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	domain, err := mailfull.NewDomain(domainName)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	if err := repo.DomainCreate(domain); err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}

	user, err := mailfull.NewUser("postmaster", mailfull.NeverMatchHashedPassword, nil)
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
