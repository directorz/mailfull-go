package command

import (
	"fmt"

	"github.com/directorz/mailfull-go"
)

// DomainDisableCommand represents a DomainDisableCommand.
type DomainDisableCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *DomainDisableCommand) Synopsis() string {
	return "Disable a domain temporarily."
}

// Help returns long-form help text.
func (c *DomainDisableCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s [-n] domain

Description:
    %s

Required Args:
    domain
        The domain name that you want to disable.

Optional Args:
    -n
        Don't update databases.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *DomainDisableCommand) Run(args []string) int {
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
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	domain, err := repo.Domain(domainName)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}
	if domain == nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", mailfull.ErrDomainNotExist)
		return 1
	}

	domain.SetDisabled(true)

	if err := repo.DomainUpdate(domain); err != nil {
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
