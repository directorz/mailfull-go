package main

import (
	"fmt"

	"github.com/directorz/mailfull-go"
	"github.com/directorz/mailfull-go/cmd"
)

// CmdDomainEnable represents a CmdDomainEnable.
type CmdDomainEnable struct {
	cmd.Meta
}

// Synopsis returns a one-line synopsis.
func (c *CmdDomainEnable) Synopsis() string {
	return "Enable a domain."
}

// Help returns long-form help text.
func (c *CmdDomainEnable) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] domain

Description:
    %s

Required Args:
    domain
        The domain name that you want to enable.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *CmdDomainEnable) Run(args []string) int {
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

	domain, err := repo.Domain(domainName)
	if err != nil {
		c.Meta.Errorf("%v\n", err)
		return 1
	}
	if domain == nil {
		c.Meta.Errorf("%v\n", mailfull.ErrDomainNotExist)
		return 1
	}

	domain.SetDisabled(false)

	if err := repo.DomainUpdate(domain); err != nil {
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
