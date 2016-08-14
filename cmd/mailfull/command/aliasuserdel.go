package command

import (
	"fmt"
	"strings"

	mailfull "github.com/directorz/mailfull-go"
)

// AliasUserDelCommand represents a AliasUserDelCommand.
type AliasUserDelCommand struct {
	Meta
}

// Synopsis returns a one-line synopsis.
func (c *AliasUserDelCommand) Synopsis() string {
	return "Delete a aliasuser."
}

// Help returns long-form help text.
func (c *AliasUserDelCommand) Help() string {
	txt := fmt.Sprintf(`
Usage:
    %s %s address

Description:
    %s

Required Args:
    address
        The email address that you want to delete.
`,
		c.CmdName, c.SubCmdName,
		c.Synopsis())

	return txt[1:]
}

// Run runs the command and returns the exit status.
func (c *AliasUserDelCommand) Run(args []string) int {
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
	aliasUserName := words[0]
	domainName := words[1]

	repo, err := mailfull.OpenRepository(".")
	if err != nil {
		fmt.Fprintf(c.UI.ErrorWriter, "[ERR] %v\n", err)
		return 1
	}

	if err := repo.AliasUserRemove(domainName, aliasUserName); err != nil {
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
