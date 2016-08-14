package command

import (
	"fmt"
	"strings"

	mailfull "github.com/directorz/mailfull-go"
)

// AliasUserAddCommand represents a AliasUserAddCommand.
type AliasUserAddCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *AliasUserAddCommand) Synopsis() string {
	return "Create a new aliasuser."
}

// Help returns long-form help text.
func (c *AliasUserAddCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s address target [target...]

Description:
    %s

Required Args:
    address
        The email address that you want to create.
    target
        Target email addresses.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *AliasUserAddCommand) Run(args []string) int {
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
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	aliasUser, err := mailfull.NewAliasUser(aliasUserName, targets)
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if err := repo.AliasUserCreate(domainName, aliasUser); err != nil {
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
