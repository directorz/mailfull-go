package command

import (
	"fmt"

	"github.com/directorz/mailfull-go"
)

// DomainDelCommand represents a DomainDelCommand.
type DomainDelCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *DomainDelCommand) Synopsis() string {
	return "Delete and backup a domain."
}

// Help returns long-form help text.
func (c *DomainDelCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s domain

Description:
    %s

Required Args:
    domain
        The domain name that you want to delete.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *DomainDelCommand) Run(args []string) int {
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

	if err := repo.DomainRemove(domainName); err != nil {
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
